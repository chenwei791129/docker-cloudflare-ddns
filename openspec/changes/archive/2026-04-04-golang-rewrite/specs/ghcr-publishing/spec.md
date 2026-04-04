## MODIFIED Requirements

### Requirement: Release-driven image tagging

The system SHALL generate image tags from the release version using `docker/metadata-action`. The tag strategy SHALL produce the following tags for a release `v1.2.3`:

- `1.2.3` (full semver)
- `1.2` (minor level)
- `1` (major level)
- `latest` (default branch only)

The `alpine` tag SHALL NOT be produced, as the image is no longer based on Alpine Linux.

#### Scenario: Semver tags generated from release version

- **WHEN** the build workflow receives version `v1.2.3`
- **THEN** `docker/metadata-action` SHALL produce tags `1.2.3`, `1.2`, `1`, and `latest`

#### Scenario: Alpine tag removed

- **WHEN** image tags are generated
- **THEN** the tag list SHALL NOT include an `alpine` tag

#### Scenario: Tags applied to both registries

- **WHEN** image tags are generated
- **THEN** the same set of tags SHALL be applied to both GHCR and Docker Hub images
