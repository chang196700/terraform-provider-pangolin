---
page_title: "pangolin_global_oidc_identity_provider Data Source - pangolin"
description: |-
  Reads an existing global oidc identity provider.
---

# pangolin_global_oidc_identity_provider (Data Source)

Reads an existing global oidc identity provider.

## Example Usage

```terraform
data "pangolin_global_oidc_identity_provider" "example" {
  id = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `auth_url` | String | Required |  |
| `auto_provision` | Bool | Read-only |  |
| `client_id` | String | Required |  |
| `client_secret` | String | Required, Sensitive |  |
| `default_org_mapping` | String | Read-only |  |
| `default_role_mapping` | String | Read-only |  |
| `email_path` | String | Read-only |  |
| `id` | String | Required |  |
| `identifier_path` | String | Required |  |
| `name` | String | Required |  |
| `name_path` | String | Read-only |  |
| `redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `scopes` | String | Required |  |
| `tags` | String | Read-only |  |
| `token_url` | String | Required |  |
| `variant` | String | Read-only |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |
