package versions

import (
	"fmt"
	"sort"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// VersionedProvider provides resources and data sources for a specific API version.
type VersionedProvider interface {
	Resources() []func() resource.Resource
	DataSources() []func() datasource.DataSource
}

var (
	mu       sync.Mutex
	registry = map[string]VersionedProvider{}
)

// Register adds a versioned provider to the registry. Called from init() in each version package.
func Register(ver string, vp VersionedProvider) {
	mu.Lock()
	defer mu.Unlock()
	registry[ver] = vp
}

// SelectVersion returns the highest registered version <= serverVersion.
// Returns an error if serverVersion is below all registered versions.
func SelectVersion(serverVersion string) (VersionedProvider, string, error) {
	mu.Lock()
	defer mu.Unlock()

	sv, err := version.NewVersion(serverVersion)
	if err != nil {
		return nil, "", fmt.Errorf("invalid server version %q: %w", serverVersion, err)
	}

	var best *version.Version
	var bestKey string

	for key := range registry {
		v, err := version.NewVersion(key)
		if err != nil {
			continue
		}
		if v.GreaterThan(sv) {
			continue
		}
		if best == nil || v.GreaterThan(best) {
			best = v
			bestKey = key
		}
	}

	if best == nil {
		versions := registeredVersions()
		return nil, "", fmt.Errorf("server version %s is below all supported versions %v", serverVersion, versions)
	}

	return registry[bestKey], bestKey, nil
}

// LatestVersion returns the highest registered version.
func LatestVersion() (VersionedProvider, string) {
	mu.Lock()
	defer mu.Unlock()

	var best *version.Version
	var bestKey string

	for key := range registry {
		v, err := version.NewVersion(key)
		if err != nil {
			continue
		}
		if best == nil || v.GreaterThan(best) {
			best = v
			bestKey = key
		}
	}

	if best == nil {
		return nil, ""
	}
	return registry[bestKey], bestKey
}

func registeredVersions() []string {
	versions := make([]string, 0, len(registry))
	for key := range registry {
		versions = append(versions, key)
	}
	sort.Strings(versions)
	return versions
}
