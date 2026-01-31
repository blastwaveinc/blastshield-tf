# Blastshield Terraform Provider

Terraform provider for managing [Blastshield](https://blastwave.com) resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)
- [Python](https://python.org) >= 3.8 (for code generation)

## Installation

### From Source

```bash
git clone https://github.com/blastwaveinc/blastshield-tf.git
cd blastshield-tf

# The provider supports multiple API versions
# By default, it builds from openapi/latest.json (currently v1.13)

# Build and install
make install
```

This will generate the provider code from the OpenAPI spec, build the binary, and install it to your local Terraform plugins directory.

### Building for a Specific API Version

If you need to build for a specific Blastshield API version (after `make fetch-openapi`):

```bash
# Build for v1.13
make generate OPENAPI_VERSION=v1.13
make build

# Build for v1.14 (when available)
make generate OPENAPI_VERSION=v1.14
make build
```

## Configuration

Configure the provider with your Blastshield API credentials:

```hcl
provider "blastshield" {
  host  = "https://orchestrator.blastshield.io:8000"
  token = "your-api-token"
}
```

### Environment Variables

You can also configure the provider using environment variables:

```bash
export BLASTSHIELD_HOST="https://orchestrator.blastshield.io:8000"
export BLASTSHIELD_TOKEN="your-api-token"
```

## Usage

### Managing Nodes

```hcl
resource "blastshield_node" "example" {
  name       = "my-node"
  node_type  = "A"
  api_access = false
  tags = {
    environment = "production"
  }
}

# The invitation field contains base64-encoded registration data
output "node_invitation" {
  value     = blastshield_node.example.invitation
  sensitive = true
}
```

### Managing Groups

```hcl
resource "blastshield_group" "developers" {
  name      = "developers"
  users     = []
  endpoints = []
  tags = {
    team = "engineering"
  }
}
```

### Managing Endpoints

```hcl
resource "blastshield_endpoint" "web_server" {
  name    = "web-server"
  node_id = blastshield_node.example.id
  enabled = true
  address = "10.0.0.1"

  groups = [
    {
      id      = blastshield_group.developers.id
      expires = 0
    }
  ]
}
```

### Managing Policies

```hcl
resource "blastshield_service" "https" {
  name = "https"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
  tags = {
    port = "443"
  }
}

resource "blastshield_policy" "allow_web" {
  name        = "allow-web-access"
  enabled     = true
  log         = true
  from_groups = [blastshield_group.developers.id]
  to_groups   = [blastshield_group.developers.id]
  services    = [blastshield_service.https.id]
}
```

### Managing Egress Policies

```hcl
resource "blastshield_egresspolicy" "allow_external" {
  name                  = "allow-external"
  enabled               = true
  allow_all_dns_queries = false
  groups                = [blastshield_group.developers.id]
  services              = []
  destinations          = ["example.com", "*.github.com"]
  dns_names             = []
}
```

### Managing Proxies

```hcl
resource "blastshield_proxy" "web_proxy" {
  name        = "web-proxy"
  proxy_port  = 8080
  domains     = ["internal.example.com"]
  groups      = [blastshield_group.developers.id]
  exit_agents = [blastshield_node.example.id]
}
```

### Data Sources

Query existing resources:

```hcl
# Get a single node by ID
data "blastshield_node" "existing" {
  id = "node-id-here"
}

# List all nodes
data "blastshield_nodes" "all" {}

# List nodes with filters
data "blastshield_nodes" "filtered" {
  name = "production-*"
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `blastshield_node` | Manages Blastshield nodes (agents) |
| `blastshield_endpoint` | Manages endpoints on nodes |
| `blastshield_group` | Manages groups for access control |
| `blastshield_service` | Manages service definitions |
| `blastshield_policy` | Manages network policies between groups |
| `blastshield_egresspolicy` | Manages egress policies for external access |
| `blastshield_proxy` | Manages proxy configurations |
| `blastshield_eventlogrule` | Manages event logging rules |

## Data Sources

Each resource has corresponding data sources for reading existing resources:

- `blastshield_node` / `blastshield_nodes`
- `blastshield_endpoint` / `blastshield_endpoints`
- `blastshield_group` / `blastshield_groups`
- `blastshield_service` / `blastshield_services`
- `blastshield_policy` / `blastshield_policies`
- `blastshield_egresspolicy` / `blastshield_egresspolicies`
- `blastshield_proxy` / `blastshield_proxies`
- `blastshield_eventlogrule` / `blastshield_eventlogrules`

## Development

### Code Generation

The provider code is generated from an OpenAPI specification using Jinja2 templates. The generated code is not committed to the repository.

```bash
# Generate code from OpenAPI spec (creates a temporary Python venv)
make generate

# Generate for a specific API version
make generate OPENAPI_VERSION=v1.13

# Build the provider
make build

# Run acceptance tests (requires a running Blastshield API)
make testacc

# Clean build artifacts and generated code
make clean
```

### API Version Management

This provider supports multiple Blastshield API versions. OpenAPI specifications are stored in the `openapi/` directory:

```
openapi/
├── v1.13.json          # API v1.13 specification
├── latest.json         # Symlink to current version
└── README.md           # Versioning guide
```

**Downloading a new API version:**

```bash
# Automatically downloads and names file based on version in spec
make fetch-openapi BLASTSHIELD_HOST=https://orchestrator.example.com:8000

# Update symlink if this is the latest version
cd openapi && ln -sf v1.14.json latest.json

# Generate and test
make generate OPENAPI_VERSION=v1.14
make build
make test
```

### Project Structure

```
.
├── generate.py              # Code generator script
├── openapi/                 # OpenAPI specifications
│   ├── v1.13.json          # Versioned specs
│   └── latest.json         # Symlink to current
├── codegen-templates/       # Jinja2 templates for Go code generation
│   ├── client.go.j2
│   ├── data_source.go.j2
│   ├── helpers.go.j2
│   ├── macros.j2
│   ├── provider.go.j2
│   ├── resource.go.j2
│   ├── schemas.go.j2
│   └── types.go.j2
├── internal/provider/
│   ├── client.go            # HTTP client implementation
│   ├── generated/           # Generated code (not in git)
│   └── *_test.go            # Acceptance tests
└── examples/                # Example Terraform configurations
```

### Testing

```bash
# Test API connectivity
make test-api

# Run unit tests
make test

# Run acceptance tests (requires live API)
make testacc

# Run tests for a specific resource
make testacc-node
```

### Releases and Versioning

This provider uses **semantic versioning** that matches the Blastshield API version it supports.

**Provider Version = API Version**
- Provider `v1.13.0` → Built from Blastshield API `v1.13` OpenAPI spec
- Provider `v1.14.0` → Built from Blastshield API `v1.14` OpenAPI spec

**What gets committed to Git:**
- Source code (generator, templates, client)
- OpenAPI specifications (`openapi/v*.json`)
- Tests
- **NOT** generated code (`internal/provider/generated/` is in `.gitignore`)

**Release Process:**

1. Download new API version spec:
   ```bash
   make fetch-openapi BLASTSHIELD_HOST=https://orchestrator.example.com:8000
   ```

2. Update symlink to use the new version:
   ```bash
   cd openapi && ln -sf v1.14.0.json latest.json
   ```

3. Test the provider (VERSION is automatically extracted from spec):
   ```bash
   make clean
   make build
   make test
   make testacc
   ```

4. Tag and push the release:
   ```bash
   git tag v1.14.0
   git push origin v1.14.0
   ```

5. GitHub Actions (or registry publish process) builds binaries for all platforms and publishes to the Terraform Registry

**Note:** The provider VERSION is automatically derived from the `info.version` field in the OpenAPI spec. No manual version updates needed!

**Using the provider in Terraform:**

```hcl
terraform {
  required_providers {
    blastshield = {
      source  = "blastwaveinc/blastshield"
      version = "~> 1.13.0"  # Matches Blastshield API v1.13
    }
  }
}
```

The `~>` constraint allows automatic updates to patch versions (e.g., `1.13.1`, `1.13.2`) but not minor versions (`1.14.0`).

## License

Apache 2.0 - See [LICENSE](LICENSE) for details.
