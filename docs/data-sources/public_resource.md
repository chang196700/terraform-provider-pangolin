---
page_title: "pangolin_public_resource Data Source - pangolin"
description: |-
  Reads an existing public resource.
---

# pangolin_public_resource (Data Source)

Reads an existing public resource.

## Example Usage

```terraform
data "pangolin_public_resource" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access` | Object | Read-only | Authoritative inline access settings. Omit this object when resource_policy_id is set. Removing it resets the inline policy to SSO with no explicit users, roles, whitelist, credentials, or rules. Secrets are stored in Terraform state and cannot be recovered on import. Retained in state when the API does not return it. |
| `ai_models` | Set of Object | Read-only | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `ai_providers` | Set of Object | Read-only | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `auth_daemon_mode` | String | Read-only |  Allowed: `site`, `remote`, `native`. |
| `auth_daemon_port` | Int | Read-only |  |
| `block_access` | Bool | Read-only |  |
| `domain_id` | String | Read-only |  |
| `enabled` | Bool | Read-only |  |
| `id` | String | Required |  |
| `maintenance_estimated_time` | String | Read-only |  |
| `maintenance_message` | String | Read-only |  |
| `maintenance_mode_enabled` | Bool | Read-only |  |
| `maintenance_mode_type` | String | Read-only |  Allowed: `forced`, `automatic`. |
| `maintenance_title` | String | Read-only |  |
| `mode` | String | Required, replacement |  Allowed: `http`, `ssh`, `rdp`, `vnc`, `tcp`, `udp`, `inference`. |
| `name` | String | Required |  |
| `nice_id` | String | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pam_mode` | String | Read-only |  Allowed: `passthrough`, `push`. |
| `post_auth_path` | String | Read-only |  |
| `proxy_port` | Int | Read-only |  |
| `proxy_protocol` | Bool | Read-only |  |
| `proxy_protocol_version` | Int | Read-only |  |
| `request_headers` | List of Object | Read-only, Sensitive |  |
| `resource_policy_id` | String (numeric API ID) | Read-only |  |
| `response_headers` | List of Object | Read-only, Sensitive |  |
| `set_host_header` | String | Read-only |  |
| `ssl` | Bool | Read-only |  |
| `sticky_session` | Bool | Read-only |  |
| `subdomain` | String | Read-only |  |
| `tls_server_name` | String | Read-only |  |

### Nested access

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access.apply_rules` | Bool | Optional |  Default: `false`. |
| `access.email_whitelist_enabled` | Bool | Optional |  Default: `false`. |
| `access.emails` | Set of String | Optional |  Default: `[]`. |
| `access.header_auth` | Object | Optional, Sensitive |  |
| `access.password` | String | Optional, Sensitive |  |
| `access.pincode` | String | Optional, Sensitive |  |
| `access.role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `access.rules` | Set of Object | Optional |  Default: `[]`. |
| `access.skip_to_idp_id` | String (numeric API ID) | Optional |  |
| `access.sso` | Bool | Optional |  Default: `true`. |
| `access.user_ids` | Set of String | Optional |  Default: `[]`. |

### Nested access.header_auth

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access.header_auth.extended_compatibility` | Bool | Optional |  Default: `false`. |
| `access.header_auth.password` | String | Required, Sensitive |  |
| `access.header_auth.user` | String | Required |  |

### Nested access.rules

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access.rules[].action` | String | Required |  Allowed: `ACCEPT`, `DROP`, `PASS`. |
| `access.rules[].enabled` | Bool | Optional |  Default: `true`. |
| `access.rules[].match` | String | Required |  Allowed: `CIDR`, `IP`, `PATH`, `COUNTRY`, `COUNTRY_IS_NOT`, `ASN`, `REGION`, `METHOD`. |
| `access.rules[].priority` | Int | Required |  |
| `access.rules[].value` | String | Required |  |

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

### Nested request_headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `request_headers[].name` | String | Required |  |
| `request_headers[].value` | String | Required |  |

### Nested response_headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `response_headers[].name` | String | Required |  |
| `response_headers[].value` | String | Required |  |
