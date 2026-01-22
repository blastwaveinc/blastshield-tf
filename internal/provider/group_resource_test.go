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

func TestAccGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupResourceConfig("test-group-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_group.test", "name", "test-group-1"),
					resource.TestCheckResourceAttrSet("blastshield_group.test", "id"),
					resource.TestCheckResourceAttr("blastshield_group.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGroupResourceConfig("test-group-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_group.test", "name", "test-group-1-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGroupResource_withEndpoints(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResourceConfigWithEndpoints("test-group-endpoints"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_group.test", "name", "test-group-endpoints"),
					resource.TestCheckResourceAttr("blastshield_group.test", "endpoints.#", "1"),
				),
				// The endpoint's groups field is updated automatically when added to a group,
				// which causes a non-empty plan. We expect this refresh-only change.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccGroupResourceConfig(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "test" {
  name = %[1]q
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}
`, name, TestTag)
}

func testAccGroupResourceConfigWithEndpoints(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test_gateway" {
  name          = "group-test-gateway"
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
  tags = {
    test = %[2]q
  }
}

resource "blastshield_endpoint" "test" {
  name     = "group-test-endpoint"
  node_id  = blastshield_node.test_gateway.id
  endpoint = "192.168.2.100"
  enabled  = true
  tags = {
    test = %[2]q
  }
}

resource "blastshield_group" "test" {
  name = %[1]q
  tags = {
    test = %[2]q
  }
  users = []
  endpoints = [
    {
      id      = blastshield_endpoint.test.id
      expires = 0
    }
  ]
}
`, name, TestTag)
}
