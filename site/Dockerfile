FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:2f9432b119a87428bd106582a9da0ecb0ebb9684faf9eee23c0aca7559cd64d1 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:1fb8fe7f417abd5802d86747d297728aa48078d763159e7dff6559fe3f8bb8e7 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
