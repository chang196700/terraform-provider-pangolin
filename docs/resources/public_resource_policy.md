---
page_title: "pangolin_public_resource_policy Resource - pangolin"
description: |-
  A shared public-resource policy and its authoritative authentication, authorization, whitelist, and rules. Passwords and header credentials are write-only and cannot be recovered by import.
---

# pangolin_public_resource_policy (Resource)

A shared public-resource policy and its authoritative authentication, authorization, whitelist, and rules. Passwords and header credentials are write-only and cannot be recovered by import.

## Example Usage

```terraform
resource "pangolin_public_resource_policy" "example" {
  name   = "example"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `apply_rules` | Bool | Optional |  Default: `false`. |
| `email_whitelist_enabled` | Bool | Optional |  Default: `false`. |
| `emails` | Set of String | Optional |  Default: `[]`. |
| `header_auth` | Object | Optional, Sensitive |  |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `nice_id` | String | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `password` | String | Optional, Sensitive |  |
| `pincode` | String | Optional, Sensitive |  |
| `role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `rules` | Set of Object | Optional |  Default: `[]`. |
| `skip_to_idp_id` | String (numeric API ID) | Optional |  |
| `sso` | Bool | Optional |  Default: `true`. |
| `user_ids` | Set of String | Optional |  Default: `[]`. |

### Nested header_auth

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `header_auth.extended_compatibility` | Bool | Optional |  Default: `false`. |
| `header_auth.password` | String | Required, Sensitive |  |
| `header_auth.user` | String | Required |  |

### Nested rules

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `rules[].action` | String | Required |  Allowed: `ACCEPT`, `DROP`, `PASS`. |
| `rules[].enabled` | Bool | Optional |  Default: `true`. |
| `rules[].match` | String | Required |  Allowed: `CIDR`, `IP`, `PATH`, `COUNTRY`, `COUNTRY_IS_NOT`, `ASN`, `REGION`, `METHOD`. |
| `rules[].priority` | Int | Required |  |
| `rules[].value` | String | Required |  |

## Import

```shell
terraform import pangolin_public_resource_policy.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
