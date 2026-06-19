FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:e508a9dd017e3311fe652e4656feb122e5ae0bf6925d684d0b0d701ae581ae8f AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:89154b89cf899b1d26672a6ea911d88ed661720723c55e17cef3ad098ebdfada AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
