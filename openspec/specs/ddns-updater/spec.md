# ddns-updater Specification

## Purpose

TBD - created by archiving change 'golang-rewrite'. Update Purpose after archive.

## Requirements

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


<!-- @trace
source: fix-http-tech-debt
updated: 2026-04-04
code:
  - main.go
  - tech-debt.md
-->

---
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


<!-- @trace
source: skip-redundant-update
updated: 2026-04-04
code:
  - main.go
-->

---
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


<!-- @trace
source: skip-redundant-update
updated: 2026-04-04
code:
  - main.go
-->

---
### Requirement: Periodic execution via CHECK_INTERVAL

The system SHALL execute the DDNS update cycle periodically based on the `CHECK_INTERVAL` environment variable. The value SHALL be parsed using Go's `time.ParseDuration` format (e.g., `5m`, `30s`, `1h30m`). The default value SHALL be `5m`. The system SHALL execute one update cycle immediately on startup, then repeat at the configured interval.

#### Scenario: Default interval

- **WHEN** `CHECK_INTERVAL` is not set
- **THEN** the system SHALL use a default interval of `5m` (5 minutes)

#### Scenario: Custom interval

- **WHEN** `CHECK_INTERVAL` is set to `30s`
- **THEN** the system SHALL execute the update cycle every 30 seconds

#### Scenario: Invalid interval format

- **WHEN** `CHECK_INTERVAL` contains an unparseable value (e.g., `abc`, `*/5 * * * *`)
- **THEN** the system SHALL log an error and exit with a non-zero exit code

#### Scenario: Immediate execution on startup

- **WHEN** the system starts
- **THEN** the system SHALL execute one update cycle immediately before starting the periodic timer


<!-- @trace
source: golang-rewrite
updated: 2026-04-04
code:
  - .github/workflows/release-please.yml
  - go.mod
  - Dockerfile
  - cloudflare-ddns.sh
  - Makefile
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - version.txt
  - README.md
  - main.go
-->

---
### Requirement: Environment variable configuration

The system SHALL read configuration from the following environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `CLOUDFLARE_TOKEN` | Yes | — | Cloudflare API Bearer token |
| `CLOUDFLARE_ZONE_ID` | Yes | — | Cloudflare Zone ID |
| `CLOUDFLARE_DOMAIN_NAME` | Yes | — | Full domain name for the A record |
| `CLOUDFLARE_PROXIED` | No | `false` | Enable Cloudflare proxy |
| `TTL` | No | `1` | DNS record TTL (1 = automatic) |
| `CHECK_URL` | No | `http://whatismyip.akamai.com/` | URL that returns public IP |
| `CHECK_INTERVAL` | No | `5m` | Update check interval (Go duration format) |

The system SHALL validate that all required environment variables are set on startup. If any required variable is missing, the system SHALL log an error and exit with a non-zero exit code.

#### Scenario: All required variables set

- **WHEN** `CLOUDFLARE_TOKEN`, `CLOUDFLARE_ZONE_ID`, and `CLOUDFLARE_DOMAIN_NAME` are all set
- **THEN** the system SHALL start normally

#### Scenario: Missing required variable

- **WHEN** any of `CLOUDFLARE_TOKEN`, `CLOUDFLARE_ZONE_ID`, or `CLOUDFLARE_DOMAIN_NAME` is not set
- **THEN** the system SHALL log which variable is missing and exit with a non-zero exit code


<!-- @trace
source: golang-rewrite
updated: 2026-04-04
code:
  - .github/workflows/release-please.yml
  - go.mod
  - Dockerfile
  - cloudflare-ddns.sh
  - Makefile
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - version.txt
  - README.md
  - main.go
-->

---
### Requirement: Structured log output

The system SHALL output log messages to stdout with a timestamp prefix in the format `[YYYY-MM-DD HH:MM:SS ±ZZZZ]`. This format SHALL match the existing shell script's log format for consistency.

#### Scenario: Log format

- **WHEN** the system logs any message
- **THEN** the log line SHALL be prefixed with `[YYYY-MM-DD HH:MM:SS ±ZZZZ]`


<!-- @trace
source: golang-rewrite
updated: 2026-04-04
code:
  - .github/workflows/release-please.yml
  - go.mod
  - Dockerfile
  - cloudflare-ddns.sh
  - Makefile
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - version.txt
  - README.md
  - main.go
-->

---
### Requirement: Static binary in scratch Docker image

The system SHALL be compiled as a statically-linked Go binary (`CGO_ENABLED=0`). The Dockerfile SHALL use a multi-stage build: `golang:1.27-alpine` as the builder stage and `scratch` as the runtime stage. The runtime image SHALL include only the compiled binary and CA certificates (`/etc/ssl/certs/ca-certificates.crt` copied from the builder stage).

#### Scenario: Minimal runtime image

- **WHEN** the Docker image is built
- **THEN** the runtime stage SHALL contain only the Go binary and CA certificates, with no shell, package manager, or other utilities

#### Scenario: HTTPS connectivity

- **WHEN** the Go binary makes HTTPS requests to the Cloudflare API
- **THEN** TLS certificate verification SHALL succeed using the bundled CA certificates

#### Scenario: Multi-platform build

- **WHEN** the Docker image is built via CI
- **THEN** the image SHALL be built for both `linux/amd64` and `linux/arm64` platforms

<!-- @trace
source: golang-rewrite
updated: 2026-04-04
code:
  - .github/workflows/release-please.yml
  - go.mod
  - Dockerfile
  - cloudflare-ddns.sh
  - Makefile
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - version.txt
  - README.md
  - main.go
-->