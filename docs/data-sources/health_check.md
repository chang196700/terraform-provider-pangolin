---
page_title: "pangolin_health_check Data Source - pangolin"
description: |-
  Reads an existing health check.
---

# pangolin_health_check (Data Source)

Reads an existing health check.

## Example Usage

```terraform
data "pangolin_health_check" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `hc_enabled` | Bool | Read-only |  Default: `false`. |
| `hc_follow_redirects` | Bool | Read-only |  Default: `true`. |
| `hc_headers` | String | Read-only, Sensitive |  |
| `hc_health` | String | Read-only |  |
| `hc_healthy_threshold` | Int | Read-only |  Default: `1`. |
| `hc_hostname` | String | Read-only |  |
| `hc_interval` | Int | Read-only |  Default: `30`. |
| `hc_method` | String | Read-only |  Default: `"GET"`. |
| `hc_mode` | String | Read-only |  Default: `"http"`. |
| `hc_path` | String | Read-only |  |
| `hc_port` | Int | Read-only |  |
| `hc_scheme` | String | Read-only |  |
| `hc_status` | Int | Read-only |  |
| `hc_timeout` | Int | Read-only |  Default: `1`. |
| `hc_tls_server_name` | String | Read-only |  |
| `hc_unhealthy_interval` | Int | Read-only |  Default: `30`. |
| `hc_unhealthy_threshold` | Int | Read-only |  Default: `1`. |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `site_id` | String (numeric API ID) | Required |  |
| `site_name` | String | Read-only |  |
| `site_nice_id` | String | Read-only |  |
