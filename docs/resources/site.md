---
page_title: "pangolin_site Resource - pangolin"
description: |-
  A Newt, WireGuard, or local site. Newt credentials are returned only at creation; imported sites cannot recover their secret.
---

# pangolin_site (Resource)

A Newt, WireGuard, or local site. Newt credentials are returned only at creation; imported sites cannot recover their secret.

## Example Usage

```terraform
resource "pangolin_site" "example" {
  name   = "example"
  org_id = "example"
  type   = "newt"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `address` | String | Optional, replacement |  |
| `auto_update_enabled` | Bool | Optional |  |
| `auto_update_override_org` | Bool | Optional |  |
| `docker_socket_enabled` | Bool | Optional |  |
| `exit_node_id` | String (numeric API ID) | Optional, replacement |  |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `newt_id` | String | Optional, replacement |  |
| `nice_id` | String | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pub_key` | String | Optional, replacement |  |
| `secret` | String | Optional, Sensitive, replacement |  Retained in state when the API does not return it. |
| `subnet` | String | Optional, replacement |  |
| `type` | String | Required, replacement |  Allowed: `newt`, `wireguard`, `local`. |

## Import

```shell
terraform import pangolin_site.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
