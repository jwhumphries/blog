FROM ghcr.io/gohugoio/hugo:latest AS hugo

FROM hugo AS develop
CMD ["server", "--bind", "0.0.0.0", "--buildDrafts", "--disableFastRender"]

FROM hugo AS builder
CMD ["--gc", "--minify"]
