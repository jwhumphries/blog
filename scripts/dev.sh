#!/bin/sh
set -e

echo "🔄 Syncing Hugo module dependencies..."
hugo mod npm pack

echo "📦 Installing npm dependencies..."
npm install

echo "🚀 Starting Hugo development server..."
hugo server \
    --bind 0.0.0.0 \
    --buildDrafts \
    --disableFastRender
