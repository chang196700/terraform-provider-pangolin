---
page_title: "pangolin_private_resource Resource - pangolin"
description: |-
  A private resource and its authoritative site, role, user, client, and inference-model associations. Changing mode or destination_port replaces the resource because Pangolin derives network, ICMP, and port-range values from them.
---

# pangolin_private_resource (Resource)

A private resource and its authoritative site, role, user, client, and inference-model associations. Changing mode or destination_port replaces the resource because Pangolin derives network, ICMP, and port-range values from them.

## Example Usage

```terraform
resource "pangolin_private_resource" "example" {
  destination           = "10.0.0.10"
  mode                  = "host"
  name                  = "example"
  org_id                = "example"
  site_ids              = ["1"]
  tcp_port_range_string = "443"
  udp_port_range_string = ""
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `ai_models` | Set of Object | Optional | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `ai_providers` | Set of Object | Optional | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `alias` | String | Optional | Fully qualified domain name with optional wildcards, e.g., example.internal, *.example.internal, or host-0?.example.internal |
| `auth_daemon_mode` | String | Optional |  Allowed: `site`, `remote`, `native`. |
| `auth_daemon_port` | Int | Optional |  |
| `client_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `destination` | String | Optional |  |
| `destination_port` | Int | Optional, replacement |  |
| `disable_icmp` | Bool | Optional |  |
| `domain_id` | String | Optional |  |
| `enabled` | Bool | Optional |  |
| `id` | String | Read-only |  |
| `mode` | String | Required, replacement |  Allowed: `host`, `cidr`, `http`, `ssh`, `inference`, `gateway`. |
| `name` | String | Required |  |
| `nice_id` | String | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pam_mode` | String | Optional |  Allowed: `passthrough`, `push`. |
| `role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `scheme` | String | Optional |  Allowed: `http`, `https`, ``. |
| `site_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `ssl` | Bool | Optional |  |
| `subdomain` | String | Optional |  |
| `tcp_port_range_string` | String | Optional | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `udp_port_range_string` | String | Optional | Port range string. Use "*" for all ports, a comma-separated list of ports, or ranges (e.g., "80,443,8000-9000"). Ports must be between 1 and 65535. |
| `user_ids` | Set of String | Optional |  Default: `[]`. |

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

## Import

```shell
terraform import pangolin_private_resource.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
