---
page_title: "pangolin_organization_api_keys Data Source - pangolin"
description: |-
  Lists organization api key objects across all API pages. Secret fields are excluded.
---

# pangolin_organization_api_keys (Data Source)

Lists organization api key objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_organization_api_keys" "example" {
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
| `items[].action_ids` | Set of String | Required | Nonempty authoritative set of API action identifiers. The server cannot replace this set with an empty set. |
| `items[].created_at` | String | Read-only |  |
| `items[].id` | String | Read-only |  |
| `items[].last_chars` | String | Read-only |  |
| `items[].name` | String | Required, replacement |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
