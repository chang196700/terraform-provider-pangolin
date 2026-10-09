---
page_title: "pangolin_user_server_admin Resource - pangolin"
description: |-
  Specification adapter for administrator status. NOT AVAILABLE with the Pangolin 1.24.0 Integration API: its create, update, and destroy action is mounted only in the browser-session API, with no Bearer-key alternative. The root-key individual user read exists, but does not make this a usable lifecycle. Requires a server version or custom deployment that exposes the documented action. Mock lifecycle tests verify adapter behavior only. Destroy demotes the user rather than deleting it.
---

# pangolin_user_server_admin (Resource)

Specification adapter for administrator status. NOT AVAILABLE with the Pangolin 1.24.0 Integration API: its create, update, and destroy action is mounted only in the browser-session API, with no Bearer-key alternative. The root-key individual user read exists, but does not make this a usable lifecycle. Requires a server version or custom deployment that exposes the documented action. Mock lifecycle tests verify adapter behavior only. Destroy demotes the user rather than deleting it.

## Example Usage

```terraform
resource "pangolin_user_server_admin" "example" {
  server_admin = false
  user_id      = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `id` | String | Read-only |  |
| `server_admin` | Bool | Required |  |
| `user_id` | String | Required, replacement |  Retained in state when the API does not return it. |

## Import

```shell
terraform import pangolin_user_server_admin.example 'object-id'
```

Import identity: `id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
