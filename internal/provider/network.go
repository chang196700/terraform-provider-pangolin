package provider

// The network definitions follow the supplied Integration API and the response
// shapes in Pangolin 1.24.0 (a235635d9e5a6431ea426431e3e862e965928cd3).
// Authentication attached to public resources is deliberately owned by one
// access object: some upstream resource setters mutate an attached shared
// policy, so combining these two owners would be unsafe.

import (
	"context"
	"fmt"
	"net/netip"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func init() { registerDefinitions(networkDefinitions) }

func networkDefinitions() []Definition {
	return []Definition{networkOrganization(), networkSite(), networkPublicResource(), networkPrivateResource(), networkTarget(), networkPolicy()}
}

func networkOrganization() Definition {
	d := Definition{
		Name: "organization", Description: "An organization and its network and security settings. Pangolin 1.24.0 has an upstream Integration API defect: organization deletion requires a session user and fails for API keys. Destroy retains state when the server rejects deletion; deleting organizations requires an upstream fix.", IDField: "orgId",
		Fields: SpecFields("/org", "put", "/org/{orgId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org"}, Read: Endpoint{Method: "GET", Path: "/org/{id}", Key: "org"},
		Update: Endpoint{Method: "POST", Path: "/org/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/org/{id}"},
		List: Endpoint{Method: "GET", Path: "/orgs"}, ListKey: "orgs", Pagination: "offset", ImportFields: []string{"id"},
	}
	for _, key := range []string{"org_id", "subnet", "utility_subnet"} {
		SetField(&d, key, func(f *Field) { f.Replace = true; f.Update = false })
	}
	d.Validate = func(v Values) error {
		for _, key := range []string{"subnet", "utility_subnet"} {
			if v[key] == nil {
				continue
			}
			prefix, err := netip.ParsePrefix(Text(v[key]))
			if err != nil || !prefix.Addr().Is4() {
				return fmt.Errorf("%s must be an IPv4 CIDR", key)
			}
		}
		if v["subnet"] != nil && v["utility_subnet"] != nil {
			a, _ := netip.ParsePrefix(Text(v["subnet"]))
			b, _ := netip.ParsePrefix(Text(v["utility_subnet"]))
			if a.Overlaps(b) {
				return fmt.Errorf("subnet and utility_subnet must not overlap")
			}
		}
		return nil
	}
	d.AfterWrite = networkConfigureAfterCreate(d)
	return d
}

func networkSite() Definition {
	d := Definition{
		Name: "site", Description: "A Newt, WireGuard, or local site. Newt credentials are returned only at creation; imported sites cannot recover their secret.", IDField: "siteId",
		Fields: SpecFields("/org/{orgId}/site", "put", "/site/{siteId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/site"}, Read: Endpoint{Method: "GET", Path: "/site/{id}"},
		Update: Endpoint{Method: "POST", Path: "/site/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/site/{id}?deleteResources=false"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/sites"}, ListKey: "sites", Pagination: "page", ImportFields: []string{"org_id", "id"},
	}
	d.Fields["org_id"] = Immutable(Required(Local(String("orgId"))))
	for _, key := range []string{"type", "exit_node_id", "pub_key", "subnet", "newt_id", "secret", "address"} {
		SetField(&d, key, func(f *Field) { f.Replace = true; f.Update = false })
	}
	SetField(&d, "subnet", func(f *Field) { f.ReadAPI = "exitNodeSubnet" })
	SetField(&d, "secret", func(f *Field) { f.Sensitive = true; f.Preserve = true })
	d.Validate = func(v Values) error {
		if Text(v["type"]) == "wireguard" {
			return Require(v, "exit_node_id", "subnet", "pub_key")
		}
		if Text(v["type"]) == "newt" {
			for _, key := range []string{"exit_node_id", "subnet", "pub_key"} {
				if v[key] != nil {
					return fmt.Errorf("%s is not configurable for Newt sites; Pangolin chooses tunnel settings when Newt connects", key)
				}
			}
		}
		if v["address"] != nil {
			if _, err := netip.ParseAddr(Text(v["address"])); err != nil {
				return fmt.Errorf("address must be an IP address without a CIDR suffix")
			}
		}
		return nil
	}
	d.EnrichRead = func(_ context.Context, _ *client.Client, _ Values, o Object) error {
		if s, ok := o["address"].(string); ok {
			o["address"] = strings.SplitN(s, "/", 2)[0]
		}
		// Newt's exit-node assignment is operational state, not create input.
		if Text(o["type"]) == "newt" {
			delete(o, "exitNodeId")
			delete(o, "exitNodeSubnet")
			delete(o, "pubKey")
		}
		return nil
	}
	d.AfterWrite = networkConfigureAfterCreate(d)
	return d
}

func networkTarget() Definition {
	d := Definition{
		Name: "resource_target", Description: "A public-resource or AI-provider target, routing rules, and target health checks. Set exactly one of resource_id and ai_provider_id. Import public targets as public/<resource-id>/<target-id> and AI targets as ai/<provider-id>/<target-id>.", IDField: "targetId",
		Fields: SpecFields("/public-resource/{resourceId}/target", "put", "/target/{targetId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/public-resource/{resource_id}/target"}, Read: Endpoint{Method: "GET", Path: "/target/{id}"},
		Update: Endpoint{Method: "POST", Path: "/target/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/target/{id}"},
		List: Endpoint{Method: "GET", Path: "/public-resource/{resource_id}/targets"}, ListKey: "targets", Pagination: "offset", ImportFields: []string{"id"},
	}
	d.Fields["resource_id"] = Immutable(Local(ID("resourceId")))
	d.Fields["ai_provider_id"] = Immutable(Local(ID("providerId")))
	for _, key := range []string{"resource_id", "ai_provider_id"} {
		SetField(&d, key, func(f *Field) { f.OptionalOnly = true })
	}
	SetField(&d, "hc_headers", func(f *Field) { f.Sensitive = true })
	for _, key := range []string{"hc_interval", "hc_unhealthy_interval", "hc_timeout", "hc_healthy_threshold", "hc_unhealthy_threshold"} {
		SetField(&d, key, func(f *Field) { min := float64(1); f.Min = &min })
	}
	d.Validate = func(v Values) error {
		if err := ExactlyOne(v, "resource_id", "ai_provider_id"); err != nil {
			return err
		}
		if v["hc_enabled"] == true {
			return Require(v, "hc_hostname")
		}
		return nil
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		endpoint := d.Create
		if v["ai_provider_id"] != nil {
			endpoint.Path = "/ai-provider/{ai_provider_id}/target"
		}
		return Call(ctx, c, endpoint, v, body)
	}
	d.ImportFunc = func(id string) (Values, error) {
		parts := strings.Split(id, "/")
		v := Values{}
		if len(parts) == 3 && (parts[0] == "public" || parts[0] == "ai") {
			key := "resource_id"
			if parts[0] == "ai" {
				key = "ai_provider_id"
			}
			v[key] = parts[1]
			v["id"] = parts[2]
		} else {
			return nil, fmt.Errorf("expected public/<resource-id>/<target-id> or ai/<provider-id>/<target-id>")
		}
		for _, value := range v {
			if _, err := networkNumericID(value); err != nil {
				return nil, err
			}
		}
		return v, nil
	}
	d.AfterWrite = networkConfigureAfterCreate(d)
	return d
}

func targetDataDefinitions() []DataDefinition {
	d := networkTarget()
	fields := map[string]Field{"id": Computed(ID("targetId"))}
	for key, f := range d.Fields {
		if !f.Sensitive {
			fields[key] = f
		}
	}
	return []DataDefinition{{Name: "ai_provider_targets", Description: "Lists an AI provider's targets across all API pages. Secret health-check headers are excluded.", Entity: d, IsList: true, Inputs: map[string]Field{"ai_provider_id": Required(ID("providerId"))}, Outputs: fields, Endpoint: Endpoint{Method: "GET", Path: "/ai-provider/{ai_provider_id}/targets"}, ListKey: "targets", Pagination: "offset"}}
}

func networkPublicResource() Definition {
	d := Definition{
		Name: "public_resource", Description: "A public resource. The access object owns inline authentication and authorization, and must not be combined with resource_policy_id because Pangolin 1.24.0 resource setters can modify that shared policy.", IDField: "resourceId",
		Fields: SpecFields("/org/{orgId}/public-resource", "put", "/public-resource/{resourceId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/public-resource"}, Read: Endpoint{Method: "GET", Path: "/public-resource/{id}"},
		Update: Endpoint{Method: "POST", Path: "/public-resource/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/public-resource/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/public-resources"}, ListKey: "resources", Pagination: "page", ImportFields: []string{"org_id", "id"},
	}
	d.Fields["org_id"] = Immutable(Required(Local(String("orgId"))))
	delete(d.Fields, "http")
	delete(d.Fields, "protocol")
	delete(d.Fields, "headers")
	for _, k := range []string{"sso", "email_whitelist_enabled", "apply_rules", "skip_to_idp_id"} {
		delete(d.Fields, k)
	}
	SetField(&d, "mode", func(f *Field) {
		f.Required = true
		f.Replace = true
		f.Enum = []string{"http", "ssh", "rdp", "vnc", "tcp", "udp", "inference"}
	})
	SetField(&d, "domain_id", func(f *Field) { f.Required = false })
	SetField(&d, "proxy_port", func(f *Field) { f.Required = false })
	for _, key := range []string{"request_headers", "response_headers"} {
		SetField(&d, key, func(f *Field) { f.Sensitive = true })
	}
	d.Fields["access"] = networkAccessField()
	networkAddAIFields(&d, false)
	d.Validate = networkValidatePublic
	d.UpdateFunc = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		// Policy changes clear resource-level settings. Separate this request
		// from the remaining settings; the server dispatches using old policy state.
		if networkChanged(plan, prior, "resource_policy_id") {
			body, err := networkPayload(d, plan, "resource_policy_id")
			if err != nil {
				return err
			}
			if err = networkSend(ctx, c, d.Update.Method, d.Update.Path, plan, body); err != nil {
				return err
			}
		}
		body, err := networkUpdateBody(d, plan, prior, "resource_policy_id")
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return nil
		}
		return networkSend(ctx, c, d.Update.Method, d.Update.Path, plan, body)
	}
	d.AfterWrite = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		if prior == nil {
			if err := networkConfigureAfterCreate(d)(ctx, c, plan, prior); err != nil {
				return err
			}
		}
		if plan["resource_policy_id"] == nil && (plan["access"] != nil || prior["access"] != nil) && (prior == nil || networkChanged(plan, prior, "access", "resource_policy_id")) {
			if err := networkWritePublicAccess(ctx, c, plan, prior); err != nil {
				return err
			}
		}
		return networkWriteAI(ctx, c, d, plan, prior, false)
	}
	d.EnrichRead = func(ctx context.Context, c *client.Client, v Values, o Object) error {
		if v["access"] != nil && o["resourcePolicyId"] == nil {
			if err := networkReadPublicAccess(ctx, c, v, o); err != nil {
				return err
			}
		}
		return networkReadAI(ctx, c, v, o, false)
	}
	return d
}

func networkPrivateResource() Definition {
	d := Definition{
		Name: "private_resource", Description: "A private resource and its authoritative site, role, user, client, and inference-model associations. Changing mode or destination_port replaces the resource because Pangolin derives network, ICMP, and port-range values from them.", IDField: "siteResourceId",
		Fields: SpecFields("/org/{orgId}/private-resource", "put", "/private-resource/{siteResourceId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/private-resource"}, Read: Endpoint{Method: "GET", Path: "/private-resource/{id}"},
		Update: Endpoint{Method: "POST", Path: "/private-resource/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/private-resource/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/private-resources"}, ListKey: "siteResources", Pagination: "page", ImportFields: []string{"org_id", "id"},
	}
	d.Fields["org_id"] = Immutable(Required(Local(String("orgId"))))
	SetField(&d, "mode", func(f *Field) { f.Replace = true })
	SetField(&d, "destination_port", func(f *Field) { f.Replace = true })
	delete(d.Fields, "site_id") // siteIds is the non-deprecated authoritative relationship.
	d.Fields["site_ids"] = IDSet("siteIds")
	d.Fields["user_ids"] = StringSet("userIds")
	d.Fields["role_ids"] = IDSet("roleIds")
	d.Fields["client_ids"] = IDSet("clientIds")
	networkAddAIFields(&d, true)
	d.Validate = networkValidatePrivate
	d.AfterWrite = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		if prior == nil {
			if err := networkConfigureAfterCreate(d)(ctx, c, plan, prior); err != nil {
				return err
			}
		}
		return networkWriteAI(ctx, c, d, plan, prior, true)
	}
	d.EnrichRead = func(ctx context.Context, c *client.Client, v Values, o Object) error {
		for _, relation := range []struct{ path, key, id, api string }{{"roles", "roles", "roleId", "roleIds"}, {"users", "users", "userId", "userIds"}, {"clients", "clients", "clientId", "clientIds"}} {
			rows, err := networkGetRows(ctx, c, "/private-resource/{id}/"+relation.path, v, relation.key)
			if err != nil {
				return err
			}
			o[relation.api] = networkIDs(rows, relation.id, true)
		}
		// GET by ID omits the network's sites; the organization list is the
		// documented API that actually reports them. Walk every page.
		p, err := Path("/org/{org_id}/private-resources", v)
		if err != nil {
			return err
		}
		rows, err := c.List(ctx, p, "siteResources", "page")
		if err != nil {
			return err
		}
		found := false
		for _, row := range rows {
			if Text(row["siteResourceId"]) == Text(v["id"]) {
				if _, ok := row["siteIds"].([]any); !ok {
					return fmt.Errorf("private resource listing is missing siteIds array")
				}
				o["siteIds"] = row["siteIds"]
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("private resource was returned by ID but absent from its organization listing")
		}
		return networkReadAI(ctx, c, v, o, true)
	}
	return d
}

func networkPolicy() Definition {
	d := Definition{
		Name: "public_resource_policy", Description: "A shared public-resource policy and its authoritative authentication, authorization, whitelist, and rules. Passwords and header credentials are write-only and cannot be recovered by import.", IDField: "resourcePolicyId",
		Fields: SpecFields("/org/{orgId}/resource-policy", "post", "/public-resource-policy/{resourcePolicyId}", "post"),
		Create: Endpoint{Method: "POST", Path: "/org/{org_id}/resource-policy"}, Read: Endpoint{Method: "GET", Path: "/public-resource-policy/{id}"},
		Update: Endpoint{Method: "POST", Path: "/public-resource-policy/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/resource-policy/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/resource-policies"}, ListKey: "policies", Pagination: "page", ImportFields: []string{"org_id", "id"},
	}
	d.Fields["org_id"] = Immutable(Required(Local(String("orgId"))))
	for key, f := range networkPolicyFields() {
		d.Fields[key] = f
		f.Create = true
		f.Update = false
		d.Fields[key] = f
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		// This one create endpoint expects role IDs as JSON strings; its
		// access-control update endpoint expects JSON integers.
		roles := Strings(v["role_ids"])
		if roles == nil {
			roles = []string{}
		}
		body["roleIds"] = roles
		return Call(ctx, c, d.Create, v, body)
	}
	d.Validate = networkValidateAccess
	d.UpdateFunc = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		body, err := networkUpdateBody(d, plan, prior)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return nil
		}
		return networkSend(ctx, c, d.Update.Method, d.Update.Path, plan, body)
	}
	d.AfterWrite = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		if prior == nil {
			return networkConfigureAfterCreate(d)(ctx, c, plan, prior)
		}
		return networkWritePolicyAccess(ctx, c, d, plan, prior)
	}
	d.EnrichRead = func(_ context.Context, _ *client.Client, v Values, o Object) error {
		networkFlattenPolicy(o)
		// The server returns only the hash's ID, never header credentials.
		if o["passwordId"] == nil {
			o["password"] = nil
		}
		if o["pincodeId"] == nil {
			o["pincode"] = nil
		}
		if header, ok := o["headerAuth"].(map[string]any); !ok || header["id"] == nil {
			o["headerAuth"] = nil
		} else {
			delete(o, "headerAuth")
		}
		return nil
	}
	return d
}

// Only update-only fields need a second request following creation. The generic
// runtime persists the new ID before invoking this hook, so a later error can
// be retried or destroyed without leaking the remote object.
func networkConfigureAfterCreate(d Definition) func(context.Context, *client.Client, Values, Values) error {
	return func(ctx context.Context, c *client.Client, plan, prior Values) error {
		if prior != nil || d.Update.Path == "" {
			return nil
		}
		fields := map[string]Field{}
		v := Values{}
		for key, f := range d.Fields {
			if f.Update && !f.Create && plan[key] != nil {
				fields[key] = f
				v[key] = plan[key]
			}
		}
		body, err := Body(Definition{Fields: fields}, v, false)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return nil
		}
		return networkSend(ctx, c, d.Update.Method, d.Update.Path, plan, body)
	}
}

func networkUpdateBody(d Definition, plan, prior Values, exclude ...string) (Object, error) {
	copy := d
	copy.Fields = map[string]Field{}
	excluded := map[string]bool{}
	for _, k := range exclude {
		excluded[k] = true
	}
	for k, f := range d.Fields {
		if !excluded[k] && networkChanged(plan, prior, k) {
			copy.Fields[k] = f
		}
	}
	return Body(copy, plan, false)
}

func networkPayload(d Definition, v Values, keys ...string) (Object, error) {
	fields := map[string]Field{}
	for _, key := range keys {
		f, ok := d.Fields[key]
		if !ok {
			return nil, fmt.Errorf("unknown network field %s", key)
		}
		f.Create = true
		f.Update = true
		fields[key] = f
	}
	return Body(Definition{Fields: fields}, v, false)
}

func networkChanged(a, b Values, keys ...string) bool {
	for _, k := range keys {
		if !reflect.DeepEqual(a[k], b[k]) {
			return true
		}
	}
	return false
}
func networkRows(v any) []Object {
	var out []Object
	switch rows := v.(type) {
	case []any:
		for _, v := range rows {
			if o, ok := v.(map[string]any); ok {
				out = append(out, o)
			}
		}
	case []Object:
		return rows
	}
	return out
}
func networkIDs(rows []Object, key string, skipAdmin bool) []any {
	ids := []any{}
	for _, r := range rows {
		if skipAdmin && r["isAdmin"] == true {
			continue
		}
		if r[key] != nil {
			ids = append(ids, r[key])
		}
	}
	return ids
}
func networkGetRows(ctx context.Context, c *client.Client, path string, v Values, key string) ([]Object, error) {
	o, err := Call(ctx, c, Endpoint{Method: "GET", Path: path}, v, nil)
	if err != nil {
		return nil, err
	}
	rows, ok := o[key].([]any)
	if !ok {
		return nil, fmt.Errorf("response is missing %s collection", key)
	}
	for _, row := range rows {
		if _, ok := row.(map[string]any); !ok {
			return nil, fmt.Errorf("response %s collection contains a non-object entry", key)
		}
	}
	return networkRows(o[key]), nil
}

func networkValidatePublic(v Values) error {
	mode := Text(v["mode"])
	if mode == "" {
		return nil
	}
	if mode == "tcp" || mode == "udp" {
		if err := Require(v, "proxy_port"); err != nil {
			return err
		}
		if v["domain_id"] != nil {
			return fmt.Errorf("domain_id cannot be combined with tcp or udp mode")
		}
	} else {
		if err := Require(v, "domain_id"); err != nil {
			return err
		}
		if v["proxy_port"] != nil {
			return fmt.Errorf("proxy_port is only supported for tcp or udp mode")
		}
	}
	if v["access"] != nil && v["resource_policy_id"] != nil {
		return fmt.Errorf("access and resource_policy_id are mutually exclusive: Pangolin resource authentication setters may modify an attached shared policy; configure access on public_resource_policy instead")
	}
	if v["access"] != nil && (mode == "tcp" || mode == "udp") {
		return fmt.Errorf("access authentication requires an HTTP-routed resource")
	}
	if v["access"] != nil {
		if err := networkValidateAccess(networkAccessValues(v["access"])); err != nil {
			return err
		}
	}
	return networkValidateAI(v)
}

func networkValidatePrivate(v Values) error {
	mode := Text(v["mode"])
	if mode == "" {
		return nil
	}
	if mode != "inference" && len(Strings(v["site_ids"])) == 0 {
		return fmt.Errorf("site_ids must contain at least one site for %s mode", mode)
	}
	if mode != "inference" && mode != "gateway" && !(mode == "ssh" && Text(v["auth_daemon_mode"]) == "native") {
		if err := Require(v, "destination"); err != nil {
			return err
		}
	}
	if mode == "cidr" && v["destination"] != nil {
		if _, err := netip.ParsePrefix(Text(v["destination"])); err != nil {
			return fmt.Errorf("destination must be a CIDR for cidr mode")
		}
	}
	if mode == "http" {
		if err := Require(v, "domain_id", "scheme", "destination_port"); err != nil {
			return err
		}
	}
	if (mode == "host" || mode == "ssh") && v["destination"] != nil {
		if _, err := netip.ParseAddr(Text(v["destination"])); err != nil && v["alias"] == nil {
			return fmt.Errorf("alias is required when destination is a hostname")
		}
	}
	if mode == "ssh" && (Text(v["auth_daemon_mode"]) == "native" || (Text(v["pam_mode"]) == "push" && (Text(v["auth_daemon_mode"]) == "site" || v["auth_daemon_mode"] == nil))) && len(Strings(v["site_ids"])) > 1 {
		return fmt.Errorf("this SSH daemon mode supports only one site")
	}
	derived := Values{}
	switch mode {
	case "http", "inference":
		derived = Values{"tcp_port_range_string": "443,80", "udp_port_range_string": "", "disable_icmp": true}
	case "ssh":
		port := "22"
		if v["destination_port"] != nil {
			port = Text(v["destination_port"])
		}
		derived = Values{"tcp_port_range_string": port, "udp_port_range_string": "", "disable_icmp": true}
	case "gateway":
		derived = Values{"destination": "0.0.0.0/0", "tcp_port_range_string": "*", "udp_port_range_string": "*", "disable_icmp": false}
	}
	for key, expected := range derived {
		if v[key] != nil && !isUnknown(v[key]) && !reflect.DeepEqual(v[key], expected) {
			return fmt.Errorf("%s is derived by Pangolin in %s mode; omit it or set it to %v", key, mode, expected)
		}
	}
	return networkValidateAI(v)
}

// Keep numeric ID serialization explicit: ID fields are Terraform strings and
// are converted by Body, while textual organization/domain IDs remain strings.
func networkNumericID(v any) (int64, error) {
	n, err := strconv.ParseInt(Text(v), 10, 64)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid positive numeric ID")
	}
	return n, nil
}

func networkSortedRules(rows []Object) []any {
	sort.SliceStable(rows, func(i, j int) bool {
		a, _ := strconv.ParseInt(Text(rows[i]["priority"]), 10, 64)
		b, _ := strconv.ParseInt(Text(rows[j]["priority"]), 10, 64)
		return a < b
	})
	out := []any{}
	for _, row := range rows {
		rule := Object{}
		for _, key := range []string{"action", "match", "value", "priority", "enabled"} {
			if x, ok := row[key]; ok {
				rule[key] = x
			}
		}
		out = append(out, rule)
	}
	return out
}
