FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:4036410be2d0c2c145f87c658a3a68bf8c282cd59bfac73f52f9cea969da3ee4 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:4cc9f383086c0fba180efcf7132cf2f025969c462fedf5fcd2bf921915416797 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
