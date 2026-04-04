# Changelog

## [2.0.0](https://github.com/chenwei791129/docker-cloudflare-ddns/compare/v1.0.0...v2.0.0) (2026-04-04)


### ⚠ BREAKING CHANGES

* Environment variables renamed from CLOUDFLARE_* to CF_* (CF_TOKEN, CF_ZONE_ID, CF_DOMAIN_NAME, CF_PROXIED). CRON_TIME replaced by CHECK_INTERVAL (Go duration format, e.g. 5m, 30s, 1h).

### Features

* rewrite cloudflare-ddns in Go with scratch Docker image ([14010d3](https://github.com/chenwei791129/docker-cloudflare-ddns/commit/14010d3ec9401d2c96af0237de0c7cf98c7ca560))
* skip redundant DNS update when IP matches existing record ([dd081d7](https://github.com/chenwei791129/docker-cloudflare-ddns/commit/dd081d7b9fa738bcd2a6be2786cd0b7663721c04))

## 1.0.0 (2026-04-04)


### Features

* add release-please action and migrate Docker image to GHCR ([2a6e032](https://github.com/chenwei791129/docker-cloudflare-ddns/commit/2a6e032970ea20dc72f795f2fa4654fc6dbb4f29))


### Bug Fixes

* add packages:write permission to release-please workflow ([b04b08e](https://github.com/chenwei791129/docker-cloudflare-ddns/commit/b04b08eae3a79c7541bf6abe2cd8e908c6b6ce6f))
