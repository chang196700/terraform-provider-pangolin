---
page_title: "pangolin_organization_oidc_identity_provider Data Source - pangolin"
description: |-
  Reads an existing organization oidc identity provider.
---

# pangolin_organization_oidc_identity_provider (Data Source)

Reads an existing organization oidc identity provider.

## Example Usage

```terraform
data "pangolin_organization_oidc_identity_provider" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `auth_url` | String | Required |  |
| `auto_provision` | Bool | Read-only |  |
| `client_id` | String | Required |  |
| `client_secret` | String | Required, Sensitive |  |
| `email_path` | String | Read-only |  |
| `id` | String | Required |  |
| `identifier_path` | String | Required |  |
| `name` | String | Required |  |
| `name_path` | String | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `org_mapping` | String | Read-only |  |
| `redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `role_mapping` | String | Read-only |  |
| `scopes` | String | Required |  |
| `tags` | String | Read-only |  |
| `token_url` | String | Required |  |
| `variant` | String | Read-only, replacement |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |
