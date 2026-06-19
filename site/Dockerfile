FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:f7a61f3a483693b20d19119c516bbf75af3e45edd7662f6bbd601e363adcecef AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:89154b89cf899b1d26672a6ea911d88ed661720723c55e17cef3ad098ebdfada AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
