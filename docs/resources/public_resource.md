---
page_title: "pangolin_public_resource Resource - pangolin"
description: |-
  A public resource. The access object owns inline authentication and authorization, and must not be combined with resource_policy_id because Pangolin 1.24.0 resource setters can modify that shared policy.
---

# pangolin_public_resource (Resource)

A public resource. The access object owns inline authentication and authorization, and must not be combined with resource_policy_id because Pangolin 1.24.0 resource setters can modify that shared policy.

## Example Usage

```terraform
resource "pangolin_public_resource" "example" {
  domain_id = "your-domain-id"
  mode      = "http"
  name      = "example"
  org_id    = "example"
  subdomain = "app"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access` | Object | Optional | Authoritative inline access settings. Omit this object when resource_policy_id is set. Removing it resets the inline policy to SSO with no explicit users, roles, whitelist, credentials, or rules. Secrets are stored in Terraform state and cannot be recovered on import. Retained in state when the API does not return it. |
| `ai_models` | Set of Object | Optional | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `ai_providers` | Set of Object | Optional | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `auth_daemon_mode` | String | Optional |  Allowed: `site`, `remote`, `native`. |
| `auth_daemon_port` | Int | Optional |  |
| `block_access` | Bool | Optional |  |
| `domain_id` | String | Optional |  |
| `enabled` | Bool | Optional |  |
| `id` | String | Read-only |  |
| `maintenance_estimated_time` | String | Optional |  |
| `maintenance_message` | String | Optional |  |
| `maintenance_mode_enabled` | Bool | Optional |  |
| `maintenance_mode_type` | String | Optional |  Allowed: `forced`, `automatic`. |
| `maintenance_title` | String | Optional |  |
| `mode` | String | Required, replacement |  Allowed: `http`, `ssh`, `rdp`, `vnc`, `tcp`, `udp`, `inference`. |
| `name` | String | Required |  |
| `nice_id` | String | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pam_mode` | String | Optional |  Allowed: `passthrough`, `push`. |
| `post_auth_path` | String | Optional |  |
| `proxy_port` | Int | Optional |  |
| `proxy_protocol` | Bool | Optional |  |
| `proxy_protocol_version` | Int | Optional |  |
| `request_headers` | List of Object | Optional, Sensitive |  |
| `resource_policy_id` | String (numeric API ID) | Optional |  |
| `response_headers` | List of Object | Optional, Sensitive |  |
| `set_host_header` | String | Optional |  |
| `ssl` | Bool | Optional |  |
| `sticky_session` | Bool | Optional |  |
| `subdomain` | String | Optional |  |
| `tls_server_name` | String | Optional |  |

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

## Import

```shell
terraform import pangolin_public_resource.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
