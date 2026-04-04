## 1. 建立專用 HTTP client

- [x] 1.1 在 `main.go` 建立全域 `http.Client`，設定 30 秒 timeout，取代所有 `http.DefaultClient` 和 `http.Get` 呼叫（對應 Public IP detection 與 Cloudflare DNS record update 的 timeout 需求）

## 2. 抽取共用 HTTP request helper

- [x] 2.1 [P] 在 `main.go` 建立 `newCFRequest(method, url string, body io.Reader) (*http.Request, error)` helper，統一設定 `Authorization: Bearer` header，僅在有 body 時設定 `Content-Type: application/json`

## 3. 改善 getRecordID — 使用 query filter（Cloudflare DNS record update）

- [x] 3.1 [P] 修改 `getRecordID` 的 URL 加入 `?name={domain}&type=A` query parameter，移除手動逐筆比對邏輯，改為直接取用 API 回傳的第一筆結果

## 4. 修正 updateRecord — 結構化 JSON 序列化（Cloudflare DNS record update）

- [x] 4.1 [P] 在 `main.go` 新增 `dnsUpdatePayload` struct，包含 `Type`、`Name`、`Content`、`TTL`、`Proxied` 欄位並加上 JSON tag
- [x] 4.2 修改 `updateRecord` 使用 `dnsUpdatePayload` struct + `json.Marshal` 取代 `fmt.Sprintf` 組裝 JSON body

## 5. 快取 recordID

- [x] 5.1 修改 `runUpdate` 新增 `cachedRecordID *string` 參數（或等效機制），首次查詢後快取 record ID，後續 update cycle 直接使用快取值（對應 Cached record ID reuse 與 First update cycle record ID query 場景）

## 6. 整合與驗證

- [x] 6.1 將 `getRecordID`、`updateRecord`、`getPublicIP` 全部改用步驟 1.1 的專用 HTTP client 與步驟 2.1 的 `newCFRequest` helper
- [x] 6.2 執行 `make build` 與 `make lint` 確認編譯通過且無 lint 錯誤
