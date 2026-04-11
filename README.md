# docker-cloudflare-ddns
## How to use
[![GHCR](https://img.shields.io/badge/ghcr.io-chenwei791129%2Fcloudflare--ddns-blue)](https://ghcr.io/chenwei791129/cloudflare-ddns)
[![Docker Hub](https://img.shields.io/docker/pulls/awei/cloudflare-ddns.svg)](https://hub.docker.com/r/awei/cloudflare-ddns/)

[View on GHCR](https://ghcr.io/chenwei791129/cloudflare-ddns) | [View on Docker Hub](https://hub.docker.com/r/awei/cloudflare-ddns)

```shell
# Pull from GHCR (recommended)
$ docker run -d -e CF_TOKEN="<cloudflare-token>" -e CF_ZONE_ID="<cloudflare-zone-id>" -e CF_DOMAIN_NAME=<your.domain> ghcr.io/chenwei791129/cloudflare-ddns

# Pull from Docker Hub
$ docker run -d -e CF_TOKEN="<cloudflare-token>" -e CF_ZONE_ID="<cloudflare-zone-id>" -e CF_DOMAIN_NAME=<your.domain> awei/cloudflare-ddns
```

[How to get Cloudflare Zone ID and Token](https://github.com/chenwei791129/docker-cloudflare-ddns/wiki/How-to-get-Cloudflare-Zone-ID-and-Token%3F)


### Necessary Environment Variables
* `CF_TOKEN` your cloudflare token (string)
* `CF_ZONE_ID` your cloudflare zone id (string)
* `CF_DOMAIN_NAME` your domain full name, e.g. `blog.example.com` or `*.example.com` (string)

### Option Environment Variables
* `CF_PROXIED` CF proxied function (boolean, default: false)
* `TTL` TTL (integer, default: 1)
* `CHECK_URL` a url can return your ip (string, default: "http://whatismyip.akamai.com/") also can use http://ipv4.icanhazip.com, http://api.ipify.org, https://checkip.amazonaws.com/ ...
* `CHECK_INTERVAL` interval between IP checks using Go duration format, e.g. `30s`, `5m`, `1h` (string, default: "5m")

## License
The repository is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
