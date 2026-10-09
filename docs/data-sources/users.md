---
page_title: "pangolin_users Data Source - pangolin"
description: |-
  OpenAPI compatibility schema for listing global user metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires a server administrator session. Requires a server version or extension that mounts this route for API keys. A root API key or listUsers action does not enable a missing route.
---

# pangolin_users (Data Source)

OpenAPI compatibility schema for listing global user metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires a server administrator session. Requires a server version or extension that mounts this route for API keys. A root API key or listUsers action does not enable a missing route.

## Example Usage

```terraform
data "pangolin_users" "example" {
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `idp_id` | String | Optional | Filter by a numeric identity provider ID, or internal for internal users. |
| `items` | List of Object | Read-only |  |
| `query` | String | Optional |  Retained in state when the API does not return it. |
| `server_admin` | Bool | Optional |  Retained in state when the API does not return it. |
| `two_factor` | Bool | Optional |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].date_created` | String | Optional |  |
| `items[].email` | String | Optional |  |
| `items[].idp_id` | String (numeric API ID) | Optional |  |
| `items[].idp_name` | String | Optional |  |
| `items[].name` | String | Optional |  |
| `items[].server_admin` | Bool | Optional |  |
| `items[].two_factor_enabled` | Bool | Optional |  |
| `items[].two_factor_setup_requested` | Bool | Optional |  |
| `items[].type` | String | Optional |  |
| `items[].user_id` | String | Optional |  |
| `items[].username` | String | Optional |  |
