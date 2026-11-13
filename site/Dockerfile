FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:1f46904fc7cebcae408869dd5da0d761da11cb5e3051441dcfedc86bb10eed43 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:90663d477d1fe08c85643d5752bbf1f4e6391e83f95f85e9942b702649d16cb2 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
