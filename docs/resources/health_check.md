---
page_title: "pangolin_health_check Resource - pangolin"
description: |-
  Manages a standalone site health check. Read refreshes the paginated organization list using targetHealthCheckId, the identifier returned by the API.
---

# pangolin_health_check (Resource)

Manages a standalone site health check. Read refreshes the paginated organization list using targetHealthCheckId, the identifier returned by the API.

## Example Usage

```terraform
resource "pangolin_health_check" "example" {
  hc_enabled  = true
  hc_hostname = "127.0.0.1"
  hc_mode     = "http"
  hc_port     = 8080
  name        = "example"
  org_id      = "example"
  site_id     = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `hc_enabled` | Bool | Optional |  Default: `false`. |
| `hc_follow_redirects` | Bool | Optional |  Default: `true`. |
| `hc_headers` | String | Optional, Sensitive |  |
| `hc_health` | String | Read-only |  |
| `hc_healthy_threshold` | Int | Optional |  Default: `1`. |
| `hc_hostname` | String | Optional |  |
| `hc_interval` | Int | Optional |  Default: `30`. |
| `hc_method` | String | Optional |  Default: `"GET"`. |
| `hc_mode` | String | Optional |  Default: `"http"`. |
| `hc_path` | String | Optional |  |
| `hc_port` | Int | Optional |  |
| `hc_scheme` | String | Optional |  |
| `hc_status` | Int | Optional |  |
| `hc_timeout` | Int | Optional |  Default: `1`. |
| `hc_tls_server_name` | String | Optional |  |
| `hc_unhealthy_interval` | Int | Optional |  Default: `30`. |
| `hc_unhealthy_threshold` | Int | Optional |  Default: `1`. |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `site_id` | String (numeric API ID) | Required |  |
| `site_name` | String | Read-only |  |
| `site_nice_id` | String | Read-only |  |

## Import

```shell
terraform import pangolin_health_check.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
