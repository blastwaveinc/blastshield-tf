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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	Host       string
	Token      string
	HTTPClient *http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		Host:  strings.TrimSuffix(host, "/"),
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	fullURL := c.Host + path
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Log request if TF_LOG is set
	if os.Getenv("TF_LOG") != "" {
		log.Printf("[DEBUG] Blastshield API Request: %s %s", method, fullURL)
		if len(jsonBody) > 0 {
			log.Printf("[DEBUG] Blastshield API Request Body: %s", string(jsonBody))
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response if TF_LOG is set
	if os.Getenv("TF_LOG") != "" {
		log.Printf("[DEBUG] Blastshield API Response: %d %s", resp.StatusCode, string(respBody))
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Generic CRUD operations

func (c *Client) Create(path string, body interface{}, result interface{}) error {
	respBody, err := c.doRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// CreateRaw returns the raw JSON response from a POST request
func (c *Client) CreateRaw(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body)
}

func (c *Client) Read(path string, result interface{}) error {
	respBody, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

func (c *Client) Update(path string, body interface{}, result interface{}) error {
	respBody, err := c.doRequest(http.MethodPut, path, body)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

func (c *Client) Delete(path string) error {
	_, err := c.doRequest(http.MethodDelete, path, nil)
	return err
}

func (c *Client) DeleteWithBody(path string, body interface{}) error {
	_, err := c.doRequest(http.MethodDelete, path, body)
	return err
}

func (c *Client) List(path string, params map[string]string, result interface{}) error {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		path = path + "?" + values.Encode()
	}
	respBody, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

func (c *Client) ListWithMultiParams(path string, params url.Values, result interface{}) error {
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}
	respBody, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// GetGroupsRaw gets group memberships for a resource - for generated code compatibility
// The result parameter should be a pointer to a slice (e.g., *[]GroupMembership)
func (c *Client) GetGroupsRaw(basePath string, id interface{}, result interface{}) error {
	path := fmt.Sprintf("%s%v/groups", basePath, id)
	return c.Read(path, result)
}

// UpdateGroupsRaw updates group memberships for a resource - for generated code compatibility
// groups should be a slice of group memberships, result should be a pointer to receive results
func (c *Client) UpdateGroupsRaw(basePath string, id interface{}, groups interface{}, result interface{}) error {
	path := fmt.Sprintf("%s%v/groups", basePath, id)
	req := map[string]interface{}{
		"op":     "replace",
		"groups": groups,
	}
	return c.Update(path, req, result)
}

// Node operations

type NodeCreateRequest struct {
	Name           string                 `json:"name"`
	NodeType       string                 `json:"node_type"`
	DNSName        []string               `json:"dns_name,omitempty"`
	Tags           map[string]string      `json:"tags,omitempty"`
	PublicKey      *string                `json:"public_key,omitempty"`
	EndpointMode   *string                `json:"endpoint_mode,omitempty"`
	Administrator  *string                `json:"administrator,omitempty"`
	Expires        *int64                 `json:"expires,omitempty"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	Master         *string                `json:"master,omitempty"`
	HaActive       *string                `json:"ha_active,omitempty"`
	EndpointSubnet *string                `json:"endpoint_subnet,omitempty"`
	Address        *string                `json:"address,omitempty"`
	APIAccess      *bool                  `json:"api_access,omitempty"`
}

type NodeResponse struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	NodeType          string                 `json:"node_type"`
	DNSName           []string               `json:"dns_name"`
	Tags              map[string]string      `json:"tags"`
	PublicKey         *string                `json:"public_key"`
	EndpointMode      *string                `json:"endpoint_mode"`
	Administrator     *string                `json:"administrator"`
	IDPUsername       *string                `json:"idp_username"`
	IDPAutoCreated    bool                   `json:"idp_auto_created"`
	IDPActive         *bool                  `json:"idp_active"`
	Expires           int64                  `json:"expires"`
	Settings          map[string]interface{} `json:"settings"`
	Master            *string                `json:"master"`
	HaActive          *string                `json:"ha_active"`
	EndpointSubnet    *string                `json:"endpoint_subnet"`
	Address           *string                `json:"address"`
	APIAccess         *bool                  `json:"api_access"`
	Status            map[string]interface{} `json:"status"`
	RegistrationToken *string                `json:"registration_token,omitempty"`
}

// InvitationResponse is returned when creating a node
type InvitationResponse struct {
	Type              string  `json:"type"`
	NetworkID         string  `json:"network_id"`
	NodeID            string  `json:"node_id"`
	RegistrationToken *string `json:"registration_token"`
	Offline           bool    `json:"offline"`
}

func (c *Client) CreateNode(req *NodeCreateRequest) (*InvitationResponse, error) {
	var result InvitationResponse
	err := c.Create("/nodes/", req, &result)
	return &result, err
}

func (c *Client) GetNode(id string) (*NodeResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("node ID is empty")
	}
	var result NodeResponse
	err := c.Read("/nodes/"+id, &result)
	return &result, err
}

func (c *Client) UpdateNode(id string, req interface{}) (*NodeResponse, error) {
	var result NodeResponse
	err := c.Update("/nodes/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteNode(id string) error {
	return c.Delete("/nodes/" + id)
}

// GroupWithExpiry represents a group membership with optional expiry
type GroupWithExpiry struct {
	ID      int64 `json:"id"`
	Expires int64 `json:"expires"`
}

// GroupList is used to update group memberships
type GroupList struct {
	Op     string            `json:"op,omitempty"` // "add", "replace", "remove"
	Groups []GroupWithExpiry `json:"groups"`
}

func (c *Client) GetNodeGroups(id string) ([]GroupWithExpiry, error) {
	if id == "" {
		return nil, fmt.Errorf("node ID is empty")
	}
	var result []GroupWithExpiry
	err := c.Read("/nodes/"+id+"/groups", &result)
	return result, err
}

func (c *Client) UpdateNodeGroups(id string, groups []GroupWithExpiry) ([]GroupWithExpiry, error) {
	if id == "" {
		return nil, fmt.Errorf("node ID is empty")
	}
	req := GroupList{
		Op:     "replace",
		Groups: groups,
	}
	var result []GroupWithExpiry
	err := c.Update("/nodes/"+id+"/groups", req, &result)
	return result, err
}

func (c *Client) ListNodes(params map[string]string) ([]NodeResponse, error) {
	var result []NodeResponse
	err := c.List("/nodes/", params, &result)
	return result, err
}

// Endpoint operations

type EndpointCreateRequest struct {
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	NodeID    string            `json:"node_id"`
	Enabled   bool              `json:"enabled"`
	DNSName   []string          `json:"dns_name,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Endpoint  *string           `json:"endpoint,omitempty"`
	APIAccess *bool             `json:"api_access,omitempty"`
	DefaultGW *string           `json:"default_gw,omitempty"`
}

type EndpointResponse struct {
	ID        int64                  `json:"id"`
	Name      string                 `json:"name"`
	Address   *string                `json:"address"`
	NodeID    string                 `json:"node_id"`
	Enabled   bool                   `json:"enabled"`
	DNSName   []string               `json:"dns_name"`
	Tags      map[string]string      `json:"tags"`
	Endpoint  *string                `json:"endpoint"`
	APIAccess *bool                  `json:"api_access"`
	DefaultGW *string                `json:"default_gw"`
	Status    map[string]interface{} `json:"status"`
}

func (c *Client) CreateEndpoint(req *EndpointCreateRequest) (*EndpointResponse, error) {
	var result EndpointResponse
	err := c.Create("/endpoints/", req, &result)
	return &result, err
}

func (c *Client) GetEndpoint(id string) (*EndpointResponse, error) {
	var result EndpointResponse
	err := c.Read("/endpoints/"+id, &result)
	return &result, err
}

func (c *Client) UpdateEndpoint(id string, req interface{}) (*EndpointResponse, error) {
	var result EndpointResponse
	err := c.Update("/endpoints/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteEndpoint(id string) error {
	return c.Delete("/endpoints/" + id)
}

func (c *Client) GetEndpointGroups(id string) ([]GroupWithExpiry, error) {
	if id == "" {
		return nil, fmt.Errorf("endpoint ID is empty")
	}
	var result []GroupWithExpiry
	err := c.Read("/endpoints/"+id+"/groups", &result)
	return result, err
}

func (c *Client) UpdateEndpointGroups(id string, groups []GroupWithExpiry) ([]GroupWithExpiry, error) {
	if id == "" {
		return nil, fmt.Errorf("endpoint ID is empty")
	}
	req := GroupList{
		Op:     "replace",
		Groups: groups,
	}
	var result []GroupWithExpiry
	err := c.Update("/endpoints/"+id+"/groups", req, &result)
	return result, err
}

func (c *Client) ListEndpoints(params map[string]string) ([]EndpointResponse, error) {
	var result []EndpointResponse
	err := c.List("/endpoints/", params, &result)
	return result, err
}

// Group operations

type GroupCreateRequest struct {
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags,omitempty"`
	Endpoints []int64           `json:"endpoints"`
	Users     []string          `json:"users"`
}

// GroupMemberInt represents an endpoint member with expiry
type GroupMemberInt struct {
	ID      int64 `json:"id"`
	Expires int64 `json:"expires"`
}

// GroupMemberStr represents a user/node member with expiry
type GroupMemberStr struct {
	ID      string `json:"id"`
	Expires int64  `json:"expires"`
}

type GroupResponse struct {
	ID             int64             `json:"id"`
	Name           string            `json:"name"`
	Tags           map[string]string `json:"tags"`
	IDPProvisioned bool              `json:"idp_provisioned"`
	IDPExternalID  *string           `json:"idp_externalid"`
	Endpoints      []GroupMemberInt  `json:"endpoints"`
	Users          []GroupMemberStr  `json:"users"`
}

func (c *Client) CreateGroup(req *GroupCreateRequest) (*GroupResponse, error) {
	var result GroupResponse
	err := c.Create("/groups/", req, &result)
	return &result, err
}

func (c *Client) GetGroup(id string) (*GroupResponse, error) {
	var result GroupResponse
	err := c.Read("/groups/"+id, &result)
	return &result, err
}

func (c *Client) UpdateGroup(id string, req interface{}) (*GroupResponse, error) {
	var result GroupResponse
	err := c.Update("/groups/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteGroup(id string) error {
	return c.Delete("/groups/" + id)
}

func (c *Client) ListGroups(params map[string]string) ([]GroupResponse, error) {
	var result []GroupResponse
	err := c.List("/groups/", params, &result)
	return result, err
}

func (c *Client) ListGroupsWithParams(params url.Values) ([]GroupResponse, error) {
	var result []GroupResponse
	err := c.ListWithMultiParams("/groups/", params, &result)
	return result, err
}

// Service operations

type ProtocolSpec struct {
	Protocol string `json:"protocol"`
	Port     *int64 `json:"port,omitempty"`
	EndPort  *int64 `json:"end_port,omitempty"`
	ICMPType *int64 `json:"icmp_type,omitempty"`
}

type ServiceCreateRequest struct {
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags,omitempty"`
	Protocols []ProtocolSpec    `json:"protocols"`
}

type ServiceResponse struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags"`
	Protocols []ProtocolSpec    `json:"protocols"`
}

func (c *Client) CreateService(req *ServiceCreateRequest) (*ServiceResponse, error) {
	var result ServiceResponse
	err := c.Create("/services/", req, &result)
	return &result, err
}

func (c *Client) GetService(id string) (*ServiceResponse, error) {
	var result ServiceResponse
	err := c.Read("/services/"+id, &result)
	return &result, err
}

func (c *Client) UpdateService(id string, req interface{}) (*ServiceResponse, error) {
	var result ServiceResponse
	err := c.Update("/services/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteService(id string) error {
	return c.Delete("/services/" + id)
}

func (c *Client) ListServices(params map[string]string) ([]ServiceResponse, error) {
	var result []ServiceResponse
	err := c.List("/services/", params, &result)
	return result, err
}

// Policy operations

type PolicyCreateRequest struct {
	Name       string            `json:"name"`
	Tags       map[string]string `json:"tags,omitempty"`
	Enabled    bool              `json:"enabled"`
	Log        bool              `json:"log"`
	FromGroups []int64           `json:"from_groups"`
	ToGroups   []int64           `json:"to_groups"`
	Services   []int64           `json:"services"`
}

type PolicyResponse struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	Tags       map[string]string `json:"tags"`
	Enabled    bool              `json:"enabled"`
	Log        bool              `json:"log"`
	FromGroups []int64           `json:"from_groups"`
	ToGroups   []int64           `json:"to_groups"`
	Services   []int64           `json:"services"`
}

func (c *Client) CreatePolicy(req *PolicyCreateRequest) (*PolicyResponse, error) {
	var result PolicyResponse
	err := c.Create("/policies/", req, &result)
	return &result, err
}

func (c *Client) GetPolicy(id string) (*PolicyResponse, error) {
	var result PolicyResponse
	err := c.Read("/policies/"+id, &result)
	return &result, err
}

func (c *Client) UpdatePolicy(id string, req interface{}) (*PolicyResponse, error) {
	var result PolicyResponse
	err := c.Update("/policies/"+id, req, &result)
	return &result, err
}

func (c *Client) DeletePolicy(id string) error {
	return c.Delete("/policies/" + id)
}

func (c *Client) ListPolicies(params map[string]string) ([]PolicyResponse, error) {
	var result []PolicyResponse
	err := c.List("/policies/", params, &result)
	return result, err
}

// Egress Policy operations

type DNSNameSpec struct {
	Name      string `json:"name"`
	Recursive bool   `json:"recursive"`
}

type EgressPolicyCreateRequest struct {
	Name               string            `json:"name"`
	Tags               map[string]string `json:"tags,omitempty"`
	Enabled            bool              `json:"enabled"`
	AllowAllDNSQueries bool              `json:"allow_all_dns_queries"`
	Services           []int64           `json:"services"`
	Groups             []int64           `json:"groups"`
	Destinations       []string          `json:"destinations"`
	DNSNames           []DNSNameSpec     `json:"dns_names"`
}

type EgressPolicyResponse struct {
	ID                 int64             `json:"id"`
	Name               string            `json:"name"`
	Tags               map[string]string `json:"tags"`
	Enabled            bool              `json:"enabled"`
	AllowAllDNSQueries bool              `json:"allow_all_dns_queries"`
	Services           []int64           `json:"services"`
	Groups             []int64           `json:"groups"`
	Destinations       []string          `json:"destinations"`
	DNSNames           []DNSNameSpec     `json:"dns_names"`
}

func (c *Client) CreateEgressPolicy(req *EgressPolicyCreateRequest) (*EgressPolicyResponse, error) {
	var result EgressPolicyResponse
	err := c.Create("/egress_policies/", req, &result)
	return &result, err
}

func (c *Client) GetEgressPolicy(id string) (*EgressPolicyResponse, error) {
	var result EgressPolicyResponse
	err := c.Read("/egress_policies/"+id, &result)
	return &result, err
}

func (c *Client) UpdateEgressPolicy(id string, req interface{}) (*EgressPolicyResponse, error) {
	var result EgressPolicyResponse
	err := c.Update("/egress_policies/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteEgressPolicy(id string) error {
	return c.Delete("/egress_policies/" + id)
}

func (c *Client) ListEgressPolicies(params map[string]string) ([]EgressPolicyResponse, error) {
	var result []EgressPolicyResponse
	err := c.List("/egress_policies/", params, &result)
	return result, err
}

// Proxy operations

type ProxyDomain struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

type ProxyCreateRequest struct {
	Name       string            `json:"name"`
	Tags       map[string]string `json:"tags,omitempty"`
	ProxyPort  *int64            `json:"proxy_port,omitempty"`
	Domains    []ProxyDomain     `json:"domains,omitempty"`
	Groups     []int64           `json:"groups,omitempty"`
	ExitAgents []string          `json:"exit_agents,omitempty"`
}

type ProxyResponse struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	Tags       map[string]string `json:"tags"`
	ProxyPort  int64             `json:"proxy_port"`
	Domains    []ProxyDomain     `json:"domains"`
	Groups     []int64           `json:"groups"`
	ExitAgents []string          `json:"exit_agents"`
}

func (c *Client) CreateProxy(req *ProxyCreateRequest) (*ProxyResponse, error) {
	var result ProxyResponse
	err := c.Create("/proxies/", req, &result)
	return &result, err
}

func (c *Client) GetProxy(id string) (*ProxyResponse, error) {
	var result ProxyResponse
	err := c.Read("/proxies/"+id, &result)
	return &result, err
}

func (c *Client) UpdateProxy(id string, req interface{}) (*ProxyResponse, error) {
	var result ProxyResponse
	err := c.Update("/proxies/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteProxy(id string) error {
	return c.Delete("/proxies/" + id)
}

func (c *Client) ListProxies(params map[string]string) ([]ProxyResponse, error) {
	var result []ProxyResponse
	err := c.List("/proxies/", params, &result)
	return result, err
}

// Event Log Rule operations

type EventLogRuleCreateRequest struct {
	Name            string            `json:"name"`
	Tags            map[string]string `json:"tags,omitempty"`
	Enabled         bool              `json:"enabled"`
	Conditions      []string          `json:"conditions"`
	Actions         []string          `json:"actions"`
	EmailRecipients []string          `json:"email_recipients,omitempty"`
	ApplyToGroups   []int64           `json:"apply_to_groups"`
}

type EventLogRuleResponse struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	Tags            map[string]string `json:"tags"`
	Enabled         bool              `json:"enabled"`
	Conditions      []string          `json:"conditions"`
	Actions         []string          `json:"actions"`
	EmailRecipients []string          `json:"email_recipients"`
	ApplyToGroups   []int64           `json:"apply_to_groups"`
}

func (c *Client) CreateEventLogRule(req *EventLogRuleCreateRequest) (*EventLogRuleResponse, error) {
	var result EventLogRuleResponse
	err := c.Create("/event_log_rules/", req, &result)
	return &result, err
}

func (c *Client) GetEventLogRule(id string) (*EventLogRuleResponse, error) {
	var result EventLogRuleResponse
	err := c.Read("/event_log_rules/"+id, &result)
	return &result, err
}

func (c *Client) UpdateEventLogRule(id string, req interface{}) (*EventLogRuleResponse, error) {
	var result EventLogRuleResponse
	err := c.Update("/event_log_rules/"+id, req, &result)
	return &result, err
}

func (c *Client) DeleteEventLogRule(id string) error {
	return c.Delete("/event_log_rules/" + id)
}

func (c *Client) ListEventLogRules(params map[string]string) ([]EventLogRuleResponse, error) {
	var result []EventLogRuleResponse
	err := c.List("/event_log_rules/", params, &result)
	return result, err
}

// Settings operations

type SettingsResponse struct {
	DNS             map[string]interface{} `json:"dns"`
	OverlaySubnet   map[string]interface{} `json:"overlay_subnet"`
	ConsolePassword map[string]interface{} `json:"console_password"`
	IDP             map[string]interface{} `json:"idp"`
	Syslog          map[string]interface{} `json:"syslog"`
	EULA            map[string]interface{} `json:"eula"`
	SMTP            map[string]interface{} `json:"smtp"`
	Tunnel          map[string]interface{} `json:"tunnel"`
	RemoteDesktop   map[string]interface{} `json:"remote_desktop"`
}

func (c *Client) GetSettings() (*SettingsResponse, error) {
	var result SettingsResponse
	err := c.Read("/settings/", &result)
	return &result, err
}
