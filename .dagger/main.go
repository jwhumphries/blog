package main

import (
	"context"
	"dagger/blog-build/internal/dagger"
)

type BlogBuild struct{}

// Build recreates the steps from the “builder” target in the Dockerfile.
func (m *BlogBuild) Build(
	ctx context.Context,
	// The source directory containing the Hugo project.
	source *dagger.Directory,
) (*dagger.Directory, error) {
	// 1. Get Tailwind CSS binary
	// FROM ghcr.io/jwhumphries/tailwindcss:latest... AS tailwind
	tailwind := dagger.Dag().Container().
		From("ghcr.io/jwhumphries/tailwindcss:latest@sha256:a4fdf32e156f84f0221a77b2c5afc2448a6b143b088df5e2d3e3fa6ac31f4656").
		File("/usr/local/bin/tailwindcss")

	// 2. Setup Hugo container
	// FROM ghcr.io/gohugoio/hugo:latest AS hugo
	// COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/
	builder := dagger.Dag().Container().
		From("ghcr.io/gohugoio/hugo:latest").
		WithExec([]string{"apk", "add", "--no-cache", "nodejs", "npm"}).
		WithFile("/usr/local/bin/tailwindcss", tailwind)

	// 3. Build steps
	// FROM hugo AS builder
	// COPY . .
	// RUN hugo mod npm pack
	// RUN npm install
	// RUN hugo --gc --minify (implied from CMD)
	builder = builder.
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"hugo", "mod", "npm", "pack"}).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"hugo", "--gc", "--minify"})

	// Return the "public/" directory
	return builder.Directory("public"), nil
}
