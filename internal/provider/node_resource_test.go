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

func TestAccNodeResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNodeResourceConfig("test-node-1"),
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
				Config: testAccNodeResourceConfig("test-node-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.test", "name", "test-node-1-updated"),
					resource.TestCheckResourceAttr("blastshield_node.test", "node_type", "A"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccNodeResource_gateway(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeResourceGatewayConfig("test-gateway-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_node.gateway", "name", "test-gateway-1"),
					resource.TestCheckResourceAttr("blastshield_node.gateway", "node_type", "G"),
					resource.TestCheckResourceAttrSet("blastshield_node.gateway", "id"),
				),
			},
		},
	})
}

func testAccNodeResourceConfig(name string) string {
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

func testAccNodeResourceGatewayConfig(name string) string {
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
