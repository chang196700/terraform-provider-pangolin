---
page_title: "pangolin_resource_target Resource - pangolin"
description: |-
  A public-resource or AI-provider target, routing rules, and target health checks. Set exactly one of resource_id and ai_provider_id. Import public targets as public/<resource-id>/<target-id> and AI targets as ai/<provider-id>/<target-id>.
---

# pangolin_resource_target (Resource)

A public-resource or AI-provider target, routing rules, and target health checks. Set exactly one of resource_id and ai_provider_id. Import public targets as public/<resource-id>/<target-id> and AI targets as ai/<provider-id>/<target-id>.

## Example Usage

```terraform
resource "pangolin_resource_target" "example" {
  ip          = "127.0.0.1"
  method      = "http"
  port        = 8080
  resource_id = "1"
  site_id     = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_provider_id` | String (numeric API ID) | Optional, replacement |  Retained in state when the API does not return it. |
| `enabled` | Bool | Optional |  Default: `true`. |
| `hc_enabled` | Bool | Optional |  |
| `hc_follow_redirects` | Bool | Optional |  |
| `hc_headers` | List of Object | Optional, Sensitive |  |
| `hc_healthy_threshold` | Int | Optional |  |
| `hc_hostname` | String | Optional |  |
| `hc_interval` | Int | Optional |  |
| `hc_method` | String | Optional |  |
| `hc_mode` | String | Optional |  |
| `hc_path` | String | Optional |  |
| `hc_port` | Int | Optional |  |
| `hc_scheme` | String | Optional |  |
| `hc_status` | Int | Optional |  |
| `hc_timeout` | Int | Optional |  |
| `hc_tls_server_name` | String | Optional |  |
| `hc_unhealthy_interval` | Int | Optional |  |
| `hc_unhealthy_threshold` | Int | Optional |  |
| `id` | String | Read-only |  |
| `ip` | String | Required |  |
| `method` | String | Optional |  |
| `mode` | String | Optional |  Allowed: `http`, `tcp`, `udp`, `ssh`, `rdp`, `vnc`, ``. |
| `path` | String | Optional |  |
| `path_match_type` | String | Optional |  Allowed: `exact`, `prefix`, `regex`, ``. |
| `port` | Int | Required |  |
| `priority` | Int | Optional |  |
| `resource_id` | String (numeric API ID) | Optional, replacement |  Retained in state when the API does not return it. |
| `rewrite_path` | String | Optional |  |
| `rewrite_path_type` | String | Optional |  Allowed: `exact`, `prefix`, `regex`, `stripPrefix`, ``. |
| `site_id` | String (numeric API ID) | Required |  |

### Nested hc_headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `hc_headers[].name` | String | Required |  |
| `hc_headers[].value` | String | Required |  |

## Import

```shell
terraform import pangolin_resource_target.example 'public/1/2'
```

Import identity: `public/resource_id/target_id or ai/provider_id/target_id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

AI targets use `ai/provider_id/target_id`; public targets use `public/resource_id/target_id`.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
