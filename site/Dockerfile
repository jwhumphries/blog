FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:68180dec78f611bce38eb04acba73bdf501c50d14cfba28ad89f0c7f6eeebe5e AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:4cc9f383086c0fba180efcf7132cf2f025969c462fedf5fcd2bf921915416797 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
