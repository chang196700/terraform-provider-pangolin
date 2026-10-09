---
page_title: "pangolin_private_resource Data Source - pangolin"
description: |-
  Reads an existing private resource.
---

# pangolin_private_resource (Data Source)

Reads an existing private resource.

## Example Usage

```terraform
data "pangolin_private_resource" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_models` | Set of Object | Read-only | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `ai_providers` | Set of Object | Read-only | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `alias` | String | Read-only | Fully qualified domain name with optional wildcards, e.g., example.internal, *.example.internal, or host-0?.example.internal |
| `auth_daemon_mode` | String | Read-only |  Allowed: `site`, `remote`, `native`. |
| `auth_daemon_port` | Int | Read-only |  |
| `client_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `destination` | String | Read-only |  |
| `destination_port` | Int | Read-only, replacement |  |
| `disable_icmp` | Bool | Read-only |  |
| `domain_id` | String | Read-only |  |
| `enabled` | Bool | Read-only |  |
| `id` | String | Required |  |
| `mode` | String | Required, replacement |  Allowed: `host`, `cidr`, `http`, `ssh`, `inference`, `gateway`. |
| `name` | String | Required |  |
| `nice_id` | String | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pam_mode` | String | Read-only |  Allowed: `passthrough`, `push`. |
| `role_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `scheme` | String | Read-only |  Allowed: `http`, `https`, ``. |
| `site_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `ssl` | Bool | Read-only |  |
| `subdomain` | String | Read-only |  |
| `tcp_port_range_string` | String | Read-only | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `udp_port_range_string` | String | Read-only | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `user_ids` | Set of String | Read-only |  Default: `[]`. |

### Nested ai_models

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_models[].list_type` | String | Required |  Allowed: `allow`, `block`. |
| `ai_models[].model_id` | String (numeric API ID) | Required |  |

### Nested ai_providers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_providers[].access_mode` | String | Optional |  Default: `"inherit"`. Allowed: `inherit`, `select`. |
| `ai_providers[].enabled` | Bool | Optional |  Default: `true`. |
| `ai_providers[].provider_id` | String (numeric API ID) | Required |  |
