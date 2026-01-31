# Blastshield OpenAPI Specifications

This directory contains OpenAPI specifications for different versions of the Blastshield API.

## Current Versions

- *1.12.2.json - Current GA release version
- latest.json - Symlink to the most recent spec


### 1. Update the latest symlink

```bash
cd openapi
ln -sf v1.14.json latest.json
```

### 1. Generate and test the provider

```bash
make generate OPENAPI_VERSION=v1.14
make build
make test
```

### 1. Tag the release

```bash
git add openapi/v1.13.json
git commit -m "feat: add support for Blastshield API v1.13"
git tag v1.13.0
git push origin v1.13.0
```

## Naming Convention

- Use the API version from Blastshield: `v1.12`, `v1.13`, etc.
- Provider version matches API version: Provider v1.13.x for API v1.13
- Patch releases use the third digit: v1.13.0, v1.13.1, etc.