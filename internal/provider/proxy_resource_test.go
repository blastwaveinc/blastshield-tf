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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// V1.12 Tests (without tags)

func TestAccProxyResource_basic_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProxyResourceConfig_v112("test-proxy-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.test", "name", "test-proxy-1"),
					resource.TestCheckResourceAttrSet("blastshield_proxy.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_proxy.test", "proxy_port"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_proxy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProxyResourceConfigUpdated_v112("test-proxy-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.test", "name", "test-proxy-1-updated"),
					resource.TestCheckResourceAttr("blastshield_proxy.test", "domains.#", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProxyResource_withGroups_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyResourceConfigWithGroups_v112("test-proxy-groups"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.groups", "name", "test-proxy-groups"),
					resource.TestCheckResourceAttr("blastshield_proxy.groups", "groups.#", "1"),
				),
			},
		},
	})
}

func testAccProxyResourceConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name    = %[1]q
  domains = ["example.com"]
}
`, name)
}

func testAccProxyResourceConfigUpdated_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name    = %[1]q
  domains = ["example.com", "test.example.com"]
}
`, name)
}

func testAccProxyResourceConfigWithGroups_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "proxy_group" {
  name      = "proxy-test-group"
  users     = []
  endpoints = []
}

resource "blastshield_proxy" "groups" {
  name    = %[1]q
  domains = ["proxy.example.com"]
  groups  = [blastshield_group.proxy_group.id]
}
`, name)
}

// V1.13 Tests (with tags)

func TestAccProxyResource_basic_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProxyResourceConfig_v113("test-proxy-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.test", "name", "test-proxy-1"),
					resource.TestCheckResourceAttrSet("blastshield_proxy.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_proxy.test", "proxy_port"),
					resource.TestCheckResourceAttr("blastshield_proxy.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_proxy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProxyResourceConfigUpdated_v113("test-proxy-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.test", "name", "test-proxy-1-updated"),
					resource.TestCheckResourceAttr("blastshield_proxy.test", "domains.#", "2"),
					resource.TestCheckResourceAttr("blastshield_proxy.test", "tags.test", TestTag),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProxyResource_withGroups_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyResourceConfigWithGroups_v113("test-proxy-groups"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_proxy.groups", "name", "test-proxy-groups"),
					resource.TestCheckResourceAttr("blastshield_proxy.groups", "groups.#", "1"),
					resource.TestCheckResourceAttr("blastshield_proxy.groups", "tags.test", TestTag),
				),
			},
		},
	})
}

func testAccProxyResourceConfig_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name    = %[1]q
  domains = ["example.com"]
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccProxyResourceConfigUpdated_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name    = %[1]q
  domains = ["example.com", "test.example.com"]
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccProxyResourceConfigWithGroups_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "proxy_group" {
  name = "proxy-test-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_proxy" "groups" {
  name    = %[1]q
  domains = ["proxy.example.com"]
  groups  = [blastshield_group.proxy_group.id]
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}
