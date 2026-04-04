## Why

目前 cloudflare-ddns 使用 shell script（bash + curl + jq）實作，執行環境依賴 Alpine Linux，需要安裝多個套件（bash、curl、jq），且使用 crond 排程。這導致 Docker image 攻擊面較大、體積較大（~30MB），且 shell script 難以測試與維護。

改用 Go 重構為靜態編譯 binary，搭配 scratch base image，可大幅縮小 image 體積（~5-8MB）、消除所有外部依賴、降低安全風險。

## What Changes

- **BREAKING**: 移除 `CRON_TIME` 環境變數，改為 `CHECK_INTERVAL`（使用 Go `time.Duration` 格式，如 `5m`、`30s`、`1h`）
- 以 Go 重寫整個 `cloudflare-ddns.sh` 邏輯：取得公網 IP → 比對快取 → 呼叫 Cloudflare API 更新 A record
- 使用 Go 標準庫 `net/http` 取代 curl、`encoding/json` 取代 jq
- 使用 `time.Ticker` + `time.ParseDuration` 實作內建排程，取代 crond
- Dockerfile 改為 multi-stage build：`golang:1.26-alpine` builder + `scratch` runtime
- release-please `release-type` 從 `simple` 改為 `go`
- 移除 `cloudflare-ddns.sh`
- CI/CD image tags 移除 `alpine` tag

## Non-Goals

- 不新增多網域支援（維持單一 A record 更新）
- 不支援 IPv6 / AAAA record
- 不引入設定檔（繼續使用環境變數）
- 不新增 Web UI 或 health check endpoint

## Capabilities

### New Capabilities

- `ddns-updater`: Go 實作的 DDNS 更新核心邏輯，包含公網 IP 偵測、IP 快取比對、Cloudflare API 互動、定時排程

### Modified Capabilities

- `release-management`: `release-type` 從 `simple` 改為 `go`，版本來源改為 Go module
- `ghcr-publishing`: Dockerfile 改為 multi-stage build，移除 Alpine 相關 tags，新增 scratch 相關 tags

## Impact

- 移除檔案：`cloudflare-ddns.sh`
- 新增檔案：`main.go`、`go.mod`、`go.sum`
- 修改檔案：`Dockerfile`、`.github/workflows/build-image.yml`、`.github/workflows/release-please.yml`
- 環境變數變更：`CRON_TIME` → `CHECK_INTERVAL`（breaking change）
- 依賴變更：移除 bash、curl、jq、crond；新增 Go 1.26 build toolchain
