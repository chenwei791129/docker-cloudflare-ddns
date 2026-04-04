## Context

目前專案僅有一個 `.github/workflows/docker-image.yml`，在每次 push to `master` 時自動建置多平台 Docker image 並推送至 Docker Hub。沒有正式的版本管理機制，tag 依賴 semver pattern（需手動打 git tag）和 SHA-based 自動 tag。

參考專案 `chenwei791129/launchpal` 已成功導入 `googleapis/release-please-action@v4`，採用 release-please 產生 release PR → merge 後建立 GitHub Release → 觸發後續建置的模式。

## Goals / Non-Goals

**Goals:**

- 自動化版本管理與 changelog 產生
- Docker image 發佈至 GHCR 作為主要 registry
- 維持 Docker Hub 雙推確保過渡期相容
- 升級 GitHub Actions 至最新版本

**Non-Goals:**

- 不變更 Dockerfile 內容或建置邏輯
- 不移除 Docker Hub 支援
- 不變更多平台建置目標

## Decisions

### 使用 release-type: simple 搭配 version.txt

release-please 支援多種 release-type（go, node, python, simple 等）。本專案是純 shell script + Dockerfile，無 `package.json` 或 `go.mod`，`simple` type 透過 `version.txt` 追蹤版本是最直接的選擇。

**替代方案**：
- `node`：需要建立 `package.json`，對純 shell 專案不合理
- `go`：需要 `go.mod`，同樣不適用

### Workflow 拆分為 release-please.yml 和 build-image.yml

將 release 管理與 image 建置分離為兩個 workflow：

```
release-please.yml (push to master)
  └── 建立/更新 release PR
  └── merge 後建立 GitHub Release
       └── 觸發 build-image.yml (workflow_call)
            └── 建置並推送至 GHCR + Docker Hub
```

`build-image.yml` 同時支援 `workflow_call`（由 release-please 觸發）以接收 version 參數。

**替代方案**：
- 單一 workflow 處理所有邏輯：職責耦合，難以單獨測試建置流程

### GHCR 登入使用 GITHUB_TOKEN

GHCR 可直接使用 `${{ secrets.GITHUB_TOKEN }}` 認證，無需額外設定 PAT。Docker Hub 維持使用現有的 `DOCKER_USER` / `DOCKER_TOKEN` secrets。

### Image tag 策略

release-please 產生的 tag（如 `v1.2.3`）作為版本來源：

| Tag | 說明 |
|-----|------|
| `1.2.3` | 完整版本號 |
| `1.2` | Minor 層級 |
| `1` | Major 層級 |
| `latest` | 最新穩定版 |
| `alpine` | 維持現有慣例 |

使用 `docker/metadata-action` 的 `type=semver` pattern 自動產生。

## Risks / Trade-offs

- **初始版本號選擇**：`version.txt` 需要設定起始版本。建議從 `0.1.0` 開始，讓 release-please 在第一次 merge 時自動 bump。如果已有既有使用者期望特定版本號，可調整起始值。→ 可在 `version.txt` 中設定任意起始版本。
- **Docker Hub 雙推增加建置時間**：每次 release 需推送至兩個 registry。→ 影響微小，multi-platform build 本身才是主要耗時。
- **現有 SHA-based tag 使用者**：遷移後不再產生 `YYYYMMDD-sha` 格式的 tag。→ semver tag 是更標準的做法，影響有限。
