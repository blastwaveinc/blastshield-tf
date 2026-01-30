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

// getAPIVersion reads the API version from the local openapi.json file
// It caches the result to avoid multiple file reads during test runs
func getAPIVersion() (string, error) {
	apiVersionOnce.Do(func() {
		// Find openapi.json - should be in project root
		// We're in internal/provider, so go up two directories
		openapiPath := filepath.Join("..", "..", "openapi.json")

		data, err := os.ReadFile(openapiPath)
		if err != nil {
			apiVersionFetchError = fmt.Errorf("failed to read openapi.json: %w", err)
			return
		}

		// Parse JSON to extract version
		var openapi struct {
			Info struct {
				Version string `json:"version"`
			} `json:"info"`
		}

		if err := json.Unmarshal(data, &openapi); err != nil {
			apiVersionFetchError = fmt.Errorf("failed to parse openapi.json: %w", err)
			return
		}

		if openapi.Info.Version == "" {
			apiVersionFetchError = fmt.Errorf("version not found in openapi.json")
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
		t.Skipf("Test requires API version >= %s, but openapi.json version is %s", requiredVersion, apiVersionStr)
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
		t.Skipf("Test is for API version < %s, but openapi.json version is %s", thresholdVersion, apiVersionStr)
	}
}

// requireAPIVersion fails the test if the API version doesn't match the requirement
// This is stricter than skipIfAPIVersionLessThan and will fail rather than skip
func requireAPIVersion(t *testing.T, requiredVersion string) {
	t.Helper()

	apiVersionStr, err := getAPIVersion()
	if err != nil {
		t.Fatalf("Could not determine API version: %v", err)
	}

	apiVersion, err := version.NewVersion(apiVersionStr)
	if err != nil {
		t.Fatalf("Could not parse API version %s: %v", apiVersionStr, err)
	}

	required, err := version.NewVersion(requiredVersion)
	if err != nil {
		t.Fatalf("Could not parse required version %s: %v", requiredVersion, err)
	}

	if apiVersion.LessThan(required) {
		t.Fatalf("Test requires API version >= %s, but openapi.json version is %s", requiredVersion, apiVersionStr)
	}
}
