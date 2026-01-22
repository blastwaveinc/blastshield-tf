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

func TestAccEventLogRuleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEventLogRuleResourceConfig("test-eventlogrule-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "name", "test-eventlogrule-1"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("blastshield_eventlogrule.test", "id"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "tags.test", TestTag),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "conditions.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "blastshield_eventlogrule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEventLogRuleResourceConfigUpdated("test-eventlogrule-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "name", "test-eventlogrule-1-updated"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "conditions.#", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccEventLogRuleResource_withEmailRecipients(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventLogRuleResourceConfigWithEmail("test-eventlogrule-email"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "name", "test-eventlogrule-email"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "email_recipients.#", "2"),
				),
			},
		},
	})
}

func testAccEventLogRuleResourceConfig(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = %[1]q
  enabled = true
  tags = {
    test = %[2]q
  }
  conditions = [
    {
      condition_type = "category"
      operator       = "eq"
      value          = "security"
    }
  ]
  actions         = ["email-notification"]
  apply_to_groups = []
}
`, name, TestTag)
}

func testAccEventLogRuleResourceConfigUpdated(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = %[1]q
  enabled = false
  tags = {
    test = %[2]q
  }
  conditions = [
    {
      condition_type = "category"
      operator       = "eq"
      value          = "security"
    },
    {
      condition_type = "priority"
      operator       = "gt"
      value          = "5"
    }
  ]
  actions         = ["email-notification"]
  apply_to_groups = []
}
`, name, TestTag)
}

func testAccEventLogRuleResourceConfigWithEmail(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "eventlog_group" {
  name = "eventlog-test-group"
  tags = {
    test = %[2]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_eventlogrule" "email" {
  name    = %[1]q
  enabled = true
  tags = {
    test = %[2]q
  }
  conditions = [
    {
      condition_type = "category"
      operator       = "eq"
      value          = "security"
    }
  ]
  actions          = ["email-notification"]
  email_recipients = ["admin@example.com", "security@example.com"]
  apply_to_groups  = [blastshield_group.eventlog_group.id]
}
`, name, TestTag)
}
