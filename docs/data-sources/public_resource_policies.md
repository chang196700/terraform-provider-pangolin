---
page_title: "pangolin_public_resource_policies Data Source - pangolin"
description: |-
  Lists public resource policy objects across all API pages. Secret fields are excluded.
---

# pangolin_public_resource_policies (Data Source)

Lists public resource policy objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_public_resource_policies" "example" {
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
| `items[].apply_rules` | Bool | Optional |  Default: `false`. |
| `items[].email_whitelist_enabled` | Bool | Optional |  Default: `false`. |
| `items[].emails` | Set of String | Optional |  Default: `[]`. |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].nice_id` | String | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].rules` | Set of Object | Optional |  Default: `[]`. |
| `items[].skip_to_idp_id` | String (numeric API ID) | Optional |  |
| `items[].sso` | Bool | Optional |  Default: `true`. |
| `items[].user_ids` | Set of String | Optional |  Default: `[]`. |

### Nested items[].rules

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].rules[].action` | String | Required |  Allowed: `ACCEPT`, `DROP`, `PASS`. |
| `items[].rules[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].rules[].match` | String | Required |  Allowed: `CIDR`, `IP`, `PATH`, `COUNTRY`, `COUNTRY_IS_NOT`, `ASN`, `REGION`, `METHOD`. |
| `items[].rules[].priority` | Int | Required |  |
| `items[].rules[].value` | String | Required |  |
