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

func TestAccPolicyResource_basic_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPolicyResourceConfig_v112("test-policy-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_policy.test", "name", "test-policy-1"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "log", "true"),
					resource.TestCheckResourceAttrSet("blastshield_policy.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPolicyResourceConfigUpdated_v112("test-policy-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_policy.test", "name", "test-policy-1-updated"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "log", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPolicyResourceConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name      = "policy-test-from-group"
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name      = "policy-test-to-group"
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "policy-test-service"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
}

resource "blastshield_policy" "test" {
  name        = %[1]q
  enabled     = true
  log         = true
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
}
`, name)
}

func testAccPolicyResourceConfigUpdated_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name      = "policy-test-from-group"
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name      = "policy-test-to-group"
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "policy-test-service"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
}

resource "blastshield_policy" "test" {
  name        = %[1]q
  enabled     = false
  log         = false
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
}
`, name)
}

// V1.13 Tests (with tags)

func TestAccPolicyResource_basic_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPolicyResourceConfig_v113("test-policy-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_policy.test", "name", "test-policy-1"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "log", "true"),
					resource.TestCheckResourceAttrSet("blastshield_policy.test", "id"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "tags.test", TestTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPolicyResourceConfigUpdated_v113("test-policy-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_policy.test", "name", "test-policy-1-updated"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "log", "false"),
					resource.TestCheckResourceAttr("blastshield_policy.test", "tags.test", TestTag),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPolicyResourceConfig_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name = "policy-test-from-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name = "policy-test-to-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "policy-test-service"
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

resource "blastshield_policy" "test" {
  name        = %[1]q
  enabled     = true
  log         = true
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}

func testAccPolicyResourceConfigUpdated_v113(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name = "policy-test-from-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name = "policy-test-to-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "policy-test-service"
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

resource "blastshield_policy" "test" {
  name        = %[1]q
  enabled     = false
  log         = false
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
  tags = {
    test = %[2]q
  }
}
`, name, TestTag)
}
