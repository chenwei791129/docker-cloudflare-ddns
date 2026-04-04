## Why

首次啟動時，in-memory cache 為空，系統一定會呼叫 Cloudflare update API，即使 DNS 上的 IP 與目前 public IP 完全相同。這造成不必要的 API 呼叫。現有的 `getRecordID` API 回應本身就包含 record 的 `content`（即當前 DNS IP），只是程式沒有解析這個欄位。

## What Changes

- 從 `getRecordID` 的 API 回應中一併解析 DNS record 的 `content` 欄位，用於初始化 in-memory IP cache
- 首次啟動時，如果 DNS record IP 與 public IP 相同，跳過 update，不再無條件觸發更新

## Non-Goals

- 不增加額外的 Cloudflare API 呼叫
- 不改變 `CHECK_INTERVAL` 或其他定時邏輯
- 不改變後續週期的 IP 比對行為（仍使用 in-memory cache）

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `ddns-updater`: 修改「IP change detection via in-memory cache」requirement — 首次啟動時從 DNS record query 回應初始化 cache，而非視為空值直接觸發更新。修改「Cloudflare DNS record update」requirement — `getRecordID` 同時回傳 record 的 content 欄位。

## Impact

- Affected specs: `ddns-updater`（修改 2 個 requirement）
- Affected code: `main.go`（`dnsRecord` struct、`getRecordID` 函式、`runUpdate` 函式）
