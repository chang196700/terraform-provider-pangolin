---
page_title: "pangolin_virtual_api_keys Data Source - pangolin"
description: |-
  Lists virtual api key objects across all API pages. Secret fields are excluded.
---

# pangolin_virtual_api_keys (Data Source)

Lists virtual api key objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_virtual_api_keys" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].all_resources` | Bool | Optional |  Default: `false`. |
| `items[].created_at` | Int | Read-only |  |
| `items[].created_by_user_id` | String | Read-only |  |
| `items[].description` | String | Optional |  |
| `items[].expires_at` | Int | Read-only |  |
| `items[].id` | String | Read-only |  |
| `items[].kind` | String | Read-only |  |
| `items[].last_chars` | String | Read-only |  |
| `items[].last_used_at` | Int | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].resource_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].user_id` | String | Optional |  |
| `items[].valid_for_seconds` | Int | Optional | Lifetime from creation, or from the most recent change to this value. Removing it clears the expiry. This input is not returned by the API. Retained in state when the API does not return it. |
