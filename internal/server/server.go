package server

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/jwhumphries/blog/internal/compress"
)

// compressibleExtensions lists file extensions that benefit from compression.
var compressibleExtensions = map[string]bool{
	".html": true,
	".css":  true,
	".js":   true,
	".json": true,
	".xml":  true,
	".svg":  true,
	".txt":  true,
	".md":   true,
}

// cachedFile holds both original and compressed versions of a file.
type cachedFile struct {
	content     []byte
	compressed  []byte
	contentType string
}

// Server serves pre-compressed static files from an embedded filesystem.
type Server struct {
	cache  map[string]*cachedFile
	mu     sync.RWMutex
	logger *log.Logger
}

// New creates a new Server by loading and pre-compressing files from the embedded FS.
func New(fsys embed.FS, root string, logger *log.Logger) (*Server, error) {
	s := &Server{
		cache:  make(map[string]*cachedFile),
		logger: logger,
	}

	subFS, err := fs.Sub(fsys, root)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(subFS, path)
		if err != nil {
			return err
		}

		// Normalize path for URL matching
		urlPath := "/" + path
		if strings.HasSuffix(path, "/index.html") {
			// Also cache the directory path
			dirPath := "/" + strings.TrimSuffix(path, "index.html")
			s.cacheFile(dirPath, path, content)
		}

		s.cacheFile(urlPath, path, content)
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info("server initialized", "files_cached", len(s.cache))
	return s, nil
}

// cacheFile adds a file to the cache, optionally pre-compressing it.
func (s *Server) cacheFile(urlPath, filePath string, content []byte) {
	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	cf := &cachedFile{
		content:     content,
		contentType: contentType,
	}

	// Pre-compress if beneficial
	if compressibleExtensions[ext] && len(content) > 1024 {
		compressed, err := compress.Brotli(content)
		if err == nil && len(compressed) < len(content) {
			cf.compressed = compressed
			s.logger.Debug("pre-compressed file",
				"path", urlPath,
				"original", len(content),
				"compressed", len(compressed),
				"ratio", float64(len(compressed))/float64(len(content)),
			)
		}
	}

	s.mu.Lock()
	s.cache[urlPath] = cf
	s.mu.Unlock()
}

// Handler returns an http.Handler that serves the cached files.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Try exact match first
		s.mu.RLock()
		cf, ok := s.cache[path]
		s.mu.RUnlock()

		// Try with trailing slash for directories
		if !ok && !strings.HasSuffix(path, "/") {
			s.mu.RLock()
			cf, ok = s.cache[path+"/"]
			s.mu.RUnlock()
		}

		// Try index.html for directory paths
		if !ok {
			indexPath := strings.TrimSuffix(path, "/") + "/index.html"
			s.mu.RLock()
			cf, ok = s.cache[indexPath]
			s.mu.RUnlock()
		}

		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", cf.contentType)

		// Check if client accepts Brotli and we have compressed version
		if cf.compressed != nil && strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header().Set("Content-Encoding", "br")
			w.Header().Set("Vary", "Accept-Encoding")
			_, _ = w.Write(cf.compressed)
			return
		}

		_, _ = w.Write(cf.content)
	})
}

// HealthHandler returns a simple health check handler.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
