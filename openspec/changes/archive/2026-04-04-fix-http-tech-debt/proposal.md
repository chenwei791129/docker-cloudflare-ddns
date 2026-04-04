## Why

Code review 過程中發現 `main.go` 的 HTTP 相關程式碼存在 5 項技術債，包含無意義的 header 設定、手動組裝 JSON 導致潛在 malformed payload、無 timeout 設定可能永久阻塞、重複查詢可快取的資料、以及未善用 API query filter 造成不必要的網路傳輸。這些問題影響程式的正確性、健壯性與效率。

## What Changes

- 移除 `getRecordID` 中 GET request 無意義的 `Content-Type: application/json` header
- 抽取共用的 HTTP request helper，統一 Authorization header 設定邏輯
- `updateRecord` 改用 struct + `json.Marshal` 取代 `fmt.Sprintf` 組裝 JSON body，避免特殊字元造成 malformed JSON
- 建立專用 `http.Client` 並設定 timeout（取代 `http.DefaultClient`），避免上游服務無回應時永久阻塞
- 快取 `recordID`，首次查詢後重複使用，將每次 IP 變更的 HTTP call 從兩次減為一次
- `getRecordID` 改用 Cloudflare API 的 `?name=<domain>&type=A` query parameter 直接查詢目標 record，取代取得全部 records 再逐筆比對

## Non-Goals

- 不變更現有的環境變數設定介面或 config struct 的公開欄位
- 不引入外部 HTTP client library（如 resty），維持使用標準庫
- 不變更 log 格式或日誌行為
- 不新增單元測試（可作為後續獨立變更）

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `ddns-updater`：HTTP client 行為變更（timeout、request helper、query filter、recordID 快取），不影響外部可觀察的功能需求，屬於實作層級改善

## Impact

- 受影響程式碼：`main.go`（`getPublicIP`、`getRecordID`、`updateRecord`、`runUpdate`）
- 受影響 spec：`ddns-updater`（實作細節變更，需求層級無變動）
