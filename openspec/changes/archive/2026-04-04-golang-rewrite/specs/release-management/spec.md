## MODIFIED Requirements

### Requirement: Automated version management via release-please

The system SHALL use `googleapis/release-please-action@v4` with `release-type: go` to manage versions automatically. A `version.txt` file in the repository root SHALL serve as the version source of truth.

#### Scenario: Release PR creation on push to master

- **WHEN** a commit is pushed to the `master` branch
- **THEN** release-please SHALL create or update a release PR with the bumped version and generated changelog

#### Scenario: GitHub Release creation on PR merge

- **WHEN** a release-please PR is merged
- **THEN** release-please SHALL create a GitHub Release with the version tag (e.g., `v1.2.3`)
