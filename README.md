# docker-cloudflare-ddns
## How to use
[![This image on DockerHub](https://img.shields.io/docker/pulls/awei/cloudflare-ddns.svg)](https://hub.docker.com/r/awei/cloudflare-ddns/)

[View on Docker Hub](https://hub.docker.com/r/awei/cloudflare-ddns)

```shell
$ docker run -d -e CLOUDFLARE_TOKEN=<godaddy-api-key> -e CLOUDFLARE_ZONE_ID=<godaddy-api-secret> -e CLOUDFLARE_DOMAIN_NAME=<your.domain> awei/cloudflare-ddns
```

### Necessary Environment Variables
* `CLOUDFLARE_TOKEN` your cloudflare token (string)
* `CLOUDFLARE_ZONE_ID` your cloudflare zone id (string)
* `CLOUDFLARE_DOMAIN_NAME` your domain full name (string)

### Option Environment Variables
* `CLOUDFLARE_PROXIED` CF proxied function (boolean, default: false)
* `TTL` TTL (integer,default: 1)
* `CHECK_URL` a url can return your ip (string,default: "http://whatismyip.akamai.com/")
* `CRON_TIME` crontab job time (default: "*/5 * * * *")

## License
The repository is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
