# Blastshield Terraform Provider

## Project Overview

This is a Terraform provider for managing Blastshield resources, generated from an OpenAPI specification.

## Architecture

### Code Generation

The provider uses `generate.py` to automatically generate Go code from `openapi.json`. This eliminates boilerplate and ensures consistency with the API specification.

**Run the generator:**
```bash
python3 generate.py
```

**Generated files location:** `internal/provider/generated/`

- `client.go` - Client interface that the provider's Client must implement
- `types.go` - Request/Response structs for all resources
- `schemas.go` - Terraform schema definitions
- `helpers.go` - Helper functions for type conversions
- `*_resource.go` - Resource implementations (CRUD operations)
- `*_data_source.go` - Data source implementations (singular and plural)

### Hand-written Files

- `internal/provider/client.go` - API client implementation with HTTP methods
- `internal/provider/provider.go` - Provider configuration and registration

## Generator Customizations

The generator has special handling for certain resources configured at the top of `generate.py`:

### RESOURCES_WITH_GROUPS
Resources with a separate `/resource/{id}/groups` endpoint for group membership:
- Node
- Endpoint

These get an additional `groups` attribute with `[{id, expires}]` structure.

### STORE_POST_RESPONSE
Resources where POST returns a different response that should be preserved:
- Node (returns invitation data with registration tokens)

The POST response is stored as base64-encoded JSON in an `invitation` field.

### POST_RESPONSE_ID_FIELD
Custom field name in POST response that contains the entity ID:
- Node uses `node_id` (from InvitationResponse)

## Automatic Nested Field Handling

The generator automatically handles nested object arrays from the OpenAPI spec:
- **Service.protocols**: Array of `{ip_protocol, ports}` objects
- **EgressPolicy.dns_names**: Array of `{name, recursive}` objects
- **Group.endpoints/users**: Array of `{id, expires}` objects

These are generated as `ListNestedAttribute` in Terraform schemas with proper type handling for both reading from API responses and writing to API requests.

## Build Commands

```bash
# Generate code from OpenAPI spec
python3 generate.py

# Build the provider
go build ./...

# Run tests
go test ./...
```

## Key Design Decisions

1. **GET after POST**: All resources perform a GET request after POST because the API's POST endpoints don't return the full entity.

2. **Interface-based Client**: The generated code uses a `Client` interface (not a concrete type) to avoid import cycles between `provider` and `generated` packages.

3. **Raw methods for groups**: `GetGroupsRaw` and `UpdateGroupsRaw` use `interface{}` parameters to handle the `[]GroupMembership` type without import cycles.

4. **CreateRaw for POST responses**: Returns raw `[]byte` so resources can parse and store the full response (needed for Node's invitation data).

## Resources

- `blastshield_node`
- `blastshield_endpoint`
- `blastshield_group`
- `blastshield_service`
- `blastshield_policy`
- `blastshield_egress_policy`
- `blastshield_proxy`
- `blastshield_api_key`
- `blastshield_event_log_rule`

## Data Sources

Each resource has both singular and plural data sources:
- `blastshield_node` / `blastshield_nodes`
- `blastshield_endpoint` / `blastshield_endpoints`
- etc.
