# docker-cloudflare-ddns
## How to use
[![GHCR](https://img.shields.io/badge/ghcr.io-chenwei791129%2Fcloudflare--ddns-blue)](https://ghcr.io/chenwei791129/cloudflare-ddns)
[![Docker Hub](https://img.shields.io/docker/pulls/awei/cloudflare-ddns.svg)](https://hub.docker.com/r/awei/cloudflare-ddns/)

[View on GHCR](https://ghcr.io/chenwei791129/cloudflare-ddns) | [View on Docker Hub](https://hub.docker.com/r/awei/cloudflare-ddns)

```shell
# Pull from GHCR (recommended)
$ docker run -d -e CLOUDFLARE_TOKEN="<cloudflare-token>" -e CLOUDFLARE_ZONE_ID="<cloudflare-zone-id>" -e CLOUDFLARE_DOMAIN_NAME=<your.domain> ghcr.io/chenwei791129/cloudflare-ddns

# Pull from Docker Hub
$ docker run -d -e CLOUDFLARE_TOKEN="<cloudflare-token>" -e CLOUDFLARE_ZONE_ID="<cloudflare-zone-id>" -e CLOUDFLARE_DOMAIN_NAME=<your.domain> awei/cloudflare-ddns
```

[How to get Cloudflare Zone ID and Token](https://github.com/chenwei791129/docker-cloudflare-ddns/wiki/How-to-get-Cloudflare-Zone-ID-and-Token%3F)


### Necessary Environment Variables
* `CLOUDFLARE_TOKEN` your cloudflare token (string)
* `CLOUDFLARE_ZONE_ID` your cloudflare zone id (string)
* `CLOUDFLARE_DOMAIN_NAME` your domain full name (string)

### Option Environment Variables
* `CLOUDFLARE_PROXIED` CF proxied function (boolean, default: false)
* `TTL` TTL (integer,default: 1)
* `CHECK_URL` a url can return your ip (string,default: "http://whatismyip.akamai.com/") also can use http://ipv4.icanhazip.com 、 http://api.ipify.org ...
* `CRON_TIME` crontab job time (default: "*/5 * * * *")

## License
The repository is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
