---
page_title: "pangolin_role Resource - pangolin"
description: |-
  An organization role. Deletion migrates its users to reassign_role_id; choose a role with appropriate permissions before destroying this resource.
---

# pangolin_role (Resource)

An organization role. Deletion migrates its users to reassign_role_id; choose a role with appropriate permissions before destroying this resource.

## Example Usage

```terraform
resource "pangolin_role" "example" {
  name             = "example"
  org_id           = "example"
  reassign_role_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `allow_ssh` | Bool | Optional | Whether to allow SSH signing. Read from the organization's role list because the individual role endpoint omits this permission. |
| `description` | String | Optional |  |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `reassign_role_id` | String (numeric API ID) | Required | Existing role in the same organization to receive users when this role is deleted. Required even when the role currently has no members. Retained in state when the API does not return it. |
| `require_device_approval` | Bool | Optional |  |
| `ssh_create_home_dir` | Bool | Optional |  |
| `ssh_sudo_commands` | Set of String | Optional |  Default: `[]`. |
| `ssh_sudo_mode` | String | Optional |  Allowed: `none`, `full`, `commands`. |
| `ssh_unix_groups` | Set of String | Optional |  Default: `[]`. |

## Import

```shell
terraform import pangolin_role.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
