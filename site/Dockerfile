FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:2f9432b119a87428bd106582a9da0ecb0ebb9684faf9eee23c0aca7559cd64d1 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:9f3cccb54b48e83a5468cd44f0372b10834b6d8418ef692d9821bb1314761829 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
