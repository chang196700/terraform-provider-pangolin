---
page_title: "pangolin_resource_target Data Source - pangolin"
description: |-
  Reads an existing resource target.
---

# pangolin_resource_target (Data Source)

Reads an existing resource target.

## Example Usage

```terraform
data "pangolin_resource_target" "example" {
  id = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_provider_id` | String (numeric API ID) | Read-only, replacement |  Retained in state when the API does not return it. |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `hc_enabled` | Bool | Read-only |  |
| `hc_follow_redirects` | Bool | Read-only |  |
| `hc_headers` | List of Object | Read-only, Sensitive |  |
| `hc_healthy_threshold` | Int | Read-only |  |
| `hc_hostname` | String | Read-only |  |
| `hc_interval` | Int | Read-only |  |
| `hc_method` | String | Read-only |  |
| `hc_mode` | String | Read-only |  |
| `hc_path` | String | Read-only |  |
| `hc_port` | Int | Read-only |  |
| `hc_scheme` | String | Read-only |  |
| `hc_status` | Int | Read-only |  |
| `hc_timeout` | Int | Read-only |  |
| `hc_tls_server_name` | String | Read-only |  |
| `hc_unhealthy_interval` | Int | Read-only |  |
| `hc_unhealthy_threshold` | Int | Read-only |  |
| `id` | String | Required |  |
| `ip` | String | Required |  |
| `method` | String | Read-only |  |
| `mode` | String | Read-only |  Allowed: `http`, `tcp`, `udp`, `ssh`, `rdp`, `vnc`, ``. |
| `path` | String | Read-only |  |
| `path_match_type` | String | Read-only |  Allowed: `exact`, `prefix`, `regex`, ``. |
| `port` | Int | Required |  |
| `priority` | Int | Read-only |  |
| `resource_id` | String (numeric API ID) | Read-only, replacement |  Retained in state when the API does not return it. |
| `rewrite_path` | String | Read-only |  |
| `rewrite_path_type` | String | Read-only |  Allowed: `exact`, `prefix`, `regex`, `stripPrefix`, ``. |
| `site_id` | String (numeric API ID) | Required |  |

### Nested hc_headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `hc_headers[].name` | String | Required |  |
| `hc_headers[].value` | String | Required |  |
