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

func TestAccServiceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceResourceConfig("test-service-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_service.test", "name", "test-service-1"),
					resource.TestCheckResourceAttrSet("blastshield_service.test", "id"),
					resource.TestCheckResourceAttr("blastshield_service.test", "tags.test", TestTag),
					resource.TestCheckResourceAttr("blastshield_service.test", "protocols.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccServiceResourceConfigUpdated("test-service-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_service.test", "name", "test-service-1-updated"),
					resource.TestCheckResourceAttr("blastshield_service.test", "protocols.#", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServiceResource_multipleProtocols(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfigMultiProtocol("test-service-multi"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_service.multi", "name", "test-service-multi"),
					resource.TestCheckResourceAttr("blastshield_service.multi", "protocols.#", "2"),
				),
			},
		},
	})
}

func testAccServiceResourceConfig(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_service" "test" {
  name = %[1]q
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6  # TCP
      ports       = ["80", "443"]
    }
  ]
}
`, name, TestTag)
}

func testAccServiceResourceConfigUpdated(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_service" "test" {
  name = %[1]q
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6  # TCP
      ports       = ["80", "443", "8080"]
    },
    {
      ip_protocol = 17  # UDP
      ports       = ["53"]
    }
  ]
}
`, name, TestTag)
}

func testAccServiceResourceConfigMultiProtocol(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_service" "multi" {
  name = %[1]q
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6  # TCP
      ports       = ["22", "3389"]
    },
    {
      ip_protocol = 1  # ICMP
      ports       = []
    }
  ]
}
`, name, TestTag)
}
