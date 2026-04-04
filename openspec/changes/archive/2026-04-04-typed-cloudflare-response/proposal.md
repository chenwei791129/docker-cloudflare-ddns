## Why

`cloudflareResponse` 使用 `json.RawMessage` 處理 `Result` 和 `Errors` 欄位，導致需要二次反序列化（`json.Unmarshal`），且 error 輸出為原始 JSON dump 而非結構化訊息。此專案僅使用兩個固定的 Cloudflare API endpoint，不需要通用的 RawMessage 彈性。

## What Changes

- 新增 `cloudflareError` struct，解析 Cloudflare 錯誤回應的 `code` 和 `message` 欄位
- 將 `cloudflareResponse` 改為 generic struct `cloudflareResponse[T any]`，以 type parameter 取代 `json.RawMessage`
- 移除 `getRecordID` 中對 `Result` 的二次 `json.Unmarshal`
- 改善 error 輸出，從原始 JSON dump 改為結構化的 code + message 格式

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

（無）

## Impact

- 受影響程式碼：`main.go`（`cloudflareResponse` struct、`getRecordID` 函式、`updateRecord` 函式）
- 無 API 或對外行為變更，純內部重構
