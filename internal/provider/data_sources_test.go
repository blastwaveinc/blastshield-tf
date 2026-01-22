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

// Node data sources

func TestAccNodeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_node.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_node.test", "name", "test-node-ds"),
					resource.TestCheckResourceAttr("data.blastshield_node.test", "node_type", "A"),
				),
			},
		},
	})
}

func TestAccNodesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_nodes.all", "nodes.#"),
				),
			},
		},
	})
}

func testAccNodeDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test" {
  name       = "test-node-ds"
  node_type  = "A"
  api_access = false
  tags = {
    test = %[1]q
  }
}

data "blastshield_node" "test" {
  id = blastshield_node.test.id
}
`, TestTag)
}

func testAccNodesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_node" "test" {
  name       = "test-node-ds-list"
  node_type  = "A"
  api_access = false
  tags = {
    test = %[1]q
  }
}

data "blastshield_nodes" "all" {
  depends_on = [blastshield_node.test]
}
`, TestTag)
}

// Group data sources

func TestAccGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_group.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_group.test", "name", "test-group-ds"),
				),
			},
		},
	})
}

func TestAccGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_groups.all", "groups.#"),
				),
			},
		},
	})
}

func testAccGroupDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "test" {
  name = "test-group-ds"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

data "blastshield_group" "test" {
  id = blastshield_group.test.id
}
`, TestTag)
}

func testAccGroupsDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "test" {
  name = "test-group-ds-list"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

data "blastshield_groups" "all" {
  depends_on = [blastshield_group.test]
}
`, TestTag)
}

// Service data sources

func TestAccServiceDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_service.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_service.test", "name", "test-service-ds"),
				),
			},
		},
	})
}

func TestAccServicesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServicesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_services.all", "services.#"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_service" "test" {
  name = "test-service-ds"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["8080"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

data "blastshield_service" "test" {
  id = blastshield_service.test.id
}
`, TestTag)
}

func testAccServicesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_service" "test" {
  name = "test-service-ds-list"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["8081"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

data "blastshield_services" "all" {
  depends_on = [blastshield_service.test]
}
`, TestTag)
}

// Policy data sources

func TestAccPolicyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_policy.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_policy.test", "name", "test-policy-ds"),
				),
			},
		},
	})
}

func TestAccPoliciesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_policies.all", "policies.#"),
				),
			},
		},
	})
}

func testAccPolicyDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name = "test-policy-ds-from"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name = "test-policy-ds-to"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "test-policy-ds-svc"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

resource "blastshield_policy" "test" {
  name        = "test-policy-ds"
  enabled     = true
  log         = false
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
  tags = {
    test = %[1]q
  }
}

data "blastshield_policy" "test" {
  id = blastshield_policy.test.id
}
`, TestTag)
}

func testAccPoliciesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "from" {
  name = "test-policy-ds-list-from"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_group" "to" {
  name = "test-policy-ds-list-to"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "test-policy-ds-list-svc"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["444"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

resource "blastshield_policy" "test" {
  name        = "test-policy-ds-list"
  enabled     = true
  log         = false
  from_groups = [blastshield_group.from.id]
  to_groups   = [blastshield_group.to.id]
  services    = [blastshield_service.test.id]
  tags = {
    test = %[1]q
  }
}

data "blastshield_policies" "all" {
  depends_on = [blastshield_policy.test]
}
`, TestTag)
}

// Egress Policy data sources

func TestAccEgressPolicyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEgressPolicyDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_egresspolicy.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_egresspolicy.test", "name", "test-egress-ds"),
				),
			},
		},
	})
}

func TestAccEgressPoliciesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEgressPoliciesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_egresspolicies.all", "egresspolicies.#"),
				),
			},
		},
	})
}

func testAccEgressPolicyDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "test" {
  name = "test-egress-ds-group"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "test-egress-ds-svc"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

resource "blastshield_egresspolicy" "test" {
  name                  = "test-egress-ds"
  enabled               = true
  allow_all_dns_queries = false
  groups                = [blastshield_group.test.id]
  services              = [blastshield_service.test.id]
  destinations          = []
  dns_names             = []
  tags = {
    test = %[1]q
  }
}

data "blastshield_egresspolicy" "test" {
  id = blastshield_egresspolicy.test.id
}
`, TestTag)
}

func testAccEgressPoliciesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_group" "test" {
  name = "test-egress-ds-list-group"
  tags = {
    test = %[1]q
  }
  users     = []
  endpoints = []
}

resource "blastshield_service" "test" {
  name = "test-egress-ds-list-svc"
  protocols = [
    {
      ip_protocol = 6
      ports       = ["443"]
    }
  ]
  tags = {
    test = %[1]q
  }
}

resource "blastshield_egresspolicy" "test" {
  name                  = "test-egress-ds-list"
  enabled               = true
  allow_all_dns_queries = false
  groups                = [blastshield_group.test.id]
  services              = [blastshield_service.test.id]
  destinations          = []
  dns_names             = []
  tags = {
    test = %[1]q
  }
}

data "blastshield_egresspolicies" "all" {
  depends_on = [blastshield_egresspolicy.test]
}
`, TestTag)
}

// Proxy data sources

func TestAccProxyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_proxy.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_proxy.test", "name", "test-proxy-ds"),
				),
			},
		},
	})
}

func TestAccProxiesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProxiesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_proxies.all", "proxies.#"),
				),
			},
		},
	})
}

func testAccProxyDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name = "test-proxy-ds"
  tags = {
    test = %[1]q
  }
}

data "blastshield_proxy" "test" {
  id = blastshield_proxy.test.id
}
`, TestTag)
}

func testAccProxiesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_proxy" "test" {
  name = "test-proxy-ds-list"
  tags = {
    test = %[1]q
  }
}

data "blastshield_proxies" "all" {
  depends_on = [blastshield_proxy.test]
}
`, TestTag)
}

// Event Log Rule data sources

func TestAccEventLogRuleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventLogRuleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_eventlogrule.test", "id"),
					resource.TestCheckResourceAttr("data.blastshield_eventlogrule.test", "name", "test-eventlogrule-ds"),
				),
			},
		},
	})
}

func TestAccEventLogRulesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventLogRulesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.blastshield_eventlogrules.all", "eventlogrules.#"),
				),
			},
		},
	})
}

func testAccEventLogRuleDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = "test-eventlogrule-ds"
  enabled = true
  tags = {
    test = %[1]q
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

data "blastshield_eventlogrule" "test" {
  id = blastshield_eventlogrule.test.id
}
`, TestTag)
}

func testAccEventLogRulesDataSourceConfig() string {
	return testAccProviderConfig() + fmt.Sprintf(`
resource "blastshield_eventlogrule" "test" {
  name    = "test-eventlogrule-ds-list"
  enabled = true
  tags = {
    test = %[1]q
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

data "blastshield_eventlogrules" "all" {
  depends_on = [blastshield_eventlogrule.test]
}
`, TestTag)
}
