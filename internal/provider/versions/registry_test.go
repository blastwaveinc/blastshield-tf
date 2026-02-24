package versions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type mockVersionedProvider struct {
	version string
}

func (m *mockVersionedProvider) Resources() []func() resource.Resource     { return nil }
func (m *mockVersionedProvider) DataSources() []func() datasource.DataSource { return nil }

func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = map[string]VersionedProvider{}
}

func TestSelectVersion_ExactMatch(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})

	vp, ver, err := SelectVersion("1.13.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver != "1.13.0" {
		t.Errorf("expected 1.13.0, got %s", ver)
	}
	if vp == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestSelectVersion_HighestBelow(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})
	Register("1.14.0", &mockVersionedProvider{version: "1.14.0"})
	Register("1.15.0", &mockVersionedProvider{version: "1.15.0"})

	vp, ver, err := SelectVersion("1.14.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver != "1.14.0" {
		t.Errorf("expected 1.14.0, got %s", ver)
	}
	if vp == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestSelectVersion_ServerAboveAll(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})
	Register("1.14.0", &mockVersionedProvider{version: "1.14.0"})

	_, ver, err := SelectVersion("1.15.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver != "1.14.0" {
		t.Errorf("expected 1.14.0, got %s", ver)
	}
}

func TestSelectVersion_ServerBelowAll(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})
	Register("1.14.0", &mockVersionedProvider{version: "1.14.0"})

	_, _, err := SelectVersion("1.12.0")
	if err == nil {
		t.Fatal("expected error for version below all supported")
	}
}

func TestSelectVersion_InvalidVersion(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})

	_, _, err := SelectVersion("not-a-version")
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}

func TestLatestVersion(t *testing.T) {
	resetRegistry()
	Register("1.13.0", &mockVersionedProvider{version: "1.13.0"})
	Register("1.15.0", &mockVersionedProvider{version: "1.15.0"})
	Register("1.14.0", &mockVersionedProvider{version: "1.14.0"})

	vp, ver := LatestVersion()
	if ver != "1.15.0" {
		t.Errorf("expected 1.15.0, got %s", ver)
	}
	if vp == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestLatestVersion_Empty(t *testing.T) {
	resetRegistry()

	vp, ver := LatestVersion()
	if ver != "" {
		t.Errorf("expected empty string, got %s", ver)
	}
	if vp != nil {
		t.Fatal("expected nil provider")
	}
}
