# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A personal blog that combines Hugo static site generation with a Go web server for deployment to fly.io. The site deploys to `https://blog.johnhumphries.dev/`.

### Architecture

The project has two main components:

1. **Hugo Site** (`site/`) - Static site generator using the Shiloh theme
2. **Go Server** (root) - Lightweight HTTP server that embeds and serves the built Hugo site with Brotli pre-compression

## Project Structure

```
blog/
├── site/                      # Hugo project
│   ├── archetypes/
│   ├── assets/
│   ├── config/                # Hugo configuration
│   ├── content/               # Blog posts and pages
│   ├── static/
│   ├── go.mod                 # Hugo modules
│   └── go.sum
├── internal/
│   ├── server/                # HTTP server and routing
│   ├── compress/              # Brotli compression utilities
│   └── metrics/               # Prometheus metrics
├── version/                   # Build version metadata
├── .dagger/                   # Dagger pipeline
├── main.go                    # Server entry point
├── go.mod                     # Server dependencies
├── Taskfile.yml
├── fly.toml                   # fly.io configuration
└── dagger.json
```

## Build & Development Commands

### Local Development
```bash
task dev          # Start Hugo dev server with live reload
task build        # Build Hugo site via Dagger
task pack         # Build complete server binary via Dagger
task lint         # Run golangci-lint on Go code
task test         # Run Go tests
task test-prod    # Build and run production server locally
task clean        # Clean all generated files
```

### Dagger Functions
```bash
dagger call hugo-build --source .              # Build Hugo site only
dagger call lint --source .                    # Run linter
dagger call test --source .                    # Run tests
dagger call pack --source . --git .git         # Full build pipeline
dagger call container --source . --git .git   # Build container image
dagger call publish --source . --git .git     # Publish to GHCR
```

### Container Tasks
```bash
task container      # Build container image
task container-run  # Build and run container locally
task publish        # Publish to ghcr.io (requires GITHUB_TOKEN)
```

## Configuration

### Hugo Configuration
Hugo uses split configuration in `site/config/_default/`:
- `hugo.toml` - Core Hugo settings
- `module.toml` - Theme import and asset mounts
- `params.toml` - Site parameters and author info
- `menus.toml` - Navigation structure
- `markup.toml` - Markdown rendering options
- `build.toml` - Build stats for CSS purging

### Server Configuration (Environment Variables)
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP server port |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `METRICS_PORT` | 9101 | Prometheus metrics port (Stage 2) |

### Theme
The Shiloh theme (`github.com/jwhumphries/shiloh`) is imported via Hugo modules.

## Content Structure

- `site/content/posts/` - Blog posts (bundle format with `index.md`)
- `site/content/about/` - About page
- `site/content/topics/` - Taxonomy listing page

### Post Front Matter
```yaml
date: 2025-01-01
draft: false
author: John Humphries
title: Post Title
description: SEO description
topics: [topic1, topic2]
subjects: [subject1]
```

## Server Features

- **Pre-compressed responses**: Files are Brotli-compressed at startup for zero-CPU serving
- **Embedded filesystem**: Hugo output is embedded in the binary via `//go:embed`
- **Graceful shutdown**: Handles SIGINT/SIGTERM for clean container stops
- **Health endpoint**: `GET /health` for container orchestration
- **Prometheus metrics**: Exposed on `METRICS_PORT` (default 9101)

## Prometheus Metrics

The server exposes metrics on a separate port (default 9101) at `/metrics`.

### Available Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests by method, path, status |
| `http_request_duration_seconds` | Histogram | Request latency by method, path |
| `http_response_size_bytes` | Histogram | Response size by method, path |
| `http_compression_ratio` | Histogram | Brotli compression ratio by path |

### fly.io Integration

Add to `fly.toml` to enable fly.io metrics scraping:

```toml
[[metrics]]
port = 9101
path = "/metrics"
```

Metrics will be available at [fly-metrics.net](https://fly-metrics.net) after deployment.

## CI/CD

Dagger pipeline orchestrates:
1. Hugo site build (with Tailwind CSS)
2. Go linting (golangci-lint)
3. Go tests
4. Binary compilation with version injection
5. Container image creation (alpine-based, ~15MB)
6. Publishing to GHCR

## Deployment

### Container Image

The container image is based on Alpine Linux with:
- Non-root user (UID 10001)
- CA certificates and timezone data
- Exposed ports: 8080 (HTTP), 9101 (metrics)

### fly.io Deployment

Configuration is in `fly.toml`. Deploy with:

```bash
# First time setup
fly launch --no-deploy
fly secrets set GITHUB_TOKEN=<token>  # If using private GHCR image

# Deploy
fly deploy
```

### Manual Deployment

```bash
# Build and publish image
task publish

# Or deploy directly with Dagger
dagger call publish --source . --git .git \
  --registry-user jwhumphries \
  --registry-password env:GITHUB_TOKEN
```
