ARG BASE_IMAGE=alpine
FROM ${BASE_IMAGE}:3.16.2

ENV CLOUDFLARE_TOKEN="" \
    CLOUDFLARE_ZONE_ID="" \
    CLOUDFLARE_DOMAIN_NAME="" \
    CLOUDFLARE_PROXIED=false \
    TTL=1 \
    CHECK_URL="http://whatismyip.akamai.com/" \
    CRON_TIME="*/5 * * * *"

WORKDIR /scripts

COPY cloudflare-ddns.sh /scripts

RUN apk add --update bash curl jq && \
    chmod +x /scripts/cloudflare-ddns.sh

CMD echo "${CRON_TIME} /scripts/cloudflare-ddns.sh" >> /etc/crontabs/root && crond -f -L 2
