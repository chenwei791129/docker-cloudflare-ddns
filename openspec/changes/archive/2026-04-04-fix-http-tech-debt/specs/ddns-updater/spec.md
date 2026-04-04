## MODIFIED Requirements

### Requirement: Public IP detection

The system SHALL retrieve the current public IP address by sending an HTTP GET request to the URL specified in the `CHECK_URL` environment variable. The default value of `CHECK_URL` SHALL be `http://whatismyip.akamai.com/`. The response body SHALL be validated as a valid IPv4 address (four octets in dotted-decimal notation). If validation fails, the system SHALL log an error and skip the current update cycle. All HTTP requests SHALL use a dedicated HTTP client with a timeout of 30 seconds.

#### Scenario: Successful IP retrieval

- **WHEN** the system sends a GET request to `CHECK_URL`
- **THEN** the response body SHALL be parsed and validated as an IPv4 address

#### Scenario: Invalid IP response

- **WHEN** the response body is not a valid IPv4 address
- **THEN** the system SHALL log an error message and skip the current update cycle without exiting

#### Scenario: Network failure during IP check

- **WHEN** the HTTP request to `CHECK_URL` fails (timeout, DNS resolution failure, connection refused)
- **THEN** the system SHALL log an error message and skip the current update cycle without exiting

#### Scenario: HTTP request timeout

- **WHEN** the upstream server does not respond within 30 seconds
- **THEN** the HTTP client SHALL abort the request and the system SHALL treat it as a network failure

### Requirement: Cloudflare DNS record update

The system SHALL update the A record for the domain specified in `CLOUDFLARE_DOMAIN_NAME` via the Cloudflare API v4. The update process SHALL:

1. Query `GET /zones/{zone_id}/dns_records?name={domain}&type=A` to find the record ID of the A record matching `CLOUDFLARE_DOMAIN_NAME`
2. Update the record via `PUT /zones/{zone_id}/dns_records/{record_id}` with the new IP, TTL, and proxied settings

The record ID query result SHALL be cached in memory after the first successful lookup. Subsequent update cycles SHALL reuse the cached record ID without re-querying the API.

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
- **THEN** the system SHALL query the Cloudflare API with name and type filters to retrieve the record ID
