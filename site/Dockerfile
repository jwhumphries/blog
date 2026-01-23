FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:a4fdf32e156f84f0221a77b2c5afc2448a6b143b088df5e2d3e3fa6ac31f4656 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:53dc48ef4d550835b0e54b0f6b41e22e5160e27065d0691b220a713218eb059d AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
