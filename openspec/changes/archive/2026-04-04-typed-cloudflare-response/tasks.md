## 1. 定義新型別

- [x] 1.1 [P] 新增 `cloudflareError` struct，包含 `Code int` 和 `Message string` 欄位，對應 JSON `code` 和 `message`
- [x] 1.2 將 `cloudflareResponse` 改為 generic struct `cloudflareResponse[T any]`，`Errors` 改為 `[]cloudflareError`，`Result` 改為 `T`

## 2. 更新呼叫端

- [x] 2.1 [P] 更新 `getRecordID`：使用 `cloudflareResponse[[]dnsRecord]` 解碼回應，移除二次 `json.Unmarshal`，直接存取 `cfResp.Result`
- [x] 2.2 [P] 更新 `updateRecord`：使用 `cloudflareResponse[dnsRecord]` 解碼回應

## 3. 改善錯誤輸出

- [x] 3.1 將 `getRecordID` 和 `updateRecord` 中的 error 格式化從 `json.Marshal(cfResp.Errors)` 改為遍歷 `[]cloudflareError` 產生結構化訊息（如 `"code NNNN: message"`）

## 4. 清理

- [x] 4.1 移除 `encoding/json` 以外不再需要的 `json.RawMessage` 相關 import（如果有的話）
- [x] 4.2 執行 `go vet` 和 `go build` 確認編譯通過
