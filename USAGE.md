# Usage Guide

This guide explains how to set up and deploy the blog to fly.io.

## Prerequisites

- [Dagger](https://docs.dagger.io/install) (v0.19.10+)
- [Task](https://taskfile.dev/installation/) (for running tasks)
- [Fly CLI](https://fly.io/docs/hands-on/install-flyctl/)
- Docker (for local container testing)
- A GitHub account (for GHCR)
- A Fly.io account

## Local Development

### Content Authoring (Hugo)

For writing and previewing content:

```bash
task dev
```

This starts Hugo's development server at http://localhost:1313 with live reload.

### Testing the Production Server

Build and run the Go server locally:

```bash
task test-prod
```

Or build and run the container:

```bash
task container-run
```

Both expose:
- http://localhost:8080 - Main site
- http://localhost:9101/metrics - Prometheus metrics

## Initial Fly.io Setup

### 1. Create the Fly.io App

```bash
# Login to Fly.io
fly auth login

# Create the app (don't deploy yet)
fly apps create blog-johnhumphries
```

> **Note**: Update `fly.toml` if you use a different app name:
> ```toml
> app = "your-app-name"
> ```

### 2. Choose a Region

The default region is `ord` (Chicago). To change it, update `fly.toml`:

```toml
primary_region = "sea"  # Seattle, or your preferred region
```

Available regions: https://fly.io/docs/reference/regions/

### 3. Set Up Secrets

The GitHub Actions workflow needs a Fly.io API token:

```bash
# Generate a token
fly tokens create deploy -x 999999h

# Add it to GitHub repository secrets
# Go to: Settings → Secrets and variables → Actions → New repository secret
# Name: FLY_API_TOKEN
# Value: <the token from above>
```

## GitHub Actions Setup

The workflow at `.github/workflows/publish.yaml` automatically:
1. Builds the container using Dagger
2. Publishes to GHCR (`ghcr.io/jwhumphries/blog`)
3. Deploys to Fly.io (on pushes to `main`)

### Required GitHub Settings

1. **Enable GHCR**: Go to your GitHub profile → Packages → Enable improved container support

2. **Package Visibility**: After the first publish, go to the package settings and make it public (or configure Fly.io to authenticate with GHCR)

3. **Repository Secrets**: Add the `FLY_API_TOKEN` secret as described above

## Manual Deployment

### Build and Publish Container

```bash
# Set your GitHub token
export GITHUB_TOKEN=$(gh auth token)

# Publish to GHCR
task publish
```

### Deploy to Fly.io

```bash
fly deploy
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP server port |
| `METRICS_PORT` | 9101 | Prometheus metrics port |
| `LOG_LEVEL` | info | Log level (debug, info, warn, error) |

Set environment variables in `fly.toml` or via the Fly CLI:

```bash
fly secrets set LOG_LEVEL=debug
```

### Scaling

The default configuration uses minimal resources:
- 256MB RAM
- Shared CPU
- Auto-stop when idle (cost-effective)

Adjust in `fly.toml`:

```toml
[[vm]]
  memory = "512mb"
  cpu_kind = "shared"
  cpus = 1
```

### Custom Domain

```bash
# Add a custom domain
fly certs create blog.example.com

# Update DNS to point to your Fly.io app
# CNAME: blog.example.com → blog-johnhumphries.fly.dev
```

## Monitoring

### Prometheus Metrics

Fly.io automatically scrapes metrics from port 9101. View them at:
- https://fly-metrics.net (requires Fly.io login)

Available metrics:
- `http_requests_total` - Request counts by method, path, status
- `http_request_duration_seconds` - Latency histograms
- `http_response_size_bytes` - Response size histograms

### Logs

```bash
# Stream logs
fly logs

# View recent logs
fly logs --no-tail
```

### Health Checks

The `/health` endpoint returns `ok` when the server is healthy. Fly.io uses this for:
- Container health monitoring
- Zero-downtime deployments

## Troubleshooting

### Container Won't Start

Check logs:
```bash
fly logs
```

Common issues:
- Missing environment variables
- Port conflicts
- Out of memory

### Deployment Fails

1. Check GitHub Actions logs
2. Verify `FLY_API_TOKEN` secret is set
3. Ensure GHCR package is accessible

### Metrics Not Appearing

1. Verify the `[[metrics]]` section in `fly.toml`
2. Check that port 9101 is exposed
3. Wait a few minutes for initial scrape

## File Changes Checklist

When customizing this setup, you may need to update:

| File | What to Change |
|------|----------------|
| `fly.toml` | App name, region, resources |
| `.github/workflows/publish.yaml` | Image name if not using default |
| `site/config/_default/hugo.toml` | Site URL, title |
| `site/config/_default/params.toml` | Author info, description |
