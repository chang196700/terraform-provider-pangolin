package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

// Extra data sources cover API objects that have no declarative CRUD resource.
// Response fields are based on the checked Pangolin 1.24.0 source snapshot.
func extraDataDefinitions() []DataDefinition {
	domainFields := map[string]Field{
		"domain_id": String("domainId"), "base_domain": String("baseDomain"), "type": String("type"),
		"verified": Bool("verified"), "failed": Bool("failed"), "tries": Int("tries"),
		"config_managed": Bool("configManaged"), "cert_resolver": String("certResolver"),
		"prefer_wildcard_cert": Bool("preferWildcardCert"), "error_message": String("errorMessage"),
	}
	domain := DataDefinition{Name: "domain", Description: "Reads a domain available to an organization.", Inputs: extraOrgInputs(), Outputs: copyFields(domainFields), Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/domain/{domain_id}"}}
	domain.Inputs["domain_id"] = Required(Local(String("domainId")))
	domain.Outputs["custom_cert_resolver"] = String("customCertResolver")
	domain.Outputs["last_checked_at"] = Int("lastCheckedAt")
	domains := DataDefinition{Name: "domains", Description: "Lists all domains available to an organization.", Inputs: extraOrgInputs(), Outputs: domainFields, IsList: true, Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/domains"}, ListKey: "domains", Pagination: "offset"}

	// The Integration router also exposes resource names at this path, but the
	// generated OpenAPI document contains only the legacy plural spelling. Keep
	// this as a read-only data source because the response is a bare array and
	// the route has no corresponding lifecycle mutation.
	resourceNames := DataDefinition{
		Name:        "resource_names",
		Description: "Lists public resource IDs and names for an organization.",
		Inputs:      extraOrgInputs(),
		Outputs: map[string]Field{
			"resource_id": ID("resourceId"),
			"name":        String("name"),
		},
		IsList:   true,
		ListKey:  "resources",
		Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/resource-names"},
	}
	resourceNames.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		path, err := Path(resourceNames.Endpoint.Path, v)
		if err != nil {
			return nil, err
		}
		var rows []Object
		if err := c.Do(ctx, resourceNames.Endpoint.Method, path, nil, &rows); err != nil {
			return nil, err
		}
		return Object{resourceNames.ListKey: rows}, nil
	}

	dns := DataDefinition{Name: "domain_dns_records", Description: "Lists DNS records required for a Pangolin domain, including verification status and the server's effective record values.", Inputs: extraOrgInputs(), IsList: true, ListKey: "records", Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/domain/{domain_id}/dns-records"}, Outputs: map[string]Field{
		"id": ID("id"), "domain_id": String("domainId"), "record_type": String("recordType"), "base_domain": String("baseDomain"), "value": String("value"), "verified": Bool("verified"),
	}}
	dns.Inputs["domain_id"] = Required(Local(String("domainId")))
	dns.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		path, err := Path(dns.Endpoint.Path, v)
		if err != nil {
			return nil, err
		}
		// Unlike other lists, this endpoint's data is a bare array.
		var records []Object
		if err := c.Do(ctx, "GET", path, nil, &records); err != nil {
			return nil, err
		}
		return Object{"records": records}, nil
	}

	certificate := DataDefinition{Name: "certificate", Description: "OpenAPI compatibility schema for certificate issuance and renewal metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires session authentication. Requires a server version or extension that mounts this route for API keys. Certificate and private-key material are not returned.", Inputs: extraOrgInputs(), Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/certificate/{domain_id}/{domain}"}, Outputs: map[string]Field{
		"cert_id": ID("certId"), "certificate_domain": String("domain"), "domain_type": String("domainType"), "wildcard": Bool("wildcard"), "status": String("status"),
		"expires_at": Int("expiresAt"), "last_renewal_attempt": Int("lastRenewalAttempt"), "created_at": Int("createdAt"), "updated_at": Int("updatedAt"), "error_message": String("errorMessage"), "renewal_count": Int("renewalCount"),
	}}
	certificate.Inputs["domain_id"] = Required(Local(String("domainId")))
	certificate.Inputs["domain"] = Required(Local(String("domain")))

	blueprintFields := map[string]Field{"blueprint_id": ID("blueprintId"), "org_id": String("orgId"), "name": String("name"), "source": String("source"), "succeeded": Bool("succeeded"), "created_at": Int("createdAt")}
	blueprint := DataDefinition{Name: "blueprint", Description: "Reads an existing blueprint execution. This data source never applies the blueprint. Contents may include credentials and are stored as sensitive state.", Inputs: extraOrgInputs(), Outputs: copyFields(blueprintFields), Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/blueprint/{blueprint_id}"}}
	blueprint.Inputs["blueprint_id"] = Required(Local(ID("blueprintId")))
	blueprint.Outputs["message"] = String("message")
	blueprint.Outputs["contents"] = Sensitive(String("contents"))
	blueprints := DataDefinition{Name: "blueprints", Description: "Lists previous blueprint executions without retrieving their potentially sensitive contents.", Inputs: extraOrgInputs(), Outputs: blueprintFields, IsList: true, Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/blueprints"}, ListKey: "blueprints", Pagination: "offset"}

	roles := Field{Kind: "list", Element: &Field{Kind: "object", Fields: map[string]Field{"role_id": ID("roleId"), "role_name": String("roleName")}}, API: "roles"}
	invitations := DataDefinition{Name: "invitations", Description: "Lists pending organization invitations and their assigned roles. This data source never sends invitations or email.", Inputs: extraOrgInputs(), Outputs: map[string]Field{"invite_id": String("inviteId"), "email": String("email"), "expires_at": Int("expiresAt"), "roles": roles}, IsList: true, Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/invitations"}, ListKey: "invitations", Pagination: "offset"}

	return []DataDefinition{domain, domains, dns, resourceNames, certificate, blueprint, blueprints, invitations, extraUserDevices(), extraCatalogModels(), extraUser(), extraUsers()}
}

func extraOrgInputs() map[string]Field {
	return map[string]Field{"org_id": Required(Local(String("orgId")))}
}

func copyFields(fields map[string]Field) map[string]Field {
	out := make(map[string]Field, len(fields))
	for k, v := range fields {
		out[k] = v
	}
	return out
}

func extraUserDevices() DataDefinition {
	fields := map[string]Field{}
	for _, name := range []string{"orgId", "name", "pubKey", "subnet", "orgName", "type", "olmVersion", "userId", "username", "userEmail", "userType", "idpName", "idpVariant", "niceId", "agent", "approvalState", "deviceModel", "fingerprintPlatform", "fingerprintOsVersion", "fingerprintKernelVersion", "fingerprintArch", "fingerprintSerialNumber", "fingerprintUsername", "fingerprintHostname"} {
		fields[snake(name)] = String(name)
	}
	fields["client_id"] = ID("clientId")
	for _, name := range []string{"online", "olmArchived", "archived", "blocked", "olmUpdateAvailable"} {
		fields[snake(name)] = Bool(name)
	}
	for _, name := range []string{"megabytesIn", "megabytesOut"} {
		fields[snake(name)] = Number(name)
	}
	for _, name := range []string{"firstSeen", "lastSeen"} {
		fields[snake(name)] = Int(name)
	}
	inputs := extraOrgInputs()
	inputs["query"] = Local(String("query"))
	inputs["online"] = Local(Bool("online"))
	inputs["agent"] = Local(String("agent"))
	status := Local(StringSet("status"))
	status.Default = nil
	status.Description = "Optional device status filter: active, pending, denied, blocked or archived. Pangolin defaults to active and pending when omitted."
	inputs["status"] = status
	d := DataDefinition{Name: "user_devices", Description: "Lists organization user devices across all pages, including identity, network and device metadata.", Inputs: inputs, Outputs: fields, IsList: true, Endpoint: Endpoint{Method: "GET", Path: "/org/{org_id}/user-devices"}, ListKey: "devices", Pagination: "page"}
	d.ReadFunc = extraFilteredList(d, []string{"query", "online", "agent", "status"})
	return d
}

func extraCatalogModels() DataDefinition {
	d := DataDefinition{Name: "ai_catalog_models", Description: "Lists known catalog model keys. Supply provider_id, or supply both org_id and type to query before creating a provider.", Inputs: map[string]Field{"provider_id": Local(ID("providerId")), "org_id": Local(String("orgId")), "type": Local(String("type")), "query": Local(String("query"))}, Outputs: map[string]Field{"model": String("model")}, IsList: true, ListKey: "models"}
	d.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		var path string
		query := url.Values{}
		if v["provider_id"] != nil && Text(v["provider_id"]) != "" {
			if v["org_id"] != nil || v["type"] != nil {
				return nil, fmt.Errorf("provider_id cannot be combined with org_id or type")
			}
			var err error
			path, err = Path("/ai-provider/{provider_id}/catalog-models", v)
			if err != nil {
				return nil, err
			}
		} else {
			if err := Require(v, "org_id", "type"); err != nil {
				return nil, err
			}
			var err error
			path, err = Path("/org/{org_id}/ai-catalog-models", v)
			if err != nil {
				return nil, err
			}
			query.Set("type", Text(v["type"]))
		}
		if v["query"] != nil {
			query.Set("query", Text(v["query"]))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
		rows, err := c.List(ctx, path, "models", "")
		if err != nil {
			return nil, err
		}
		return Object{"models": rows}, nil
	}
	return d
}

func extraUser() DataDefinition {
	fields := map[string]Field{}
	for _, name := range []string{"userId", "email", "username", "name", "type", "idpName", "dateCreated"} {
		fields[snake(name)] = String(name)
	}
	for _, name := range []string{"twoFactorEnabled", "twoFactorSetupRequested", "emailVerified", "serverAdmin"} {
		fields[snake(name)] = Bool(name)
	}
	fields["idp_id"] = ID("idpId")
	return DataDefinition{Name: "user", Description: "Reads global user metadata by ID. This server-level API requires a root API key with getUser permission.", Inputs: map[string]Field{"user_id": Required(Local(String("userId")))}, Outputs: fields, Endpoint: Endpoint{Method: "GET", Path: "/user/{user_id}"}}
}

func extraUsers() DataDefinition {
	fields := map[string]Field{"user_id": String("id")}
	for _, name := range []string{"email", "username", "name", "dateCreated", "type", "idpName"} {
		fields[snake(name)] = String(name)
	}
	for _, name := range []string{"serverAdmin", "twoFactorEnabled", "twoFactorSetupRequested"} {
		fields[snake(name)] = Bool(name)
	}
	fields["idp_id"] = ID("idpId")
	inputs := map[string]Field{"query": Local(String("query")), "idp_id": Local(String("idp_id")), "two_factor": Local(Bool("two_factor")), "server_admin": Local(Bool("server_admin"))}
	inputs["idp_id"] = Field{Kind: "string", API: "idp_id", Description: "Filter by a numeric identity provider ID, or internal for internal users."}
	d := DataDefinition{Name: "users", Description: "OpenAPI compatibility schema for listing global user metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires a server administrator session. Requires a server version or extension that mounts this route for API keys. A root API key or listUsers action does not enable a missing route.", Inputs: inputs, Outputs: fields, IsList: true, Endpoint: Endpoint{Method: "GET", Path: "/users"}, ListKey: "users", Pagination: "page"}
	d.ReadFunc = extraFilteredList(d, []string{"query", "idp_id", "two_factor", "server_admin"})
	return d
}

func extraFilteredList(d DataDefinition, filters []string) func(context.Context, *client.Client, Values) (Object, error) {
	return func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		path, err := Path(d.Endpoint.Path, v)
		if err != nil {
			return nil, err
		}
		query := url.Values{}
		for _, name := range filters {
			value := v[name]
			if value == nil {
				continue
			}
			field := d.Inputs[name]
			if field.Kind == "set" || field.Kind == "list" {
				query.Set(field.API, strings.Join(Strings(value), ","))
			} else {
				query.Set(field.API, Text(value))
			}
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
		rows, err := c.List(ctx, path, d.ListKey, d.Pagination)
		if err != nil {
			return nil, err
		}
		return Object{d.ListKey: rows}, nil
	}
}
