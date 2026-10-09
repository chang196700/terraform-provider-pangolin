---
page_title: "pangolin_public_resource_policy Data Source - pangolin"
description: |-
  Reads an existing public resource policy.
---

# pangolin_public_resource_policy (Data Source)

Reads an existing public resource policy.

## Example Usage

```terraform
data "pangolin_public_resource_policy" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `apply_rules` | Bool | Read-only |  Default: `false`. |
| `email_whitelist_enabled` | Bool | Read-only |  Default: `false`. |
| `emails` | Set of String | Read-only |  Default: `[]`. |
| `header_auth` | Object | Read-only, Sensitive |  |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `nice_id` | String | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `password` | String | Read-only, Sensitive |  |
| `pincode` | String | Read-only, Sensitive |  |
| `role_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `rules` | Set of Object | Read-only |  Default: `[]`. |
| `skip_to_idp_id` | String (numeric API ID) | Read-only |  |
| `sso` | Bool | Read-only |  Default: `true`. |
| `user_ids` | Set of String | Read-only |  Default: `[]`. |

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
