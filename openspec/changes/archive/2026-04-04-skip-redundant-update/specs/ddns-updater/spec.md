## MODIFIED Requirements

### Requirement: IP change detection via in-memory cache

The system SHALL store the last known public IP address in an in-memory variable. On each check cycle, the system SHALL compare the newly retrieved IP against the cached IP. If the IPs match, the system SHALL log that no update is required and skip the Cloudflare API call.

On the first check cycle, the system SHALL initialize the in-memory IP cache from the DNS record query response (the `content` field of the A record returned by Cloudflare API). This initialization SHALL occur as part of the existing record ID query — no additional API call SHALL be made.

#### Scenario: IP unchanged

- **WHEN** the retrieved IP matches the cached IP
- **THEN** the system SHALL log "no update required" and skip the Cloudflare API call

#### Scenario: IP changed

- **WHEN** the retrieved IP differs from the cached IP
- **THEN** the system SHALL proceed to update the Cloudflare DNS record

#### Scenario: First run — IP matches DNS record

- **WHEN** the system starts for the first time and the cache is empty
- **AND** the public IP matches the IP stored in the Cloudflare DNS A record
- **THEN** the system SHALL log "no update required" and skip the Cloudflare update API call

#### Scenario: First run — IP differs from DNS record

- **WHEN** the system starts for the first time and the cache is empty
- **AND** the public IP differs from the IP stored in the Cloudflare DNS A record
- **THEN** the system SHALL proceed to update the Cloudflare DNS record

### Requirement: Cloudflare DNS record update

The system SHALL update the A record for the domain specified in `CLOUDFLARE_DOMAIN_NAME` via the Cloudflare API v4. The update process SHALL:

1. Query `GET /zones/{zone_id}/dns_records?name={domain}&type=A` to find the record ID and current IP (`content` field) of the A record matching `CLOUDFLARE_DOMAIN_NAME`
2. Update the record via `PUT /zones/{zone_id}/dns_records/{record_id}` with the new IP, TTL, and proxied settings

The record ID query result SHALL be cached in memory after the first successful lookup. The `content` field from the query response SHALL be used to initialize the in-memory IP cache on the first check cycle. Subsequent update cycles SHALL reuse the cached record ID without re-querying the API.

Authentication SHALL use Bearer token from `CLOUDFLARE_TOKEN`. All Cloudflare API requests SHALL use a dedicated HTTP client with a timeout of 30 seconds. The JSON request body for the PUT request SHALL be constructed using structured serialization (not string formatting). On successful update, the system SHALL update the in-memory cache with the new IP.

#### Scenario: Successful DNS record update

- **WHEN** the IP has changed and the Cloudflare API returns `success: true` for both the query and update requests
- **THEN** the system SHALL update the in-memory IP cache and log the successful update with the new IP

#### Scenario: Record ID not found

- **WHEN** no A record matching `CLOUDFLARE_DOMAIN_NAME` is found in the zone
- **THEN** the system SHALL log an error indicating the record was not found and skip the update

#### Scenario: API authentication failure

- **WHEN** the Cloudflare API returns `success: false` on the record list query
- **THEN** the system SHALL log the error details from the API response and skip the update

#### Scenario: Update API call failure

- **WHEN** the Cloudflare API returns `success: false` on the PUT update request
- **THEN** the system SHALL log the error details and NOT update the in-memory cache

#### Scenario: Cached record ID reuse

- **WHEN** the record ID has been successfully retrieved in a previous update cycle
- **THEN** the system SHALL skip the record ID query and use the cached value directly

#### Scenario: First update cycle record ID query

- **WHEN** the system performs its first update cycle and no cached record ID exists
- **THEN** the system SHALL query the Cloudflare API with name and type filters to retrieve the record ID and current IP content
