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

func TestAccEndpointResource_basic_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEndpointResourceConfig_v112("test-endpoint-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "name", "test-endpoint-1"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("blastshield_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_endpoint.test", "node_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEndpointResourceConfigUpdated_v112("test-endpoint-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "name", "test-endpoint-1-updated"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "enabled", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEndpointResourceConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test_gateway" {
  name          = "endpoint-test-gateway"
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
}

resource "blastshield_endpoint" "test" {
  name     = %[1]q
  node_id  = blastshield_node.test_gateway.id
  endpoint = "192.168.1.100"  # Required for enabled endpoints in NAT mode
  enabled  = true
}
`, name)
}

func testAccEndpointResourceConfigUpdated_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test_gateway" {
  name          = "endpoint-test-gateway"
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
}

resource "blastshield_endpoint" "test" {
  name     = %[1]q
  node_id  = blastshield_node.test_gateway.id
  endpoint = "192.168.1.101"  # Updated endpoint address
  enabled  = false
}
`, name)
}

// V1.13 Tests (with tags)

func TestAccEndpointResource_basic_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEndpointResourceConfig_v113("test-endpoint-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "name", "test-endpoint-1"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("blastshield_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_endpoint.test", "node_id"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEndpointResourceConfigUpdated_v113("test-endpoint-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "name", "test-endpoint-1-updated"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_endpoint.test", "tags.test", TestTag),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEndpointResourceConfig_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test_gateway" {
  name          = "endpoint-test-gateway"
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
  tags = {
    test = %[2]q
  }
}

resource "blastshield_endpoint" "test" {
  name     = %[1]q
  node_id  = blastshield_node.test_gateway.id
  endpoint = "192.168.1.100"  # Required for enabled endpoints in NAT mode
  enabled  = true
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccEndpointResourceConfigUpdated_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test_gateway" {
  name          = "endpoint-test-gateway"
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
  tags = {
    test = %[2]q
  }
}

resource "blastshield_endpoint" "test" {
  name     = %[1]q
  node_id  = blastshield_node.test_gateway.id
  endpoint = "192.168.1.101"  # Updated endpoint address
  enabled  = false
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}
