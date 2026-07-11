FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:19200b82b8b1de206cd9ddb3649297312001dda527f59ffbf7e43dd86efb6db1 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:f8671f2299e60154536c158bff8ce27f6eef4dddbbfc73bcce66263276ae0f80 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
