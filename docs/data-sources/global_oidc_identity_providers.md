---
page_title: "pangolin_global_oidc_identity_providers Data Source - pangolin"
description: |-
  Lists global oidc identity provider objects across all API pages. Secret fields are excluded.
---

# pangolin_global_oidc_identity_providers (Data Source)

Lists global oidc identity provider objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_global_oidc_identity_providers" "example" {
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].auth_url` | String | Required |  |
| `items[].auto_provision` | Bool | Optional |  |
| `items[].client_id` | String | Required |  |
| `items[].default_org_mapping` | String | Optional |  |
| `items[].default_role_mapping` | String | Optional |  |
| `items[].email_path` | String | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].identifier_path` | String | Required |  |
| `items[].name` | String | Required |  |
| `items[].name_path` | String | Optional |  |
| `items[].redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `items[].scopes` | String | Required |  |
| `items[].tags` | String | Optional |  |
| `items[].token_url` | String | Required |  |
| `items[].variant` | String | Optional |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |
