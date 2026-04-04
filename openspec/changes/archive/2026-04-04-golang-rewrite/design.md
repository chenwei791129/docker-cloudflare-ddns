## Context

目前專案使用 shell script（`cloudflare-ddns.sh`）配合 Alpine Linux Docker image 執行 DDNS 更新。script 依賴 bash、curl、jq 三個外部工具，排程由 crond 管理。此架構導致 image 體積偏大（~30MB）、攻擊面廣（含完整 shell 環境），且 shell script 不易進行單元測試。

## Goals / Non-Goals

**Goals:**

- 以 Go 靜態編譯 binary 完全取代 shell script 及所有外部依賴
- 使用 scratch base image 將 Docker image 體積降至 ~5-8MB
- 消除容器內的 shell 環境，大幅降低安全攻擊面
- 保持環境變數介面向後相容（除 `CRON_TIME` → `CHECK_INTERVAL` 的 breaking change）
- CI/CD 流程對齊 Go 生態

**Non-Goals:**

- 不新增多網域或 AAAA record 支援
- 不引入設定檔，繼續使用環境變數
- 不新增 Web UI 或 health check endpoint
- 不為 Alpine 版本提供向後相容 image

## Decisions

### Go 專案結構：單一 main.go

所有邏輯放在 `main.go` 中，不拆分 package。理由：整個程式邏輯不超過 200 行，拆分 package 是過度抽象。
替代方案：拆成 `pkg/cloudflare`、`pkg/ip` 等 package — 對此規模的專案來說增加不必要的複雜度。

### HTTP client 使用 Go 標準庫 net/http

使用 `net/http` 發送 HTTP 請求，`encoding/json` 處理 JSON。不引入任何第三方 HTTP 或 JSON 套件。
替代方案：使用 resty 或 go-retryablehttp — 增加依賴卻沒有明顯收益，此程式的 HTTP 需求極為簡單。

### 排程機制使用 time.Ticker

使用 `time.Ticker` 配合 `time.ParseDuration` 解析 `CHECK_INTERVAL` 環境變數（預設 `5m`）。程式啟動時立即執行一次，之後按 interval 週期執行。
替代方案：引入 cron library（如 robfig/cron）支援 cron 表達式 — 過度設計，固定間隔已滿足需求。

### IP 快取使用記憶體變數

現有 script 將 IP 寫入 `/tmp/current_ip` 檔案。Go 版改為使用記憶體變數儲存上次 IP，因為程式是長駐 process，不需要跨 process 持久化。
替代方案：繼續寫檔案 — scratch image 沒有可寫檔案系統（除非 mount volume），且記憶體變數更簡單。

### Dockerfile 使用 multi-stage build

第一階段：`golang:1.26-alpine` 編譯靜態 binary（`CGO_ENABLED=0 GOOS=linux go build`）。
第二階段：`FROM scratch`，僅 COPY binary 和 CA certificates（`/etc/ssl/certs/ca-certificates.crt`）。
CA certificates 必須包含，因為程式需要 HTTPS 連線到 Cloudflare API。

### release-please 改用 go release type

`.github/workflows/release-please.yml` 中 `release-type` 從 `simple` 改為 `go`，版本管理對齊 Go module 生態。

### image tags 移除 alpine，不新增替代 tag

移除 `alpine` tag（因為不再基於 Alpine）。不新增 `scratch` tag — `latest` 足以代表預設版本。

## Risks / Trade-offs

- **[scratch 無法 exec 進容器 debug]** → 如需 debug，可臨時將 base image 改為 `gcr.io/distroless/static` 或在本地直接執行 binary。生產環境不需要 shell access。
- **[CHECK_INTERVAL breaking change]** → 在 README 和 CHANGELOG 中明確說明遷移方式：`CRON_TIME="*/5 * * * *"` 改為 `CHECK_INTERVAL=5m`。
- **[scratch 無 CA certificates]** → 在 build 階段從 Alpine 複製 CA certificates 到 scratch image，確保 HTTPS 連線正常。
- **[容器重啟後 IP 快取遺失]** → 可接受：重啟後首次執行會重新查詢並更新，Cloudflare API 冪等操作，不會造成問題。
