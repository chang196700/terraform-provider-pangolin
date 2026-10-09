package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func init() { registerDefinitions(identityDefinitions) }

func identityDefinitions() []Definition {
	defs := []Definition{roleDefinition(), organizationUserDefinition(), serverAdminDefinition(), clientDefinition(false), clientDefinition(true), oidcDefinition(true), oidcDefinition(false), identityProviderOrganizationDefinition()}
	for i := range defs {
		defs[i].ListKey = defs[i].List.Key
	}
	return defs
}

func identityOrgField() Field {
	f := Required(Immutable(Local(String("orgId"))))
	f.Description = "Organization identifier."
	return f
}

func roleDefinition() Definition {
	d := Definition{
		Name: "role", Description: "An organization role. Deletion migrates its users to reassign_role_id; choose a role with appropriate permissions before destroying this resource.", IDField: "roleId",
		Fields: SpecFields("/org/{orgId}/role", "put", "/role/{roleId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/role"}, Read: Endpoint{Method: "GET", Path: "/role/{id}"},
		Update: Endpoint{Method: "POST", Path: "/role/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/role/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/roles", Key: "roles"}, Pagination: "page", ImportFields: []string{"org_id", "id"}, Sources: []string{"role", "roles"},
	}
	d.Fields["org_id"] = identityOrgField()
	d.Fields["reassign_role_id"] = Required(Local(ID("")))
	SetField(&d, "reassign_role_id", func(f *Field) {
		f.Description = "Existing role in the same organization to receive users when this role is deleted. Required even when the role currently has no members."
	})
	SetField(&d, "allow_ssh", func(f *Field) {
		f.Description = "Whether to allow SSH signing. Read from the organization's role list because the individual role endpoint omits this permission."
	})
	for _, name := range []string{"ssh_sudo_commands", "ssh_unix_groups"} {
		SetField(&d, name, func(f *Field) { f.Kind = "set"; f.Default = []any{} })
	}
	d.EnrichRead = func(_ context.Context, _ *client.Client, _ Values, obj Object) error {
		for _, key := range []string{"sshSudoCommands", "sshUnixGroups"} {
			if raw, ok := obj[key].(string); ok {
				var values []string
				if err := json.Unmarshal([]byte(raw), &values); err != nil {
					return fmt.Errorf("invalid role %s response: %w", key, err)
				}
				obj[key] = values
			}
		}
		return nil
	}
	d.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		obj, err := Call(ctx, c, d.Read, v, nil)
		if err != nil {
			return nil, err
		}
		path, err := Path(d.List.Path, v)
		if err != nil {
			return nil, err
		}
		rows, err := c.List(ctx, path, "roles", "page")
		if err != nil {
			return nil, fmt.Errorf("role exists but its SSH permission could not be read: %v", err)
		}
		for _, row := range rows {
			if Text(row["roleId"]) == Text(v["id"]) {
				allow, ok := row["allowSsh"].(bool)
				if !ok {
					return nil, fmt.Errorf("role list response is missing allowSsh")
				}
				obj["allowSsh"] = allow
				return obj, nil
			}
		}
		return nil, fmt.Errorf("role exists but is missing from the organization's role list")
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		obj, err := Call(ctx, c, d.Create, v, body)
		if err != nil {
			return nil, err
		}
		if err = d.EnrichRead(ctx, c, v, obj); err != nil {
			return obj, err
		}
		return obj, nil
	}
	d.Validate = func(v Values) error {
		if v["id"] != nil && Text(v["id"]) == Text(v["reassign_role_id"]) {
			return fmt.Errorf("reassign_role_id must differ from the role being deleted")
		}
		return nil
	}
	d.DeleteFunc = func(ctx context.Context, c *client.Client, v Values) error {
		if err := Require(v, "reassign_role_id"); err != nil {
			return err
		}
		if err := d.Validate(v); err != nil {
			return err
		}
		// A missing migration target must never be mistaken for the managed role
		// already being absent by the shared Delete implementation.
		if _, err := Call(ctx, c, d.Read, v, nil); err != nil {
			return err
		}
		target, err := Call(ctx, c, Endpoint{Method: "GET", Path: "/role/{reassign_role_id}"}, v, nil)
		if err != nil {
			return fmt.Errorf("cannot validate reassignment role: %v", err)
		}
		if Text(target["orgId"]) != Text(v["org_id"]) {
			return fmt.Errorf("reassignment role belongs to a different organization")
		}
		n, err := strconv.ParseInt(Text(v["reassign_role_id"]), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid reassign_role_id: %w", err)
		}
		return identityWrite(ctx, c, d.Delete, v, Object{"roleId": n})
	}
	return d
}

func organizationUserDefinition() Definition {
	d := Definition{
		Name: "organization_user", Description: "An OIDC user's membership and authoritative role assignments in an organization. Destroy removes membership, not the global user. The API does not support creating internal users. In Pangolin 1.24.0, role reassignment endpoints require the enterprise/SaaS integration router. Reads use the organization list and username lookup because the individual membership GET incorrectly requires a browser session.", IDField: "userId",
		Fields: SpecFields("/org/{orgId}/user", "put", "/org/{orgId}/user/{userId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/user"}, Read: Endpoint{Method: "GET", Path: "/org/{org_id}/user/{id}"},
		Update: Endpoint{Method: "POST", Path: "/org/{org_id}/user/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/user/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/users", Key: "users"}, Pagination: "page", ImportFields: []string{"org_id", "id"}, Sources: []string{"organization_user", "organization_users"},
	}
	d.Fields["org_id"] = identityOrgField()
	delete(d.Fields, "role_id") // The singular legacy field duplicates role_ids.
	d.Fields["role_ids"] = Required(IDSet("roleIds"))
	SetField(&d, "role_ids", func(f *Field) {
		f.Default = nil
		f.Update = false
		f.Description = "Nonempty authoritative role ID set. New assignments are added before obsolete ones are removed."
	})
	SetField(&d, "type", func(f *Field) { f.Default = "oidc"; f.Enum = []string{"oidc"} })
	SetField(&d, "idp_id", func(f *Field) { f.Required = true; f.Kind = "id" })
	d.Validate = func(v Values) error {
		if v["role_ids"] != nil && len(Strings(v["role_ids"])) == 0 {
			return fmt.Errorf("role_ids must contain at least one role")
		}
		for _, name := range []string{"username", "email"} {
			if value := Text(v[name]); value != strings.ToLower(value) {
				return fmt.Errorf("%s must be lowercase because Pangolin normalizes it", name)
			}
		}
		return nil
	}
	d.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		path, err := Path(d.List.Path, v)
		if err != nil {
			return nil, err
		}
		rows, err := c.List(ctx, path, "users", "page")
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			// listUsers projects users.userId as id, unlike queryUser.
			if Text(row["id"]) == Text(v["id"]) {
				return lookupListedOrganizationUser(ctx, c, v, row)
			}
		}
		return nil, client.NotFound(path)
	}
	d.EnrichRead = func(ctx context.Context, c *client.Client, v Values, obj Object) error {
		if obj["userId"] != nil {
			return nil
		}
		// The collection has id and roles, but omits autoProvisioned. Use the
		// mounted Bearer-compatible username lookup for complete datasource rows.
		full, err := lookupListedOrganizationUser(ctx, c, v, obj)
		if err != nil {
			return err
		}
		for key, value := range full {
			obj[key] = value
		}
		return nil
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		if err = identityWrite(ctx, c, d.Create, v, body); err != nil {
			return nil, err
		}
		path, err := Path("/org/{org_id}/user-by-username", v)
		if err != nil {
			return nil, err
		}
		q := url.Values{"username": {Text(v["username"])}, "idpId": {Text(v["idp_id"])}}
		var obj Object
		if err = c.Do(ctx, "GET", path+"?"+q.Encode(), nil, &obj); err != nil {
			return nil, fmt.Errorf("membership was created but its user ID could not be resolved; import the existing org_id/user_id before retrying: %w", err)
		}
		return obj, nil
	}
	d.AfterWrite = func(ctx context.Context, c *client.Client, v, prior Values) error {
		if prior == nil && v["auto_provisioned"] != nil {
			if err := identityWrite(ctx, c, d.Update, v, Object{"autoProvisioned": v["auto_provisioned"]}); err != nil {
				return err
			}
		}
		// Read actual assignments before diffing, so retrying a partially applied
		// update does not duplicate additions or rely on stale Terraform state.
		current, err := d.ReadFunc(ctx, c, v)
		if err != nil {
			return err
		}
		return reconcileUserRoles(ctx, c, v, Strings(current["roleIds"]), Strings(v["role_ids"]))
	}
	return d
}

func lookupListedOrganizationUser(ctx context.Context, c *client.Client, v Values, row Object) (Object, error) {
	if row["id"] == nil || Text(row["username"]) == "" {
		return nil, fmt.Errorf("organization user list response is missing id or username")
	}
	path, err := Path("/org/{org_id}/user-by-username", v)
	if err != nil {
		return nil, err
	}
	q := url.Values{"username": {Text(row["username"])}}
	if row["idpId"] != nil {
		q.Set("idpId", Text(row["idpId"]))
	}
	var obj Object
	if err = c.Do(ctx, "GET", path+"?"+q.Encode(), nil, &obj); err != nil {
		// The list already proved membership exists. A secondary lookup failure
		// (including 404 after a concurrent rename) must not erase Terraform state.
		return nil, fmt.Errorf("organization membership exists but its details could not be read: %v", err)
	}
	if Text(obj["userId"]) != Text(row["id"]) {
		return nil, fmt.Errorf("organization user lookup returned a different user ID")
	}
	return obj, nil
}

func reconcileUserRoles(ctx context.Context, c *client.Client, v Values, current, desired []string) error {
	if len(desired) == 0 {
		return fmt.Errorf("role_ids must contain at least one role")
	}
	old, next := map[string]bool{}, map[string]bool{}
	for _, role := range current {
		old[role] = true
	}
	for _, role := range desired {
		next[role] = true
	}
	for _, role := range desired {
		if old[role] {
			continue
		}
		path, err := Path("/user/{id}/add-role/{role_id}", Values{"id": v["id"], "role_id": role})
		if err != nil {
			return err
		}
		if err = c.Do(ctx, "POST", path, nil, nil); err != nil {
			return fmt.Errorf("adding user role %s: %w", role, err)
		}
	}
	for _, role := range current {
		if next[role] {
			continue
		}
		path, err := Path("/user/{id}/remove-role/{role_id}", Values{"id": v["id"], "role_id": role})
		if err != nil {
			return err
		}
		if err = c.Do(ctx, "DELETE", path, nil, nil); err != nil {
			return fmt.Errorf("removing user role %s: %w", role, err)
		}
	}
	return nil
}

func serverAdminDefinition() Definition {
	d := Definition{
		Name: "user_server_admin", Description: "Specification adapter for administrator status. NOT AVAILABLE with the Pangolin 1.24.0 Integration API: its create, update, and destroy action is mounted only in the browser-session API, with no Bearer-key alternative. The root-key individual user read exists, but does not make this a usable lifecycle. Requires a server version or custom deployment that exposes the documented action. Mock lifecycle tests verify adapter behavior only. Destroy demotes the user rather than deleting it.", IDField: "userId",
		Fields: SpecFields("/user/{userId}/server-admin", "post", "/user/{userId}/server-admin", "post"),
		Create: Endpoint{Method: "POST", Path: "/user/{user_id}/server-admin"}, Read: Endpoint{Method: "GET", Path: "/user/{id}"}, Update: Endpoint{Method: "POST", Path: "/user/{id}/server-admin"},
		ImportFields: []string{"id"},
	}
	d.Fields["user_id"] = Required(Immutable(Local(String("userId"))))
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		obj, err := Call(ctx, c, d.Create, v, body)
		return obj, serverAdminActionError(err)
	}
	d.UpdateFunc = func(ctx context.Context, c *client.Client, v, _ Values) error {
		body, err := Body(d, v, false)
		if err != nil {
			return err
		}
		return serverAdminActionError(identityWrite(ctx, c, d.Update, v, body))
	}
	d.DeleteFunc = func(ctx context.Context, c *client.Client, v Values) error {
		// Verify existence separately. A 404 from the unmounted action endpoint
		// must never be interpreted as successful demotion or user removal.
		if _, err := Call(ctx, c, d.Read, v, nil); err != nil {
			return err
		}
		return serverAdminActionError(identityWrite(ctx, c, d.Update, v, Object{"serverAdmin": false}))
	}
	return d
}

func serverAdminActionError(err error) error {
	if client.IsNotFound(err) {
		return fmt.Errorf("server administrator action is unavailable; stock Pangolin 1.24.0 does not mount this action in the Bearer-key Integration API: %v", err)
	}
	return err
}

func clientDefinition(user bool) Definition {
	createPath, apiPath, name := "/org/{org_id}/client", "/org/{orgId}/client", "client"
	if user {
		createPath, apiPath, name = "/org/{org_id}/user/{user_id}/client", "/org/{orgId}/user/{userId}/client", "user_client"
	}
	d := Definition{
		Name: name, Description: "A Pangolin OLM client. Credentials and address defaults are allocated only during creation. subnet is a host IP address, without a CIDR prefix.", IDField: "clientId",
		Fields: SpecFields(apiPath, "put", "/client/{clientId}", "post"),
		Create: Endpoint{Method: "PUT", Path: createPath}, Read: Endpoint{Method: "GET", Path: "/client/{id}"}, Update: Endpoint{Method: "POST", Path: "/client/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/client/{id}"},
		ImportFields: []string{"org_id", "id"},
	}
	d.Fields["org_id"] = identityOrgField()
	SetField(&d, "type", func(f *Field) { f.Required = false; f.Default = "olm" })
	d.Fields["archived"] = WithDefault(Local(Bool("archived")), false)
	d.Fields["blocked"] = WithDefault(Local(Bool("blocked")), false)
	SetField(&d, "archived", func(f *Field) {
		f.Preserve = false
		f.Description = "Whether the client is archived. Defaults to false."
	})
	SetField(&d, "blocked", func(f *Field) {
		f.Preserve = false
		f.Description = "Whether the client is blocked. Blocking terminates connected OLM sessions. Defaults to false."
	})
	SetField(&d, "subnet", func(f *Field) {
		f.Required = false
		f.Description = "Host IP address without CIDR prefix. Allocated from the organization only on create when omitted."
	})
	if user {
		d.Description += " In Pangolin 1.24.0, the user-client CREATE route is commented out in the Integration API, and the shared DELETE handler rejects clients owned by users. Existing user clients can be imported, read, updated, archived, and blocked through the shared client routes, but create, replacement, and destroy require a server that implements these operations. There is no Bearer-key alternative for the missing lifecycle operations. Mock create/destroy tests verify specification adapter behavior only."
		SetField(&d, "name", func(f *Field) {
			f.Description = "Client display name. Pangolin 1.24.0 replaces this with a recognized Mac/iOS model name for connected user devices; use that returned display name when importing those devices because custom names cannot be read back reliably."
		})
		d.Fields["user_id"] = Required(Immutable(Local(String("userId"))))
		SetField(&d, "olm_id", func(f *Field) {
			f.Description = "Existing OLM ID belonging to user_id. The provider does not generate an OLM for user clients."
			f.Preserve = true
		})
	} else {
		SetField(&d, "olm_id", func(f *Field) { f.Required = false; f.Preserve = true })
		SetField(&d, "secret", func(f *Field) { f.Required = false; f.Sensitive = true; f.Preserve = true })
		d.List = Endpoint{Method: "GET", Path: "/org/{org_id}/clients", Key: "clients"}
		d.Pagination = "page"
		d.Sources = []string{"client", "clients"}
	}
	d.Validate = func(v Values) error {
		if ip := Text(v["subnet"]); ip != "" && net.ParseIP(ip) == nil {
			return fmt.Errorf("subnet must be a host IP address without a CIDR prefix")
		}
		return nil
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		values := Values{}
		for k, value := range v {
			values[k] = value
		}
		needsDefaults := values["subnet"] == nil || (!user && (values["olm_id"] == nil || values["secret"] == nil))
		if needsDefaults {
			defaults, err := Call(ctx, c, Endpoint{Method: "GET", Path: "/org/{org_id}/pick-client-defaults"}, v, nil)
			if err != nil {
				return nil, err
			}
			for name, api := range map[string]string{"olm_id": "olmId", "secret": "olmSecret", "subnet": "subnet"} {
				if user && name != "subnet" {
					continue
				}
				if values[name] == nil {
					values[name] = defaults[api]
				}
			}
		}
		body, err := Body(d, values, true)
		if err != nil {
			return nil, err
		}
		obj, err := Call(ctx, c, d.Create, values, body)
		if err != nil {
			if user && client.IsNotFound(err) {
				return nil, fmt.Errorf("user-client creation is unavailable; stock Pangolin 1.24.0 does not mount this route in the Integration API: %v", err)
			}
			return nil, err
		}
		// Credentials are absent from subsequent GET responses.
		if !user {
			obj["secret"] = values["secret"]
		}
		if obj["olmId"] == nil {
			obj["olmId"] = values["olm_id"]
		}
		if v["nice_id"] != nil {
			obj["niceId"] = v["nice_id"]
		}
		if err = d.EnrichRead(ctx, c, values, obj); err != nil {
			return obj, err
		}
		return obj, nil
	}
	d.EnrichRead = func(_ context.Context, _ *client.Client, v Values, obj Object) error {
		if subnet := Text(obj["subnet"]); subnet != "" {
			obj["subnet"] = strings.SplitN(subnet, "/", 2)[0]
		}
		// The current GET handler joins an active OLM and can return null before
		// that OLM connects, despite the immutable association existing.
		if obj["olmId"] == nil && v["olm_id"] != nil {
			obj["olmId"] = v["olm_id"]
		}
		return nil
	}
	d.AfterWrite = func(ctx context.Context, c *client.Client, v, prior Values) error {
		if prior == nil && v["nice_id"] != nil {
			if err := identityWrite(ctx, c, d.Update, v, Object{"niceId": v["nice_id"]}); err != nil {
				return err
			}
		}
		return reconcileClientState(ctx, c, v, prior)
	}
	return d
}

func reconcileClientState(ctx context.Context, c *client.Client, v, prior Values) error {
	if prior != nil && !identityChanged(prior, v, "archived") && !identityChanged(prior, v, "blocked") {
		return nil
	}
	current, err := Call(ctx, c, Endpoint{Method: "GET", Path: "/client/{id}"}, v, nil)
	if err != nil {
		return err
	}
	// These action endpoints reject repeats rather than being idempotent. Read
	// current flags before each reconciliation, including a partially failed retry.
	for _, state := range []struct{ field, enabled, disabled string }{{"blocked", "block", "unblock"}, {"archived", "archive", "unarchive"}} {
		desired, ok := v[state.field].(bool)
		if !ok {
			desired = false
		}
		actual, ok := current[state.field].(bool)
		if !ok {
			return fmt.Errorf("client response is missing %s state", state.field)
		}
		if actual == desired {
			continue
		}
		action := state.disabled
		if desired {
			action = state.enabled
		}
		if err := identityWrite(ctx, c, Endpoint{Method: "POST", Path: "/client/{id}/" + action}, v, nil); err != nil {
			return fmt.Errorf("setting client %s: %w", state.field, err)
		}
	}
	return nil
}

func oidcDefinition(org bool) Definition {
	cp, up, rp, lp := "/idp/oidc", "/idp/{idpId}/oidc", "/idp/{id}", "/idp"
	name, description := "global_oidc_identity_provider", "A global OIDC identity provider. Requires Pangolin global identity provider mode."
	imports := []string{"id"}
	if org {
		cp, up, rp, lp = "/org/{orgId}/idp/oidc", "/org/{orgId}/idp/{idpId}/oidc", "/org/{org_id}/idp/{id}", "/org/{org_id}/idp"
		name, description = "organization_oidc_identity_provider", "An organization-scoped OIDC identity provider. Requires the corresponding Pangolin edition and organization identity provider mode."
		imports = []string{"org_id", "id"}
	}
	d := Definition{Name: name, Description: description, IDField: "idpId", Fields: SpecFields(cp, "put", up, "post"), Read: Endpoint{Method: "GET", Path: rp}, Delete: Endpoint{Method: "DELETE", Path: rp}, List: Endpoint{Method: "GET", Path: lp, Key: "idps"}, Pagination: "offset", ImportFields: imports, Sources: []string{name, name + "s"}}
	d.Create = Endpoint{Method: "PUT", Path: strings.ReplaceAll(cp, "{orgId}", "{org_id}")}
	d.Update = Endpoint{Method: "POST", Path: rp + "/oidc"}
	if org {
		d.Fields["org_id"] = identityOrgField()
		// listOrgIdps filters rows by organization but counts the entire idp
		// table. A short page, not that unrelated total, completes this listing.
		d.Pagination = "offset-short"
	}
	SetField(&d, "client_secret", func(f *Field) { f.Sensitive = true })
	d.Fields["redirect_url"] = Preserve(Computed(String("redirectUrl")))
	d.EnrichRead = func(_ context.Context, _ *client.Client, _ Values, obj Object) error {
		return flattenIdentityProvider(obj)
	}
	if !org {
		d.AfterWrite = func(ctx context.Context, c *client.Client, v, prior Values) error {
			if prior != nil {
				return nil
			}
			body := Object{}
			for _, name := range []string{"default_role_mapping", "default_org_mapping"} {
				if v[name] != nil {
					body[d.Fields[name].API] = v[name]
				}
			}
			if len(body) == 0 {
				return nil
			}
			return identityWrite(ctx, c, d.Update, v, body)
		}
	}
	return d
}

func flattenIdentityProvider(obj Object) error {
	for _, key := range []string{"idp", "idpOidcConfig", "idpOrg"} {
		value, exists := obj[key]
		if !exists || value == nil {
			continue
		}
		nested, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("identity provider response has invalid %s object", key)
		}
		for name, field := range nested {
			obj[name] = field
		}
	}
	if obj["idpId"] == nil {
		return fmt.Errorf("identity provider response is missing idpId")
	}
	return nil
}

func identityProviderOrganizationDefinition() Definition {
	d := Definition{
		Name: "identity_provider_organization", Description: "A global identity provider's organization association and mappings. Do not manage the same association through an organization-scoped identity provider resource.", IDField: "orgId",
		Fields: SpecFields("/idp/{idpId}/org/{orgId}", "put", "/idp/{idpId}/org/{orgId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/idp/{idp_id}/org/{org_id}"}, Update: Endpoint{Method: "POST", Path: "/idp/{idp_id}/org/{org_id}"}, Delete: Endpoint{Method: "DELETE", Path: "/idp/{idp_id}/org/{org_id}"},
		List: Endpoint{Method: "GET", Path: "/idp/{idp_id}/org", Key: "policies"}, Pagination: "offset", ImportFields: []string{"idp_id", "org_id"},
	}
	d.Fields["idp_id"] = Required(Immutable(Local(ID("idpId"))))
	d.Fields["org_id"] = identityOrgField()
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		if err = identityWrite(ctx, c, d.Create, v, body); err != nil {
			return nil, err
		}
		return Object{"orgId": v["org_id"], "idpId": v["idp_id"]}, nil
	}
	d.ImportFunc = func(value string) (Values, error) {
		parts := strings.Split(value, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("expected idp_id/org_id")
		}
		return Values{"idp_id": parts[0], "org_id": parts[1], "id": parts[1]}, nil
	}
	d.ReadFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		path, err := Path(d.List.Path, v)
		if err != nil {
			return nil, err
		}
		rows, err := c.List(ctx, path, "policies", "offset")
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			if Text(row["orgId"]) == Text(v["org_id"]) {
				return row, nil
			}
		}
		return nil, client.NotFound(path)
	}
	return d
}

func identityChanged(prior, values Values, key string) bool {
	return prior == nil || !reflect.DeepEqual(prior[key], values[key])
}

func identityWrite(ctx context.Context, c *client.Client, endpoint Endpoint, values Values, body any) error {
	path, err := Path(endpoint.Path, values)
	if err != nil {
		return err
	}
	return c.Do(ctx, endpoint.Method, path, body, nil)
}
