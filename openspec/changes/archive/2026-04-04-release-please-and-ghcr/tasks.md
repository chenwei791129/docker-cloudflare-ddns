## 1. Version tracking file 設定

- [x] 1.1 建立 `version.txt`，內容為 `0.1.0`（Version tracking file）

## 2. Release-please workflow 建立

- [x] 2.1 建立 `.github/workflows/release-please.yml`：使用 release-type: simple 搭配 version.txt，設定 `contents: write` 和 `pull-requests: write` permissions，實現 automated version management via release-please
- [x] 2.2 設定 release-please workflow 的 outputs，完成 workflow 拆分為 release-please.yml 和 build-image.yml：當 `release_created` 為 true 時透過 `workflow_call` 觸發 `build-image.yml` 並傳遞 `tag_name` 作為 version input（Release-please workflow file）

## 3. Build workflow 重構

- [x] 3.1 將 `.github/workflows/docker-image.yml` 重新命名為 `.github/workflows/build-image.yml`
- [x] 3.2 移除 push trigger，改為 `workflow_call` trigger 並接受 `version` input 參數（Build workflow as reusable workflow）
- [x] 3.3 升級所有 GitHub Actions 至最新版本：checkout@v4, metadata@v5, setup-qemu@v3, setup-buildx@v3, login@v3, build-push@v6（GitHub Actions upgraded to latest versions）

## 4. GHCR publishing 設定

- [x] 4.1 新增 GHCR login step，GHCR 登入使用 GITHUB_TOKEN，搭配 `docker/login-action@v3` 和 `ghcr.io` registry（Docker image published to GHCR）
- [x] 4.2 更新 `docker/metadata-action` 的 images 設定，加入 `ghcr.io/chenwei791129/cloudflare-ddns` 並保留 Docker Hub `awei/cloudflare-ddns`（Docker Hub dual-push preserved）
- [x] 4.3 設定 image tag 策略，實現 release-driven image tagging：使用 `type=semver` pattern 產生 `1.2.3`、`1.2`、`1` tag，加上 `latest` 和 `alpine` raw tag

## 5. README 更新

- [x] 5.1 更新 README.md 中的 Docker image 來源說明，加入 GHCR pull 指令與 badge
