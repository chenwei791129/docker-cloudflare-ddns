## ADDED Requirements

### Requirement: Docker image published to GHCR

The system SHALL publish Docker images to `ghcr.io/chenwei791129/cloudflare-ddns` as the primary registry. Authentication SHALL use the built-in `GITHUB_TOKEN`.

#### Scenario: Successful GHCR push on release

- **WHEN** a new GitHub Release is created
- **THEN** the build workflow SHALL build multi-platform images (linux/amd64, linux/arm64) and push them to GHCR

#### Scenario: GHCR authentication

- **WHEN** the build workflow runs
- **THEN** it SHALL authenticate to `ghcr.io` using `${{ secrets.GITHUB_TOKEN }}` via `docker/login-action`

### Requirement: Docker Hub dual-push preserved

The system SHALL continue publishing Docker images to Docker Hub (`awei/cloudflare-ddns`) alongside GHCR to maintain backward compatibility for existing users.

#### Scenario: Dual registry push

- **WHEN** a new release triggers the build workflow
- **THEN** the built images SHALL be pushed to both GHCR and Docker Hub in a single build step

#### Scenario: Docker Hub authentication

- **WHEN** the build workflow runs
- **THEN** it SHALL authenticate to Docker Hub using `secrets.DOCKER_USER` and `secrets.DOCKER_TOKEN`

### Requirement: Release-driven image tagging

The system SHALL generate image tags from the release version using `docker/metadata-action`. The tag strategy SHALL produce the following tags for a release `v1.2.3`:

- `1.2.3` (full semver)
- `1.2` (minor level)
- `1` (major level)
- `latest` (default branch only)
- `alpine` (default branch only)

#### Scenario: Semver tags generated from release version

- **WHEN** the build workflow receives version `v1.2.3`
- **THEN** `docker/metadata-action` SHALL produce tags `1.2.3`, `1.2`, `1`, `latest`, and `alpine`

#### Scenario: Tags applied to both registries

- **WHEN** image tags are generated
- **THEN** the same set of tags SHALL be applied to both GHCR and Docker Hub images

### Requirement: Build workflow as reusable workflow

The build workflow `.github/workflows/build-image.yml` SHALL support `workflow_call` trigger with a `version` input parameter. It SHALL NOT trigger on direct push to `master`.

#### Scenario: Called by release-please workflow

- **WHEN** `release-please.yml` invokes `build-image.yml` with `version: v1.2.3`
- **THEN** the build workflow SHALL use the provided version for image tagging

### Requirement: GitHub Actions upgraded to latest versions

The build workflow SHALL use the following action versions: `actions/checkout@v4`, `docker/metadata-action@v5`, `docker/setup-qemu-action@v3`, `docker/setup-buildx-action@v3`, `docker/login-action@v3`, `docker/build-push-action@v6`.

#### Scenario: All actions at latest major versions

- **WHEN** the workflow files are committed
- **THEN** all referenced GitHub Actions SHALL use the versions specified above
