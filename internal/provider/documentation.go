package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/contract"
)

// WriteDocumentation produces Registry pages, import examples, and the complete
// API disposition inventory from the same typed schemas served to Terraform.
func WriteDocumentation(root string, check bool) error {
	files := map[string][]byte{}
	for _, d := range definitions() {
		var b strings.Builder
		fmt.Fprintf(&b, "---\npage_title: \"pangolin_%s Resource - pangolin\"\ndescription: |-\n  %s\n---\n\n# pangolin_%s (Resource)\n\n%s\n\n", d.Name, strings.ReplaceAll(d.Description, "\n", " "), d.Name, d.Description)
		example := resourceExample(d)
		fmt.Fprintf(&b, "## Example Usage\n\n```terraform\n%s\n```\n\n## Schema\n\n", example)
		writeFieldDocs(&b, d.Fields, "")
		identity := strings.Join(d.ImportFields, "/")
		if d.Name == "resource_target" {
			identity = "public/resource_id/target_id or ai/provider_id/target_id"
		}
		fmt.Fprintf(&b, "\n## Import\n\n```shell\nterraform import pangolin_%s.example '%s'\n```\n\nImport identity: `%s`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.\n", d.Name, importExample(d), identity)
		if d.Name == "resource_target" {
			b.WriteString("\nAI targets use `ai/provider_id/target_id`; public targets use `public/resource_id/target_id`.\n")
		}
		b.WriteString("\n## API compatibility and permissions\n\nContract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.\n")
		files["docs/resources/"+d.Name+".md"] = []byte(b.String())
		files["examples/resources/pangolin_"+d.Name+"/resource.tf"] = []byte(example + "\n")
		files["examples/resources/pangolin_"+d.Name+"/import.sh"] = []byte("terraform import pangolin_" + d.Name + ".example '" + importExample(d) + "'\n")
	}
	for _, d := range dataDefinitions() {
		var b strings.Builder
		fmt.Fprintf(&b, "---\npage_title: \"pangolin_%s Data Source - pangolin\"\ndescription: |-\n  %s\n---\n\n# pangolin_%s (Data Source)\n\n%s\n\n## Example Usage\n\n```terraform\ndata \"pangolin_%s\" \"example\" {\n", d.Name, d.Description, d.Name, d.Description, d.Name)
		for _, k := range sortedKeys(d.Inputs) {
			f := d.Inputs[k]
			if f.Required {
				fmt.Fprintf(&b, "  %s = %s\n", k, hclLiteral(exampleValue(k, f)))
			}
		}
		b.WriteString("}\n```\n\n## Schema\n\n")
		writeFieldDocs(&b, d.fields(), "")
		files["docs/data-sources/"+d.Name+".md"] = []byte(b.String())
	}
	coverage, err := coverageDocument()
	if err != nil {
		return err
	}
	files["docs/api-coverage.md"] = coverage
	for p, data := range files {
		path := filepath.Join(root, filepath.FromSlash(p))
		if check {
			actual, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(actual, data) {
				return fmt.Errorf("generated documentation is stale: %s", p)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func writeFieldDocs(b *strings.Builder, fields map[string]Field, prefix string) {
	b.WriteString("| Attribute | Type | Usage | Description |\n|---|---|---|---|\n")
	for _, k := range sortedKeys(fields) {
		f := fields[k]
		usage := "Optional"
		if f.Required {
			usage = "Required"
		} else if f.Computed {
			usage = "Read-only"
		}
		if f.Sensitive {
			usage += ", Sensitive"
		}
		if f.Replace {
			usage += ", replacement"
		}
		desc := strings.ReplaceAll(strings.ReplaceAll(f.Description, "|", "\\|"), "\n", " ")
		if f.Default != nil {
			desc += " Default: `" + hclLiteral(f.Default) + "`."
		}
		if len(f.Enum) > 0 {
			desc += " Allowed: `" + strings.Join(f.Enum, "`, `") + "`."
		}
		if f.Preserve {
			desc += " Retained in state when the API does not return it."
		}
		fmt.Fprintf(b, "| `%s%s` | %s | %s | %s |\n", prefix, k, typeLabel(f), usage, desc)
	}
	for _, k := range sortedKeys(fields) {
		f := fields[k]
		if f.Kind == "object" {
			fmt.Fprintf(b, "\n### Nested %s%s\n\n", prefix, k)
			writeFieldDocs(b, f.Fields, prefix+k+".")
		} else if (f.Kind == "set" || f.Kind == "list") && element(f).Kind == "object" {
			fmt.Fprintf(b, "\n### Nested %s%s\n\n", prefix, k)
			writeFieldDocs(b, element(f).Fields, prefix+k+"[].")
		}
	}
}
func typeLabel(f Field) string {
	switch f.Kind {
	case "id":
		return "String (numeric API ID)"
	case "list", "set", "map":
		return strings.Title(f.Kind) + " of " + typeLabel(element(f))
	default:
		return strings.Title(f.Kind)
	}
}
func exampleValue(name string, f Field) any {
	if f.Default != nil {
		return f.Default
	}
	if len(f.Enum) > 0 {
		return f.Enum[0]
	}
	switch f.Kind {
	case "id":
		return "1"
	case "bool":
		return false
	case "int", "number":
		if f.Min != nil {
			return *f.Min
		}
		return 1
	case "set", "list":
		return []any{}
	case "object":
		v := Values{}
		for k, c := range f.Fields {
			if c.Required {
				v[k] = exampleValue(k, c)
			}
		}
		return v
	case "map":
		return Values{}
	}
	if strings.Contains(name, "url") {
		return "https://idp.example.com/"
	}
	if strings.Contains(name, "secret") || strings.Contains(name, "password") {
		return "replace-with-a-secure-value"
	}
	switch name {
	case "subnet":
		return "100.90.0.0/16"
	case "utility_subnet":
		return "100.91.0.0/16"
	case "email", "username":
		return "user@example.com"
	case "ip":
		return "127.0.0.1"
	case "org_id":
		return "example"
	case "id", "user_id", "domain_id":
		return "example-id"
	}
	return "example"
}
func resourceExample(d Definition) string {
	values := Values{}
	for k, f := range d.Fields {
		if f.Required {
			values[k] = exampleValue(k, f)
		}
	}
	switch d.Name {
	case "site":
		values["type"] = "newt"
	case "public_resource":
		values["mode"] = "http"
		values["domain_id"] = "your-domain-id"
		values["subdomain"] = "app"
	case "private_resource":
		values["mode"] = "host"
		values["destination"] = "10.0.0.10"
		values["site_ids"] = []any{"1"}
		values["tcp_port_range_string"] = "443"
		values["udp_port_range_string"] = ""
	case "resource_target":
		values["resource_id"] = "1"
		values["site_id"] = "1"
		values["port"] = 8080
		values["method"] = "http"
	case "role":
		values["reassign_role_id"] = "1"
	case "organization_user":
		values["type"] = "oidc"
		values["idp_id"] = "1"
		values["role_ids"] = []any{"1"}
	case "client":
		delete(values, "subnet")
		values["type"] = "olm"
	case "user_client":
		delete(values, "subnet")
		values["olm_id"] = "existing-user-olm-id"
		values["type"] = "olm"
	case "organization_api_key":
		values["action_ids"] = []any{"createSite", "getSite", "listSites", "updateSite", "deleteSite"}
	case "event_streaming_destination":
		values["config"] = "{}"
	case "alert_rule":
		values["event_type"] = "site_offline"
		values["all_sites"] = true
	case "health_check":
		values["site_id"] = "1"
		values["hc_mode"] = "http"
		values["hc_enabled"] = true
		values["hc_hostname"] = "127.0.0.1"
		values["hc_port"] = 8080
	case "ai_provider":
		values["type"] = "openai"
		values["upstream_url"] = "https://api.openai.com/v1"
		values["api_key"] = "replace-with-your-upstream-key"
	case "ai_budget":
		values["unit"] = "tokens"
		values["period"] = "monthly"
		values["amount"] = 1000000
	case "virtual_api_key":
		values["all_resources"] = true
	}
	var b strings.Builder
	fmt.Fprintf(&b, "resource \"pangolin_%s\" \"example\" {\n", d.Name)
	for _, k := range sortedKeys(values) {
		if _, ok := d.Fields[k]; ok {
			fmt.Fprintf(&b, "  %s = %s\n", k, hclLiteral(values[k]))
		}
	}
	b.WriteString("}")
	return strings.TrimSpace(string(hclwrite.Format([]byte(b.String()))))
}
func hclLiteral(v any) string {
	switch x := v.(type) {
	case map[string]any:
		var parts []string
		for _, k := range sortedKeys(x) {
			parts = append(parts, k+" = "+hclLiteral(x[k]))
		}
		return "{ " + strings.Join(parts, ", ") + " }"
	case string:
		return strconv.Quote(x)
	default:
		raw, e := json.Marshal(v)
		if e != nil {
			return "null"
		}
		return string(raw)
	}
}
func importExample(d Definition) string {
	if d.Name == "resource_target" {
		return "public/1/2"
	}
	parts := []string{}
	for _, k := range d.ImportFields {
		switch k {
		case "org_id":
			parts = append(parts, "example")
		case "id":
			if d.Name == "organization" {
				parts = append(parts, "example")
			} else if strings.Contains(d.Name, "api_key") || d.Name == "organization_user" || d.Name == "user_server_admin" {
				parts = append(parts, "object-id")
			} else {
				parts = append(parts, "2")
			}
		default:
			parts = append(parts, "1")
		}
	}
	return strings.Join(parts, "/")
}

func coverageDocument() ([]byte, error) {
	var b strings.Builder
	b.WriteString("# Integration API coverage\n\nSource: unchanged `openapi.yaml` (SHA-256 `d9cbd22791ecd449182a3d8778a39bb0e36b5c96e62c983abc859368a700ab52`). Behavior baseline: Pangolin 1.24.0, `a235635d9e5a6431ea426431e3e862e965928cd3`.\n\nEvery operation is classified below. `Mapped` means the behavior is owned by the named Terraform entity; it does not mean each alternate API route is called. `Alternative` routes expose behavior through another route or lookup. `Excluded` operations are not persistent managed state. This inventory is not evidence of live server validation.\n\n| Method | Path | Disposition | Terraform owner / rationale |\n|---|---|---|---|\n")
	count := 0
	for _, p := range sortedKeys(contract.Spec.Paths) {
		for _, m := range sortedKeys(contract.Spec.Paths[p]) {
			op := contract.Spec.Paths[p][m]
			if m != "get" && m != "put" && m != "post" && m != "patch" && m != "delete" {
				continue
			}
			status, reason := classifyOperation(m, p, op)
			fmt.Fprintf(&b, "| %s | `%s` | %s | %s |\n", strings.ToUpper(m), p, status, reason)
			count++
		}
	}
	fmt.Fprintf(&b, "\nTotal: **%d operations** in the supplied OpenAPI document.\n", count)
	b.WriteString("\n## Fixed-source routes omitted from `openapi.yaml`\n\nThe pinned Pangolin source registers the following Bearer Integration routes even though they are absent from the supplied OpenAPI document. The domain lifecycle routes are implemented by `pangolin_domain`; the resource-name route is exposed as `pangolin_resource_names`. The restart endpoint remains an imperative operation and is intentionally not modeled as Terraform state.\n\n| Method | Path | Provider disposition |\n|---|---|---|\n| PUT | `/org/{orgId}/domain` | Mapped: `pangolin_domain` create |\n| DELETE | `/org/{orgId}/domain/{domainId}` | Mapped: `pangolin_domain` destroy |\n| POST | `/org/{orgId}/domain/{domainId}/restart` | Excluded: imperative restart |\n| GET | `/org/{orgId}/resource-names` | Mapped: `pangolin_resource_names` |\n| GET | `/org/{orgId}/public-resource-names` | Alias of `resource-names` |\n\nThe domain update `POST` and deprecated `PATCH` routes are present in the OpenAPI document and are mapped to `pangolin_domain`; their runtime request body is supplied by the fixed server source.\n")
	if count != 312 {
		return nil, fmt.Errorf("OpenAPI operation count changed: %d; review coverage before regenerating", count)
	}
	return []byte(b.String()), nil
}
func classifyOperation(m, p string, o contract.Operation) (string, string) {
	tag := ""
	if len(o.Tags) > 0 {
		tag = o.Tags[0]
	}
	if m == "delete" && p == "/org/{orgId}" {
		return "Upstream defect", "organization: stock 1.24.0 handler dereferences session user after Bearer authentication; see compatibility.md."
	}
	if strings.Contains(p, "/my-virtual-api-keys") {
		return "Excluded", "Dashboard personal projection requires a user session; virtual_api_key uses organization-scoped Integration routes."
	}
	if m == "post" && p == "/org/{orgId}/user/{userId}" {
		return "Mapped", "organization_user"
	}
	if (m == "post" && p == "/user/{userId}/server-admin") ||
		(m == "put" && p == "/org/{orgId}/user/{userId}/client") ||
		(m == "get" && (p == "/users" || p == "/org/{orgId}/certificate/{domainId}/{domain}" || p == "/domains/namepaces")) {
		return "Unavailable in 1.24.0", "Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md."
	}
	if m == "get" && p == "/org/{orgId}/event-streaming-destination" {
		return "Corrected", "event_streaming_destination: runtime route is /org/{orgId}/event-streaming-destinations (plural)."
	}
	if m == "get" && p == "/org/{orgId}/user/{userId}" {
		return "Alternative", "organization_user: handler requires session context; use org users and username lookup with Bearer authentication."
	}
	if m == "get" && p == "/public-resource/{resourceId}/access-tokens" {
		return "Alternative", "public_resource_access_token: use the organization token list and filter resourceId; the direct route fails for root keys."
	}
	if tag == "Domain" && (m == "post" || m == "patch") && p == "/org/{orgId}/domain/{domainId}" {
		return "Mapped", "domain"
	}
	if strings.Contains(p, "/logs/") || p == "/" || strings.HasPrefix(p, "/openapi.") || p == "/maintenance/info" {
		return "Excluded", "Operational telemetry or API discovery."
	}
	for _, s := range []string{"/restart", "/approve", "/reject", "/reset-bandwidth", "/verify-associations-cache", "/rebuild-associations-cache", "/email-identity-keys", "/2fa", "check-namespace-availability"} {
		if strings.Contains(p, s) {
			return "Excluded", "Imperative workflow; not a symmetric Terraform lifecycle."
		}
	}
	if strings.Contains(p, "pick-site-defaults") {
		return "Excluded", "Preview allocation is not persistent; Newt creation allocates credentials."
	}
	if strings.Contains(p, "pick-client-defaults") {
		return "Helper", "client: called only during creation when defaults are needed."
	}
	if strings.Contains(p, "/certificate/") && m != "get" {
		return "Excluded", "Certificate issuance/restart action; certificate metadata is read-only."
	}
	if tag == "Domain" && m != "get" {
		return "Mapped", "Domain update is implemented by pangolin_domain; source-only create/delete routes are documented below."
	}
	if strings.Contains(p, "create-invite") || (tag == "User Invitation" && m != "get") {
		return "Excluded", "Invitation consumption and delivery are workflows; invitations are read-only."
	}
	if tag == "Blueprint" && m != "get" {
		return "Excluded", "Blueprint apply has no corresponding destroy lifecycle."
	}
	if opAlt := o.Deprecated || strings.Contains(tag, "Legacy"); opAlt {
		return "Alternative", "Legacy alias; use the canonical Terraform resource."
	}
	for _, s := range []string{"/nice/", "{niceId}", "user-by-username", "resources-names", "resource-names", "/admin/orgs", "/admin/org/", "/check"} {
		if strings.Contains(p, s) {
			return "Alternative", "Alternate lookup or admin projection; Terraform identifies entities by immutable ID."
		}
	}
	if p == "/user/{userId}" && m == "delete" {
		return "Excluded", "Global account deletion is outside organization membership ownership."
	}
	owners := map[string]string{"Site": "site", "Public Resource": "public_resource", "Private Resource": "private_resource", "Resource Target": "resource_target", "Organization": "organization", "Domain": "domain / domains / domain_dns_records data sources", "Role": "role / organization_user", "User": "organization_user / user / users / user_server_admin", "Client": "client / user_client / user_devices", "Rule": "public_resource.access.rules", "Public Resource Policy": "public_resource_policy / public_resource.access", "Access Token": "public_resource_access_token", "Identity Provider (Global)": "global_oidc_identity_provider / identity_provider_organization", "Identity Provider (Organization Only)": "organization_oidc_identity_provider", "API Key": "organization_api_key", "Site Provisioning Key": "site_provisioning_key", "Blueprint": "blueprint / blueprints data sources", "User Invitation": "invitations data source", "Event Streaming Destination": "event_streaming_destination", "Alert Rule": "alert_rule", "Health Check": "health_check", "AI Provider": "ai_provider / public_resource.ai_providers / private_resource.ai_providers", "AI Model": "ai_model / ai_catalog_models / resource AI associations", "AI Budget": "ai_budget", "Virtual API Key": "virtual_api_key"}
	if owner, ok := owners[tag]; ok {
		return "Mapped", owner
	}
	return "Excluded", "No persistent Terraform entity; tag " + tag + "."
}
