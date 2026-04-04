## 1. 擴展 DNS record 資料結構

- [x] 1.1 [P] 在 `dnsRecord` struct 新增 `Content string` 欄位（JSON tag `"content"`），使 Cloudflare DNS record update API 回應的 IP 能被解析

## 2. 修改 record ID 查詢以回傳 DNS record IP

- [x] 2.1 將 `getRecordID` 函式改為同時回傳 record ID 與 record content（當前 DNS IP），不增加額外 API 呼叫。對應 spec: Cloudflare DNS record update — First update cycle record ID query

## 3. 初始化 in-memory IP cache

- [x] 3.1 在 `runUpdate` 中，首次查詢 record ID 時，將回傳的 DNS record content 用來初始化 `state.ip`，實現 IP change detection via in-memory cache 的首次啟動行為

## 4. 驗證

- [x] 4.1 [P] 執行 `make build` 確認編譯通過
- [x] 4.2 [P] 執行 `make run` 驗證首次啟動時若 IP 未變則不觸發 update
