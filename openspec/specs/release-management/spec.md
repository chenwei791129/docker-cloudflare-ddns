### Requirement: Automated version management via release-please

The system SHALL use `googleapis/release-please-action@v4` with `release-type: simple` to manage versions automatically. A `version.txt` file in the repository root SHALL serve as the version source of truth.

#### Scenario: Release PR creation on push to master

- **WHEN** a commit is pushed to the `master` branch
- **THEN** release-please SHALL create or update a release PR with the bumped version and generated changelog

#### Scenario: GitHub Release creation on PR merge

- **WHEN** a release-please PR is merged
- **THEN** release-please SHALL create a GitHub Release with the version tag (e.g., `v1.2.3`)


<!-- @trace
source: release-please-and-ghcr
updated: 2026-04-04
code:
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - README.md
  - version.txt
  - .github/workflows/docker-image.yml
  - .github/workflows/release-please.yml
-->


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

### Requirement: Release-please workflow file

The system SHALL have a dedicated workflow file `.github/workflows/release-please.yml` that runs on push to `master` with `contents: write` and `pull-requests: write` permissions.

#### Scenario: Workflow triggers build on release creation

- **WHEN** release-please creates a new release
- **THEN** the workflow SHALL invoke the `build-image.yml` workflow via `workflow_call`, passing the release tag name as the `version` input

#### Scenario: No build triggered without release

- **WHEN** release-please updates an existing release PR without creating a release
- **THEN** the build workflow SHALL NOT be triggered


<!-- @trace
source: release-please-and-ghcr
updated: 2026-04-04
code:
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - README.md
  - version.txt
  - .github/workflows/docker-image.yml
  - .github/workflows/release-please.yml
-->

### Requirement: Version tracking file

The system SHALL maintain a `version.txt` file in the repository root containing the current version number. The initial version SHALL be `0.1.0`.

#### Scenario: Version file exists and is readable

- **WHEN** release-please runs
- **THEN** it SHALL read the current version from `version.txt` and bump it according to conventional commit messages

## Requirements

<!-- @trace
source: release-please-and-ghcr
updated: 2026-04-04
code:
  - .spectra.yaml
  - .github/workflows/build-image.yml
  - README.md
  - version.txt
  - .github/workflows/docker-image.yml
  - .github/workflows/release-please.yml
-->

### Requirement: Automated version management via release-please

The system SHALL use `googleapis/release-please-action@v4` with `release-type: go` to manage versions automatically. A `version.txt` file in the repository root SHALL serve as the version source of truth.

#### Scenario: Release PR creation on push to master

- **WHEN** a commit is pushed to the `master` branch
- **THEN** release-please SHALL create or update a release PR with the bumped version and generated changelog

#### Scenario: GitHub Release creation on PR merge

- **WHEN** a release-please PR is merged
- **THEN** release-please SHALL create a GitHub Release with the version tag (e.g., `v1.2.3`)

---
### Requirement: Release-please workflow file

The system SHALL have a dedicated workflow file `.github/workflows/release-please.yml` that runs on push to `master` with `contents: write` and `pull-requests: write` permissions.

#### Scenario: Workflow triggers build on release creation

- **WHEN** release-please creates a new release
- **THEN** the workflow SHALL invoke the `build-image.yml` workflow via `workflow_call`, passing the release tag name as the `version` input

#### Scenario: No build triggered without release

- **WHEN** release-please updates an existing release PR without creating a release
- **THEN** the build workflow SHALL NOT be triggered

---
### Requirement: Version tracking file

The system SHALL maintain a `version.txt` file in the repository root containing the current version number. The initial version SHALL be `0.1.0`.

#### Scenario: Version file exists and is readable

- **WHEN** release-please runs
- **THEN** it SHALL read the current version from `version.txt` and bump it according to conventional commit messages