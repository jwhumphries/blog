FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:19200b82b8b1de206cd9ddb3649297312001dda527f59ffbf7e43dd86efb6db1 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:608a19e34f86de36773503adbaab174fc28a6e338dc7904e03c70320b003a153 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
