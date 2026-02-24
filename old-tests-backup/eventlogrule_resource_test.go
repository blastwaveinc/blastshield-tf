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

// V1.12 Tests (without tags and apply_to_groups support)

func TestAccEventLogRuleResource_basic_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEventLogRuleResourceConfig_v112("test-eventlogrule-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "name", "test-eventlogrule-1"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("blastshield_eventlogrule.test", "id"),
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
				Config: testAccEventLogRuleResourceConfigUpdated_v112("test-eventlogrule-1-updated"),
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

func TestAccEventLogRuleResource_withEmailRecipients_v112(t *testing.T) {
	skipIfAPIVersionGreaterOrEqual(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventLogRuleResourceConfigWithEmail_v112("test-eventlogrule-email"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "name", "test-eventlogrule-email"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "email_recipients.#", "2"),
				),
			},
		},
	})
}

func testAccEventLogRuleResourceConfig_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = %[1]q
  enabled = true
  conditions = [
    {
      condition_type = "category"
      operator       = "eq"
      value          = "security"
    }
  ]
  actions = ["email-notification"]
}
`, name)
}

func testAccEventLogRuleResourceConfigUpdated_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = %[1]q
  enabled = false
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
  actions = ["email-notification"]
}
`, name)
}

func testAccEventLogRuleResourceConfigWithEmail_v112(name string) string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "eventlog_group" {
  name      = "eventlog-test-group"
  users     = []
  endpoints = []
}

resource "blastshield_eventlogrule" "email" {
  name             = %[1]q
  enabled          = true
  conditions = [
    {
      condition_type = "category"
      operator       = "eq"
      value          = "security"
    }
  ]
  actions          = ["email-notification"]
  email_recipients = ["admin@example.com", "security@example.com"]
}
`, name)
}

// V1.13 Tests (with tags and apply_to_groups support)

func TestAccEventLogRuleResource_basic_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEventLogRuleResourceConfig_v113("test-eventlogrule-1"),
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
				Config: testAccEventLogRuleResourceConfigUpdated_v113("test-eventlogrule-1-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "name", "test-eventlogrule-1-updated"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "conditions.#", "2"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.test", "tags.test", TestTag),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccEventLogRuleResource_withEmailRecipients_v113(t *testing.T) {
	skipIfAPIVersionLessThan(t, "1.13.0")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventLogRuleResourceConfigWithEmail_v113("test-eventlogrule-email"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "name", "test-eventlogrule-email"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "email_recipients.#", "2"),
					resource.TestCheckResourceAttr("blastshield_eventlogrule.email", "tags.test", TestTag),
				),
			},
		},
	})
}

func testAccEventLogRuleResourceConfig_v113(name string) string {
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

func testAccEventLogRuleResourceConfigUpdated_v113(name string) string {
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

func testAccEventLogRuleResourceConfigWithEmail_v113(name string) string {
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
