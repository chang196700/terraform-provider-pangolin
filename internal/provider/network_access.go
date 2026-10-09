package provider

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func networkPolicyFields() map[string]Field {
	fields := SpecFields("/org/{orgId}/resource-policy", "post", "", "")
	delete(fields, "name")
	fields["role_ids"] = IDSet("roleIds")
	fields["user_ids"] = StringSet("userIds")
	fields["emails"] = StringSet("emails")
	fields["sso"] = WithDefault(Bool("sso"), true)
	fields["apply_rules"] = WithDefault(Bool("applyRules"), false)
	fields["email_whitelist_enabled"] = WithDefault(Bool("emailWhitelistEnabled"), false)
	for _, key := range []string{"password", "pincode", "header_auth"} {
		f := fields[key]
		f.Sensitive = true
		f.Preserve = false
		f.Replace = false
		f.OptionalOnly = true
		f.Nullable = true
		fields[key] = f
	}
	if f, ok := fields["header_auth"]; ok {
		f.Fields["extended_compatibility"] = WithDefault(Bool("extendedCompatibility"), false)
		fields["header_auth"] = f
	}
	f := fields["skip_to_idp_id"]
	f.Nullable = true
	f.Replace = false
	f.OptionalOnly = true
	f.ReadAPI = "idpId"
	fields["skip_to_idp_id"] = f
	ruleFields := SpecFields("/public-resource/{resourceId}/rule", "put", "", "")
	ruleFields["enabled"] = WithDefault(Bool("enabled"), true)
	for key, f := range ruleFields {
		f.Replace = false
		ruleFields[key] = f
	}
	rule := Field{Kind: "object", Fields: ruleFields}
	fields["rules"] = Field{Kind: "set", API: "rules", Element: &rule, Default: []any{}, Create: true, Update: true}
	for key, f := range fields {
		f.Replace = false
		fields[key] = f
	}
	return fields
}

func networkAccessField() Field {
	return Field{Kind: "object", API: "access", Fields: networkPolicyFields(), OptionalOnly: true, Preserve: true,
		Description: "Authoritative inline access settings. Omit this object when resource_policy_id is set. Removing it resets the inline policy to SSO with no explicit users, roles, whitelist, credentials, or rules. Secrets are stored in Terraform state and cannot be recovered on import."}
}

func networkAccessValues(raw any) Values {
	v := Values{}
	if m, ok := raw.(map[string]any); ok {
		for k, x := range m {
			v[k] = x
		}
	}
	for key, f := range networkPolicyFields() {
		if v[key] == nil && f.Default != nil {
			v[key] = f.Default
		}
	}
	return v
}

func networkValidateAccess(v Values) error {
	if p := Text(v["pincode"]); p != "" && !regexp.MustCompile(`^[0-9]{6}$`).MatchString(p) {
		return fmt.Errorf("pincode must contain exactly six digits")
	}
	if len(Strings(v["emails"])) > 0 && v["email_whitelist_enabled"] != true {
		return fmt.Errorf("email_whitelist_enabled must be true when emails are configured")
	}
	for _, email := range Strings(v["emails"]) {
		if email != strings.ToLower(email) {
			return fmt.Errorf("emails must use lowercase because Pangolin normalizes email whitelist values")
		}
	}
	return nil
}

func networkWritePublicAccess(ctx context.Context, c *client.Client, plan, prior Values) error {
	a := networkAccessValues(plan["access"])
	old := networkAccessValues(prior["access"])
	d := Definition{Fields: networkPolicyFields()}
	base := "/public-resource/{id}"
	// The resource whitelist setter rejects disabled whitelists. Empty the
	// old list before disabling; enable first when adding a new list.
	if a["email_whitelist_enabled"] != true && old["email_whitelist_enabled"] == true {
		if err := networkSend(ctx, c, "POST", base+"/whitelist", plan, Object{"emails": []any{}}); err != nil {
			return err
		}
	}
	body, err := networkPayload(d, a, "sso", "email_whitelist_enabled", "apply_rules", "skip_to_idp_id")
	if err != nil {
		return err
	}
	if err = networkSend(ctx, c, "POST", base, plan, body); err != nil {
		return err
	}
	for _, group := range []struct{ suffix, key string }{{"roles", "role_ids"}, {"users", "user_ids"}, {"password", "password"}, {"pincode", "pincode"}} {
		body, err = networkPayload(d, a, group.key)
		if err != nil {
			return err
		}
		if err = networkSend(ctx, c, "POST", base+"/"+group.suffix, plan, body); err != nil {
			return err
		}
	}
	header := Object{"user": nil, "password": nil, "extendedCompatibility": nil}
	if a["header_auth"] != nil {
		wire, err := networkPayload(d, a, "header_auth")
		if err != nil {
			return err
		}
		if m, ok := wire["headerAuth"].(map[string]any); ok {
			header = m
		}
	}
	if err = networkSend(ctx, c, "POST", base+"/header-auth", plan, header); err != nil {
		return err
	}
	if a["email_whitelist_enabled"] == true {
		body, err = networkPayload(d, a, "emails")
		if err != nil {
			return err
		}
		if err = networkSend(ctx, c, "POST", base+"/whitelist", plan, body); err != nil {
			return err
		}
	}
	return networkReconcilePublicRules(ctx, c, plan, a["rules"])
}

func networkReadPublicAccess(ctx context.Context, c *client.Client, v Values, o Object) error {
	data, err := Call(ctx, c, Endpoint{Method: "GET", Path: "/public-resource/{id}/policies"}, v, nil)
	if err != nil {
		return err
	}
	policy, ok := data["defaultPolicy"].(map[string]any)
	if !ok {
		return fmt.Errorf("resource policies response is missing defaultPolicy")
	}
	networkFlattenPolicy(policy)
	old := networkAccessValues(v["access"])
	wire, err := wireValue(networkAccessField(), old)
	if err != nil {
		return err
	}
	a := wire.(map[string]any)
	for _, key := range []string{"sso", "applyRules", "emailWhitelistEnabled", "roleIds", "userIds", "emails", "rules"} {
		if x, ok := policy[key]; ok {
			a[key] = x
		}
	}
	a["idpId"] = policy["idpId"]
	if policy["passwordId"] == nil {
		a["password"] = nil
	}
	if policy["pincodeId"] == nil {
		a["pincode"] = nil
	}
	if header, ok := policy["headerAuth"].(map[string]any); !ok || header["id"] == nil {
		a["headerAuth"] = nil
	}
	o["access"] = a
	return nil
}

func networkFlattenPolicy(o Object) {
	o["roleIds"] = networkIDs(networkRows(o["roles"]), "roleId", true)
	o["userIds"] = networkIDs(networkRows(o["users"]), "userId", false)
	o["emails"] = networkIDs(networkRows(o["emailWhiteList"]), "email", false)
	o["rules"] = networkSortedRules(networkRows(o["rules"]))
}

func networkWritePolicyAccess(ctx context.Context, c *client.Client, d Definition, plan, prior Values) error {
	groups := []struct {
		suffix string
		keys   []string
	}{
		{"access-control", []string{"sso", "role_ids", "user_ids", "skip_to_idp_id"}},
		{"password", []string{"password"}}, {"pincode", []string{"pincode"}}, {"header-auth", []string{"header_auth"}},
		{"whitelist", []string{"email_whitelist_enabled", "emails"}}, {"rules", []string{"apply_rules", "rules"}},
	}
	for _, g := range groups {
		if !networkChanged(plan, prior, g.keys...) {
			continue
		}
		body, err := networkPayload(d, plan, g.keys...)
		if err != nil {
			return err
		}
		if err = networkSend(ctx, c, "POST", "/public-resource-policy/{id}/"+g.suffix, plan, body); err != nil {
			return err
		}
	}
	return nil
}

// Rules are reconciled by priority, with unchanged rows preserved. IDs belong
// to the server, while the Terraform set owns each rule's match and priority.
func networkReconcilePublicRules(ctx context.Context, c *client.Client, v Values, desired any) error {
	p, err := Path("/public-resource/{id}/rules", v)
	if err != nil {
		return err
	}
	current, err := c.List(ctx, p, "rules", "offset")
	if err != nil {
		return err
	}
	field := networkPolicyFields()["rules"]
	wire, err := wireValue(field, desired)
	if err != nil {
		return err
	}
	wanted := networkRows(wire)
	used := make([]bool, len(current))
	for _, rule := range wanted {
		match := -1
		for i, existing := range current {
			if !used[i] && Text(existing["priority"]) == Text(rule["priority"]) {
				match = i
				break
			}
		}
		if match < 0 {
			if err = networkSend(ctx, c, "PUT", "/public-resource/{id}/rule", v, rule); err != nil {
				return err
			}
			continue
		}
		used[match] = true
		old := Object{}
		for k := range rule {
			old[k] = current[match][k]
		}
		// API numbers may be json.Number while plan numbers are int64.
		if !networkRuleEqual(old, rule) {
			vv := copyValues(v)
			vv["rule_id"] = current[match]["ruleId"]
			if err = networkSend(ctx, c, "POST", "/public-resource/{id}/rule/{rule_id}", vv, rule); err != nil {
				return err
			}
		}
	}
	for i, rule := range current {
		if used[i] {
			continue
		}
		vv := copyValues(v)
		vv["rule_id"] = rule["ruleId"]
		if err = networkSend(ctx, c, "DELETE", "/public-resource/{id}/rule/{rule_id}", vv, nil); err != nil {
			return err
		}
	}
	return nil
}

func networkRuleEqual(a, b Object) bool {
	for k, v := range b {
		if k == "priority" {
			if Text(a[k]) != Text(v) {
				return false
			}
		} else if !reflect.DeepEqual(a[k], v) {
			return false
		}
	}
	return true
}
