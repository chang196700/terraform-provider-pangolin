---
page_title: "pangolin_roles Data Source - pangolin"
description: |-
  Lists role objects across all API pages. Secret fields are excluded.
---

# pangolin_roles (Data Source)

Lists role objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_roles" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].allow_ssh` | Bool | Optional | Whether to allow SSH signing. Read from the organization's role list because the individual role endpoint omits this permission. |
| `items[].description` | String | Optional |  |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].reassign_role_id` | String (numeric API ID) | Required | Existing role in the same organization to receive users when this role is deleted. Required even when the role currently has no members. Retained in state when the API does not return it. |
| `items[].require_device_approval` | Bool | Optional |  |
| `items[].ssh_create_home_dir` | Bool | Optional |  |
| `items[].ssh_sudo_commands` | Set of String | Optional |  Default: `[]`. |
| `items[].ssh_sudo_mode` | String | Optional |  Allowed: `none`, `full`, `commands`. |
| `items[].ssh_unix_groups` | Set of String | Optional |  Default: `[]`. |
