---
page_title: "pangolin_sites Data Source - pangolin"
description: |-
  Lists site objects across all API pages. Secret fields are excluded.
---

# pangolin_sites (Data Source)

Lists site objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_sites" "example" {
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
| `items[].address` | String | Optional, replacement |  |
| `items[].auto_update_enabled` | Bool | Optional |  |
| `items[].auto_update_override_org` | Bool | Optional |  |
| `items[].docker_socket_enabled` | Bool | Optional |  |
| `items[].exit_node_id` | String (numeric API ID) | Optional, replacement |  |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].newt_id` | String | Optional, replacement |  |
| `items[].nice_id` | String | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].pub_key` | String | Optional, replacement |  |
| `items[].subnet` | String | Optional, replacement |  |
| `items[].type` | String | Required, replacement |  Allowed: `newt`, `wireguard`, `local`. |
