---
page_title: "pangolin_private_resources Data Source - pangolin"
description: |-
  Lists private resource objects across all API pages. Secret fields are excluded.
---

# pangolin_private_resources (Data Source)

Lists private resource objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_private_resources" "example" {
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
| `items[].ai_models` | Set of Object | Optional | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `items[].ai_providers` | Set of Object | Optional | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `items[].alias` | String | Optional | Fully qualified domain name with optional wildcards, e.g., example.internal, *.example.internal, or host-0?.example.internal |
| `items[].auth_daemon_mode` | String | Optional |  Allowed: `site`, `remote`, `native`. |
| `items[].auth_daemon_port` | Int | Optional |  |
| `items[].client_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].destination` | String | Optional |  |
| `items[].destination_port` | Int | Optional, replacement |  |
| `items[].disable_icmp` | Bool | Optional |  |
| `items[].domain_id` | String | Optional |  |
| `items[].enabled` | Bool | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].mode` | String | Required, replacement |  Allowed: `host`, `cidr`, `http`, `ssh`, `inference`, `gateway`. |
| `items[].name` | String | Required |  |
| `items[].nice_id` | String | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].pam_mode` | String | Optional |  Allowed: `passthrough`, `push`. |
| `items[].role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].scheme` | String | Optional |  Allowed: `http`, `https`, ``. |
| `items[].site_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].ssl` | Bool | Optional |  |
| `items[].subdomain` | String | Optional |  |
| `items[].tcp_port_range_string` | String | Optional | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `items[].udp_port_range_string` | String | Optional | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `items[].user_ids` | Set of String | Optional |  Default: `[]`. |

### Nested items[].ai_models

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].ai_models[].list_type` | String | Required |  Allowed: `allow`, `block`. |
| `items[].ai_models[].model_id` | String (numeric API ID) | Required |  |

### Nested items[].ai_providers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].ai_providers[].access_mode` | String | Optional |  Default: `"inherit"`. Allowed: `inherit`, `select`. |
| `items[].ai_providers[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].ai_providers[].provider_id` | String (numeric API ID) | Required |  |
