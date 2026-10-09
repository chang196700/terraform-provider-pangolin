---
page_title: "pangolin_resource_targets Data Source - pangolin"
description: |-
  Lists resource target objects across all API pages. Secret fields are excluded.
---

# pangolin_resource_targets (Data Source)

Lists resource target objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_resource_targets" "example" {
  resource_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `resource_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].ai_provider_id` | String (numeric API ID) | Optional, replacement |  Retained in state when the API does not return it. |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].hc_enabled` | Bool | Optional |  |
| `items[].hc_follow_redirects` | Bool | Optional |  |
| `items[].hc_healthy_threshold` | Int | Optional |  |
| `items[].hc_hostname` | String | Optional |  |
| `items[].hc_interval` | Int | Optional |  |
| `items[].hc_method` | String | Optional |  |
| `items[].hc_mode` | String | Optional |  |
| `items[].hc_path` | String | Optional |  |
| `items[].hc_port` | Int | Optional |  |
| `items[].hc_scheme` | String | Optional |  |
| `items[].hc_status` | Int | Optional |  |
| `items[].hc_timeout` | Int | Optional |  |
| `items[].hc_tls_server_name` | String | Optional |  |
| `items[].hc_unhealthy_interval` | Int | Optional |  |
| `items[].hc_unhealthy_threshold` | Int | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].ip` | String | Required |  |
| `items[].method` | String | Optional |  |
| `items[].mode` | String | Optional |  Allowed: `http`, `tcp`, `udp`, `ssh`, `rdp`, `vnc`, ``. |
| `items[].path` | String | Optional |  |
| `items[].path_match_type` | String | Optional |  Allowed: `exact`, `prefix`, `regex`, ``. |
| `items[].port` | Int | Required |  |
| `items[].priority` | Int | Optional |  |
| `items[].resource_id` | String (numeric API ID) | Optional, replacement |  Retained in state when the API does not return it. |
| `items[].rewrite_path` | String | Optional |  |
| `items[].rewrite_path_type` | String | Optional |  Allowed: `exact`, `prefix`, `regex`, `stripPrefix`, ``. |
| `items[].site_id` | String (numeric API ID) | Required |  |
