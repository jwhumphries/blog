FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:e508a9dd017e3311fe652e4656feb122e5ae0bf6925d684d0b0d701ae581ae8f AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:c3e613f59c53664d0bb9e7f2d4cb65e30b28a4381b4b415c02c4068e324438bf AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
