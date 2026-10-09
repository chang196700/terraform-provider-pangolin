---
page_title: "pangolin_organization_users Data Source - pangolin"
description: |-
  Lists organization user objects across all API pages. Secret fields are excluded.
---

# pangolin_organization_users (Data Source)

Lists organization user objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_organization_users" "example" {
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
| `items[].auto_provisioned` | Bool | Optional |  |
| `items[].email` | String | Optional, replacement |  |
| `items[].id` | String | Read-only |  |
| `items[].idp_id` | String (numeric API ID) | Required, replacement |  |
| `items[].name` | String | Optional, replacement |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].role_ids` | Set of String (numeric API ID) | Required | Nonempty authoritative role ID set. New assignments are added before obsolete ones are removed. |
| `items[].type` | String | Optional, replacement |  Default: `"oidc"`. Allowed: `oidc`. |
| `items[].username` | String | Required, replacement |  |
