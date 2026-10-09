---
page_title: "pangolin_site Data Source - pangolin"
description: |-
  Reads an existing site.
---

# pangolin_site (Data Source)

Reads an existing site.

## Example Usage

```terraform
data "pangolin_site" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `address` | String | Read-only, replacement |  |
| `auto_update_enabled` | Bool | Read-only |  |
| `auto_update_override_org` | Bool | Read-only |  |
| `docker_socket_enabled` | Bool | Read-only |  |
| `exit_node_id` | String (numeric API ID) | Read-only, replacement |  |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `newt_id` | String | Read-only, replacement |  |
| `nice_id` | String | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `pub_key` | String | Read-only, replacement |  |
| `secret` | String | Read-only, Sensitive, replacement |  Retained in state when the API does not return it. |
| `subnet` | String | Read-only, replacement |  |
| `type` | String | Required, replacement |  Allowed: `newt`, `wireguard`, `local`. |
