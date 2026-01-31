// Copyright 2026 BlastWave, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/go-version"
)

var (
	cachedAPIVersion     string
	apiVersionOnce       sync.Once
	apiVersionFetchError error
)

// getAPIVersion reads the API version from the local openapi spec file
// It caches the result to avoid multiple file reads during test runs
// If OPENAPI_VERSION env var is set, it uses that spec file, otherwise uses openapi/latest.json
func getAPIVersion() (string, error) {
	apiVersionOnce.Do(func() {
		// Check if OPENAPI_VERSION env var is set (e.g., from Doppler or make)
		// This ensures tests use the same version that the provider was built with
		specVersion := os.Getenv("OPENAPI_VERSION")
		if specVersion == "" {
			specVersion = "latest"
		}

		// Find openapi spec - we're in internal/provider, so go up two directories
		// Note: specVersion already includes 'v' prefix (e.g., "v1.11.4" or "latest")
		openapiPath := filepath.Join("..", "..", "openapi", specVersion+".json")

		data, err := os.ReadFile(openapiPath)
		if err != nil {
			apiVersionFetchError = fmt.Errorf("failed to read openapi spec: %w", err)
			return
		}

		// Parse JSON to extract version
		var openapi struct {
			Info struct {
				Version string `json:"version"`
			} `json:"info"`
		}

		if err := json.Unmarshal(data, &openapi); err != nil {
			apiVersionFetchError = fmt.Errorf("failed to parse openapi spec: %w", err)
			return
		}

		if openapi.Info.Version == "" {
			apiVersionFetchError = fmt.Errorf("version not found in openapi spec")
			return
		}

		cachedAPIVersion = openapi.Info.Version
	})

	return cachedAPIVersion, apiVersionFetchError
}

// skipIfAPIVersionLessThan skips the test if the API version is less than the required version
// Usage: skipIfAPIVersionLessThan(t, "1.13.0")
func skipIfAPIVersionLessThan(t *testing.T, requiredVersion string) {
	t.Helper()

	apiVersionStr, err := getAPIVersion()
	if err != nil {
		t.Logf("Warning: Could not determine API version: %v", err)
		t.Logf("Continuing with test - it may fail if API doesn't support required features")
		return
	}

	apiVersion, err := version.NewVersion(apiVersionStr)
	if err != nil {
		t.Logf("Warning: Could not parse API version %s: %v", apiVersionStr, err)
		t.Logf("Continuing with test - it may fail if API doesn't support required features")
		return
	}

	required, err := version.NewVersion(requiredVersion)
	if err != nil {
		t.Logf("Warning: Could not parse required version %s: %v", requiredVersion, err)
		return
	}

	if apiVersion.LessThan(required) {
		t.Skipf("Test requires API version >= %s, but current API version is %s", requiredVersion, apiVersionStr)
	}
}

// skipIfAPIVersionGreaterOrEqual skips the test if the API version is greater than or equal to the specified version
// This is useful for running tests that only apply to older API versions
func skipIfAPIVersionGreaterOrEqual(t *testing.T, thresholdVersion string) {
	t.Helper()

	apiVersionStr, err := getAPIVersion()
	if err != nil {
		t.Logf("Warning: Could not determine API version: %v", err)
		t.Logf("Continuing with test - it may fail if API has incompatible features")
		return
	}

	apiVersion, err := version.NewVersion(apiVersionStr)
	if err != nil {
		t.Logf("Warning: Could not parse API version %s: %v", apiVersionStr, err)
		t.Logf("Continuing with test - it may fail if API has incompatible features")
		return
	}

	threshold, err := version.NewVersion(thresholdVersion)
	if err != nil {
		t.Logf("Warning: Could not parse threshold version %s: %v", thresholdVersion, err)
		return
	}

	if apiVersion.GreaterThanOrEqual(threshold) {
		t.Skipf("Test is for API version < %s, but current API version is %s", thresholdVersion, apiVersionStr)
	}
}
