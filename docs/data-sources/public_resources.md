---
page_title: "pangolin_public_resources Data Source - pangolin"
description: |-
  Lists public resource objects across all API pages. Secret fields are excluded.
---

# pangolin_public_resources (Data Source)

Lists public resource objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_public_resources" "example" {
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
| `items[].access` | Object | Optional | Authoritative inline access settings. Omit this object when resource_policy_id is set. Removing it resets the inline policy to SSO with no explicit users, roles, whitelist, credentials, or rules. Secrets are stored in Terraform state and cannot be recovered on import. Retained in state when the API does not return it. |
| `items[].ai_models` | Set of Object | Optional | Authoritative allow/block selections from attached providers in select mode. Default: `[]`. |
| `items[].ai_providers` | Set of Object | Optional | Authoritative AI provider attachments for inference mode. Providers are configured before model selections. Default: `[]`. |
| `items[].auth_daemon_mode` | String | Optional |  Allowed: `site`, `remote`, `native`. |
| `items[].auth_daemon_port` | Int | Optional |  |
| `items[].block_access` | Bool | Optional |  |
| `items[].domain_id` | String | Optional |  |
| `items[].enabled` | Bool | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].maintenance_estimated_time` | String | Optional |  |
| `items[].maintenance_message` | String | Optional |  |
| `items[].maintenance_mode_enabled` | Bool | Optional |  |
| `items[].maintenance_mode_type` | String | Optional |  Allowed: `forced`, `automatic`. |
| `items[].maintenance_title` | String | Optional |  |
| `items[].mode` | String | Required, replacement |  Allowed: `http`, `ssh`, `rdp`, `vnc`, `tcp`, `udp`, `inference`. |
| `items[].name` | String | Required |  |
| `items[].nice_id` | String | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].pam_mode` | String | Optional |  Allowed: `passthrough`, `push`. |
| `items[].post_auth_path` | String | Optional |  |
| `items[].proxy_port` | Int | Optional |  |
| `items[].proxy_protocol` | Bool | Optional |  |
| `items[].proxy_protocol_version` | Int | Optional |  |
| `items[].resource_policy_id` | String (numeric API ID) | Optional |  |
| `items[].set_host_header` | String | Optional |  |
| `items[].ssl` | Bool | Optional |  |
| `items[].sticky_session` | Bool | Optional |  |
| `items[].subdomain` | String | Optional |  |
| `items[].tls_server_name` | String | Optional |  |

### Nested items[].access

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].access.apply_rules` | Bool | Optional |  Default: `false`. |
| `items[].access.email_whitelist_enabled` | Bool | Optional |  Default: `false`. |
| `items[].access.emails` | Set of String | Optional |  Default: `[]`. |
| `items[].access.header_auth` | Object | Optional, Sensitive |  |
| `items[].access.password` | String | Optional, Sensitive |  |
| `items[].access.pincode` | String | Optional, Sensitive |  |
| `items[].access.role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].access.rules` | Set of Object | Optional |  Default: `[]`. |
| `items[].access.skip_to_idp_id` | String (numeric API ID) | Optional |  |
| `items[].access.sso` | Bool | Optional |  Default: `true`. |
| `items[].access.user_ids` | Set of String | Optional |  Default: `[]`. |

### Nested items[].access.header_auth

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].access.header_auth.extended_compatibility` | Bool | Optional |  Default: `false`. |
| `items[].access.header_auth.password` | String | Required, Sensitive |  |
| `items[].access.header_auth.user` | String | Required |  |

### Nested items[].access.rules

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].access.rules[].action` | String | Required |  Allowed: `ACCEPT`, `DROP`, `PASS`. |
| `items[].access.rules[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].access.rules[].match` | String | Required |  Allowed: `CIDR`, `IP`, `PATH`, `COUNTRY`, `COUNTRY_IS_NOT`, `ASN`, `REGION`, `METHOD`. |
| `items[].access.rules[].priority` | Int | Required |  |
| `items[].access.rules[].value` | String | Required |  |

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
