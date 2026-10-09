---
page_title: "pangolin_public_resource_access_tokens Data Source - pangolin"
description: |-
  Lists public resource access token objects across all API pages. Secret fields are excluded.
---

# pangolin_public_resource_access_tokens (Data Source)

Lists public resource access token objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_public_resource_access_tokens" "example" {
  org_id = "example"
  resource_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `resource_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].created_at` | Int | Read-only |  |
| `items[].description` | String | Optional, replacement |  |
| `items[].expires_at` | Int | Read-only |  |
| `items[].id` | String | Read-only |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].path` | String | Optional, replacement | Optional request path restriction. The token list API omits this value, so it is retained locally and cannot be recovered on import. Retained in state when the API does not return it. |
| `items[].persist_session` | Bool | Optional, replacement |  Default: `false`. |
| `items[].resource_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `items[].title` | String | Optional, replacement |  |
| `items[].user_id` | String | Optional, replacement |  |
| `items[].valid_for_seconds` | Int | Optional, replacement |  Retained in state when the API does not return it. |
