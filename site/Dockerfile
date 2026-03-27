FROM ghcr.io/jwhumphries/tailwindcss:latest@sha256:7c9c8f28e81d8c932abed855a35a4c9124bcdd5cb3e018385aed3900f6b91107 AS tailwind
FROM ghcr.io/gohugoio/hugo:latest@sha256:3647208bc2456d9d85efcf0906e21882cacd4a0739d4fccc86bdfc81a9573546 AS hugo

COPY --from=tailwind /usr/local/bin/tailwindcss /usr/local/bin/

FROM hugo AS dev
EXPOSE 1313
COPY --chmod=755 scripts/dev.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
