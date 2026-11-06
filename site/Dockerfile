FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:1f46904fc7cebcae408869dd5da0d761da11cb5e3051441dcfedc86bb10eed43 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:a146e0639230bbc73358e6563240f0a6103466f0390b80a457c775068ec802a9 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
