---
page_title: "pangolin_organization_oidc_identity_providers Data Source - pangolin"
description: |-
  Lists organization oidc identity provider objects across all API pages. Secret fields are excluded.
---

# pangolin_organization_oidc_identity_providers (Data Source)

Lists organization oidc identity provider objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_organization_oidc_identity_providers" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].auth_url` | String | Required |  |
| `items[].auto_provision` | Bool | Optional |  |
| `items[].client_id` | String | Required |  |
| `items[].email_path` | String | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].identifier_path` | String | Required |  |
| `items[].name` | String | Required |  |
| `items[].name_path` | String | Optional |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].org_mapping` | String | Optional |  |
| `items[].redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `items[].role_mapping` | String | Optional |  |
| `items[].scopes` | String | Required |  |
| `items[].tags` | String | Optional |  |
| `items[].token_url` | String | Required |  |
| `items[].variant` | String | Optional, replacement |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |
