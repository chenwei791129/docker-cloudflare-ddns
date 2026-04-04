## 1. Go 專案初始化

- [x] 1.1 初始化 Go module（`go mod init`），建立 `go.mod`
- [x] 1.2 建立 `main.go`，實作環境變數組態讀取與驗證（environment variable configuration）：讀取 `CLOUDFLARE_TOKEN`、`CLOUDFLARE_ZONE_ID`、`CLOUDFLARE_DOMAIN_NAME`、`CLOUDFLARE_PROXIED`、`TTL`、`CHECK_URL`、`CHECK_INTERVAL`，缺少必要變數時 log 錯誤並 exit

## 2. 核心 DDNS 邏輯

- [x] 2.1 實作 public IP detection：向 `CHECK_URL` 發送 GET 請求，驗證回應為合法 IPv4 格式，失敗時 log 錯誤並 skip 本次 cycle
- [x] 2.2 實作 IP change detection via in-memory cache：使用記憶體變數儲存上次 IP，比對決定是否需要更新（Go 專案結構：單一 main.go，IP 快取使用記憶體變數）
- [x] 2.3 實作 Cloudflare DNS record update：查詢 zone record list 取得 record ID，PUT 更新 A record，處理 API 錯誤回應（HTTP client 使用 Go 標準庫 net/http）
- [x] 2.4 實作 structured log output：所有 log 訊息加上 `[YYYY-MM-DD HH:MM:SS ±ZZZZ]` 時間戳前綴

## 3. 排程機制

- [x] 3.1 實作 periodic execution via CHECK_INTERVAL：使用 `time.ParseDuration` 解析 `CHECK_INTERVAL`，預設 `5m`，無效值時 exit（排程機制使用 time.Ticker）
- [x] 3.2 啟動時立即執行一次更新 cycle，之後按 interval 週期執行

## 4. Docker 與 CI/CD

- [x] 4.1 改寫 Dockerfile 為 multi-stage build：`golang:1.26-alpine` builder + `scratch` runtime，複製 CA certificates 確保 HTTPS connectivity（Dockerfile 使用 multi-stage build，static binary in scratch Docker image）
- [x] 4.2 移除 `cloudflare-ddns.sh`
- [x] 4.3 修改 `.github/workflows/build-image.yml`：移除 `build-args: BASE_IMAGE=alpine`，image tags 移除 alpine，不新增替代 tag（release-driven image tagging）
- [x] 4.4 修改 `.github/workflows/release-please.yml`：`release-type` 從 `simple` 改為 `go`（release-please 改用 go release type，automated version management via release-please）

## 5. 文件更新

- [x] 5.1 更新 `README.md`：將 `CRON_TIME` 替換為 `CHECK_INTERVAL`，說明 Go duration 格式，移除 Alpine 相關說明
