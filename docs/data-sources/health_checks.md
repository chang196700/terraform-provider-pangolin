---
page_title: "pangolin_health_checks Data Source - pangolin"
description: |-
  Lists health check objects across all API pages. Secret fields are excluded.
---

# pangolin_health_checks (Data Source)

Lists health check objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_health_checks" "example" {
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
| `items[].hc_enabled` | Bool | Optional |  Default: `false`. |
| `items[].hc_follow_redirects` | Bool | Optional |  Default: `true`. |
| `items[].hc_health` | String | Read-only |  |
| `items[].hc_healthy_threshold` | Int | Optional |  Default: `1`. |
| `items[].hc_hostname` | String | Optional |  |
| `items[].hc_interval` | Int | Optional |  Default: `30`. |
| `items[].hc_method` | String | Optional |  Default: `"GET"`. |
| `items[].hc_mode` | String | Optional |  Default: `"http"`. |
| `items[].hc_path` | String | Optional |  |
| `items[].hc_port` | Int | Optional |  |
| `items[].hc_scheme` | String | Optional |  |
| `items[].hc_status` | Int | Optional |  |
| `items[].hc_timeout` | Int | Optional |  Default: `1`. |
| `items[].hc_tls_server_name` | String | Optional |  |
| `items[].hc_unhealthy_interval` | Int | Optional |  Default: `30`. |
| `items[].hc_unhealthy_threshold` | Int | Optional |  Default: `1`. |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].site_id` | String (numeric API ID) | Required |  |
| `items[].site_name` | String | Read-only |  |
| `items[].site_nice_id` | String | Read-only |  |
