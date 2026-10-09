---
page_title: "pangolin_role Data Source - pangolin"
description: |-
  Reads an existing role.
---

# pangolin_role (Data Source)

Reads an existing role.

## Example Usage

```terraform
data "pangolin_role" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `allow_ssh` | Bool | Read-only | Whether to allow SSH signing. Read from the organization's role list because the individual role endpoint omits this permission. |
| `description` | String | Read-only |  |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `reassign_role_id` | String (numeric API ID) | Required | Existing role in the same organization to receive users when this role is deleted. Required even when the role currently has no members. Retained in state when the API does not return it. |
| `require_device_approval` | Bool | Read-only |  |
| `ssh_create_home_dir` | Bool | Read-only |  |
| `ssh_sudo_commands` | Set of String | Read-only |  Default: `[]`. |
| `ssh_sudo_mode` | String | Read-only |  Allowed: `none`, `full`, `commands`. |
| `ssh_unix_groups` | Set of String | Read-only |  Default: `[]`. |
