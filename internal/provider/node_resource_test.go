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

// V1.12 Tests (without tags support)

func TestAccNodeResource_basic_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNodeResourceConfig_v112("test-node-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.test", "name", "test-node-1"),
					resource.TestCheckResourceAttr("blastshield_node.test", "node_type", "A"),
					resource.TestCheckResourceAttrSet("blastshield_node.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_node.test", "invitation"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "blastshield_node.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"invitation"}, // invitation is only returned on create
			},
			// Update and Read testing
			{
				Config: testAccNodeResourceConfig_v112("test-node-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.test", "name", "test-node-1-updated"),
					resource.TestCheckResourceAttr("blastshield_node.test", "node_type", "A"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccNodeResource_gateway_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeResourceGatewayConfig_v112("test-gateway-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.gateway", "name", "test-gateway-1"),
					resource.TestCheckResourceAttr("blastshield_node.gateway", "node_type", "G"),
					resource.TestCheckResourceAttrSet("blastshield_node.gateway", "id"),
				),
			},
		},
	})
}

func testAccNodeResourceConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test" {
  name       = %[1]q
  node_type  = "A"  # Agent
  api_access = false
}
`, name)
}

func testAccNodeResourceGatewayConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "gateway" {
  name          = %[1]q
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
}
`, name)
}

// V1.13 Tests (with tags support)

func TestAccNodeResource_basic_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNodeResourceConfig_v113("test-node-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.test", "name", "test-node-1"),
					resource.TestCheckResourceAttr("blastshield_node.test", "node_type", "A"),
					resource.TestCheckResourceAttrSet("blastshield_node.test", "id"),
					resource.TestCheckResourceAttrSet("blastshield_node.test", "invitation"),
					resource.TestCheckResourceAttr("blastshield_node.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:            "blastshield_node.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"invitation"}, // invitation is only returned on create
			},
			// Update and Read testing
			{
				Config: testAccNodeResourceConfig_v113("test-node-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.test", "name", "test-node-1-updated"),
					resource.TestCheckResourceAttr("blastshield_node.test", "node_type", "A"),
					resource.TestCheckResourceAttr("blastshield_node.test", "tags.test", TestTag),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccNodeResource_gateway_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeResourceGatewayConfig_v113("test-gateway-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.gateway", "name", "test-gateway-1"),
					resource.TestCheckResourceAttr("blastshield_node.gateway", "node_type", "G"),
					resource.TestCheckResourceAttrSet("blastshield_node.gateway", "id"),
					resource.TestCheckResourceAttr("blastshield_node.gateway", "tags.test", TestTag),
				),
			},
		},
	})
}

func testAccNodeResourceConfig_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test" {
  name       = %[1]q
  node_type  = "A"  # Agent
  api_access = false
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccNodeResourceGatewayConfig_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "gateway" {
  name          = %[1]q
  node_type     = "G"  # Gateway
  endpoint_mode = "N"  # NAT mode
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}
