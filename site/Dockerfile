FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:2f9432b119a87428bd106582a9da0ecb0ebb9684faf9eee23c0aca7559cd64d1 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:7bb99a126eddeea8fbfaacedc5ef6708805411bdeaa01e668405a137fd02d048 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
