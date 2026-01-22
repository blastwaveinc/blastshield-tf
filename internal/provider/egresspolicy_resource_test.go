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

func TestAccEgressPolicyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEgressPolicyResourceConfig("test-egress-policy-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "name", "test-egress-policy-1"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "allow_all_dns_queries", "false"),
					resource.TestCheckResourceAttrSet("blastshield_egresspolicy.test", "id"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_egresspolicy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEgressPolicyResourceConfigUpdated("test-egress-policy-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "name", "test-egress-policy-1-updated"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.test", "allow_all_dns_queries", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccEgressPolicyResource_withDestinations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEgressPolicyResourceConfigWithDestinations("test-egress-dest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_egresspolicy.dest", "name", "test-egress-dest"),
					resource.TestCheckResourceAttr("blastshield_egresspolicy.dest", "destinations.#", "2"),
				),
			},
		},
	})
}

func testAccEgressPolicyResourceConfig(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "egress_test" {
  name = "egress-test-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "egress_test" {
  name = "egress-test-service"
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
}

resource "blastshield_egresspolicy" "test" {
  name                  = %[1]q
  enabled               = true
  allow_all_dns_queries = false
  groups                = [blastshield_group.egress_test.id]
  services              = [blastshield_service.egress_test.id]
  destinations          = []
  dns_names             = []
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccEgressPolicyResourceConfigUpdated(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "egress_test" {
  name = "egress-test-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "egress_test" {
  name = "egress-test-service"
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
}

resource "blastshield_egresspolicy" "test" {
  name                  = %[1]q
  enabled               = false
  allow_all_dns_queries = true
  groups                = [blastshield_group.egress_test.id]
  services              = [blastshield_service.egress_test.id]
  destinations          = []
  dns_names             = []
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccEgressPolicyResourceConfigWithDestinations(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "egress_group" {
  name = "egress-policy-test-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "egress_service" {
  name = "egress-policy-test-service"
  tags = {
    test = %[2]q
  }
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
}

resource "blastshield_egresspolicy" "dest" {
  name                  = %[1]q
  enabled               = true
  allow_all_dns_queries = false
  groups                = [blastshield_group.egress_group.id]
  services              = [blastshield_service.egress_service.id]
  destinations          = ["10.0.0.0/8", "192.168.0.0/16"]
  dns_names             = []
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}
