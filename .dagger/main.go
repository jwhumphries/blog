// Build and package the blog server

package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/blog-builder/internal/dagger"
)

type BlogBuilder struct{}

// gitVersion extracts version info from git tags.
func (m *BlogBuilder) gitVersion(ctx context.Context, git *dagger.Directory) (string, error) {
	if git == nil {
		return "dev", nil
	}
	out, err := dag.Container().
		From("alpine/git:latest").
		WithMountedDirectory("/src/.git", git).
		WithWorkdir("/src").
		WithExec([]string{"git", "describe", "--tags", "--always"}).
		Stdout(ctx)
	if err != nil {
		return "dev", nil
	}
	return strings.TrimSpace(out), nil
}

// gitCommit extracts the current commit SHA.
func (m *BlogBuilder) gitCommit(ctx context.Context, git *dagger.Directory) (string, error) {
	if git == nil {
		return "unknown", nil
	}
	out, err := dag.Container().
		From("alpine/git:latest").
		WithMountedDirectory("/src/.git", git).
		WithWorkdir("/src").
		WithExec([]string{"git", "rev-parse", "--short", "HEAD"}).
		Stdout(ctx)
	if err != nil {
		return "unknown", nil
	}
	return strings.TrimSpace(out), nil
}

// HugoBuild builds the Hugo site and returns the public directory.
func (m *BlogBuilder) HugoBuild(source *dagger.Directory) *dagger.Directory {
	tailwind := dag.Container().
		From("ghcr.io/jwhumphries/tailwindcss:latest@sha256:a4fdf32e156f84f0221a77b2c5afc2448a6b143b088df5e2d3e3fa6ac31f4656").
		File("/usr/local/bin/tailwindcss")

	builder := dag.Container().
		From("ghcr.io/gohugoio/hugo:latest@sha256:53dc48ef4d550835b0e54b0f6b41e22e5160e27065d0691b220a713218eb059d").
		WithFile("/usr/local/bin/tailwindcss", tailwind).
		WithUser("root").
		WithDirectory("/src", source.Directory("site")).
		WithWorkdir("/src").
		WithExec([]string{"hugo", "mod", "npm", "pack"}).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"hugo", "--gc", "--minify"})

	return builder.Directory("public")
}

// Lint runs golangci-lint on the Go server code.
func (m *BlogBuilder) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	// Create source with placeholder public dir for embed directive
	sourceWithPublic := dag.Container().
		From("alpine:3.21").
		WithDirectory("/src", source).
		WithExec([]string{"mkdir", "-p", "/src/public"}).
		WithExec([]string{"sh", "-c", "echo '<!DOCTYPE html><html></html>' > /src/public/index.html"}).
		Directory("/src")

	return dag.GolangciLint(dagger.GolangciLintOpts{
		Version: "v2.8.0",
	}).
		WithModuleCache(dag.CacheVolume("go-mod-cache")).
		WithLinterCache(dag.CacheVolume("golangci-lint-cache")).
		Run(sourceWithPublic).
		Stdout(ctx)
}

// Test runs go tests on the server code.
func (m *BlogBuilder) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithEnvVariable("GOCACHE", "/go-build-cache").
		WithEnvVariable("GOMODCACHE", "/go-mod-cache").
		WithMountedCache("/go-build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/go-mod-cache", dag.CacheVolume("go-mod-cache")).
		WithDirectory("/app", source).
		WithWorkdir("/app").
		// Create a minimal public dir for embed to work during tests
		WithExec([]string{"mkdir", "-p", "public"}).
		WithExec([]string{"sh", "-c", "echo '<html></html>' > public/index.html"}).
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}

// BuildBinary compiles the Go server binary with embedded Hugo site.
func (m *BlogBuilder) BuildBinary(source *dagger.Directory, publicDir *dagger.Directory, version, commit string) *dagger.Container {
	ldflags := fmt.Sprintf(
		"-s -w -X github.com/jwhumphries/blog/version.Tag=%s -X github.com/jwhumphries/blog/version.Commit=%s",
		version, commit,
	)

	return dag.Container().
		From("golang:1.25-alpine").
		WithEnvVariable("GOCACHE", "/go-build-cache").
		WithEnvVariable("GOMODCACHE", "/go-mod-cache").
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedCache("/go-build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/go-mod-cache", dag.CacheVolume("go-mod-cache")).
		WithDirectory("/app", source).
		WithDirectory("/app/public", publicDir).
		WithWorkdir("/app").
		WithExec([]string{
			"go", "build",
			"-ldflags", ldflags,
			"-o", "/blog",
			".",
		})
}

// Container builds a minimal container image with the blog server.
func (m *BlogBuilder) Container(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	git *dagger.Directory,
) (*dagger.Container, error) {
	// Get version info
	version, err := m.gitVersion(ctx, git)
	if err != nil {
		version = "dev"
	}
	commit, err := m.gitCommit(ctx, git)
	if err != nil {
		commit = "unknown"
	}

	// Build Hugo site
	publicDir := m.HugoBuild(source)

	// Run lint
	if _, err := m.Lint(ctx, source); err != nil {
		return nil, fmt.Errorf("lint failed: %w", err)
	}

	// Run tests
	if _, err := m.Test(ctx, source); err != nil {
		return nil, fmt.Errorf("tests failed: %w", err)
	}

	// Build binary
	binaryContainer := m.BuildBinary(source, publicDir, version, commit)
	binary := binaryContainer.File("/blog")

	// Create minimal runtime container
	return dag.Container().
		From("alpine:3.21").
		WithExec([]string{"apk", "add", "--no-cache", "ca-certificates", "tzdata"}).
		WithFile("/usr/local/bin/blog", binary).
		// Create non-root user
		WithExec([]string{"adduser", "-D", "-u", "10001", "blog"}).
		WithEnvVariable("TZ", "UTC").
		WithEnvVariable("PORT", "8080").
		WithEnvVariable("METRICS_PORT", "9101").
		WithEnvVariable("LOG_LEVEL", "info").
		WithExposedPort(8080).
		WithExposedPort(9101).
		WithUser("10001").
		WithEntrypoint([]string{"/usr/local/bin/blog"}), nil
}

// Publish builds and publishes the container image to a registry.
func (m *BlogBuilder) Publish(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	git *dagger.Directory,
	// Container registry address (e.g., ghcr.io/jwhumphries/blog)
	// +optional
	// +default="ghcr.io/jwhumphries/blog"
	registry string,
	// Registry username
	// +optional
	registryUser string,
	// Registry password (as a secret)
	// +optional
	registryPassword *dagger.Secret,
) (string, error) {
	if registry == "" {
		registry = "ghcr.io/jwhumphries/blog"
	}

	// Build container
	container, err := m.Container(ctx, source, git)
	if err != nil {
		return "", fmt.Errorf("container build failed: %w", err)
	}

	// Get version for tagging
	version, err := m.gitVersion(ctx, git)
	if err != nil {
		version = "dev"
	}

	// Authenticate if credentials provided
	if registryUser != "" && registryPassword != nil {
		container = container.WithRegistryAuth(registry, registryUser, registryPassword)
	}

	// Publish with version tag
	versionRef, err := container.Publish(ctx, fmt.Sprintf("%s:%s", registry, version))
	if err != nil {
		return "", fmt.Errorf("failed to publish version tag: %w", err)
	}

	// Also publish as latest
	latestRef, err := container.Publish(ctx, fmt.Sprintf("%s:latest", registry))
	if err != nil {
		return "", fmt.Errorf("failed to publish latest tag: %w", err)
	}

	return fmt.Sprintf("Published:\n  %s\n  %s", versionRef, latestRef), nil
}
