## Why

目前專案使用單一 workflow 在 push to `master` 時直接建置並推送 Docker image 至 Docker Hub，缺乏正式的版本管理流程。導入 `googleapis/release-please-action` 可自動化 semantic versioning 與 changelog 產生，同時將 image 遷移至 GHCR 作為主要 registry 以統一 GitHub 生態系整合。

## What Changes

- 新增 `release-please.yml` workflow，在 push to `master` 時由 release-please 自動建立 release PR 並管理版本
- 將現有 `docker-image.yml` 重新命名為 `build-image.yml`，改由 release event 觸發建置
- 新增 `version.txt` 檔案供 release-please `simple` type 追蹤版本
- Docker image 推送至 GHCR (`ghcr.io/chenwei791129/cloudflare-ddns`) 作為主要 registry
- 保留 Docker Hub (`awei/cloudflare-ddns`) 雙推，確保過渡期相容
- Image tag 策略改為 release-driven：`1.2.3`、`1.2`、`1`、`latest`、`alpine`
- 升級所有 GitHub Actions 至最新版本（checkout@v4, metadata@v5, setup-qemu@v3, setup-buildx@v3, build-push@v6, login@v3）

## Non-Goals

- 不變更 Dockerfile 內容或建置邏輯
- 不移除 Docker Hub 支援（僅新增 GHCR 作為主要 registry）
- 不變更多平台建置目標（維持 linux/amd64, linux/arm64）

## Capabilities

### New Capabilities

- `release-management`: 透過 release-please 自動化版本管理、changelog 產生與 GitHub Release 建立
- `ghcr-publishing`: 將 Docker image 發佈至 GitHub Container Registry，並保留 Docker Hub 雙推

### Modified Capabilities

(none)

## Impact

- Affected code:
  - `.github/workflows/docker-image.yml` → 重新命名為 `.github/workflows/build-image.yml` 並重構觸發邏輯
  - 新增 `.github/workflows/release-please.yml`
  - 新增 `version.txt`
  - `README.md`（更新 image 來源說明與 badge）
