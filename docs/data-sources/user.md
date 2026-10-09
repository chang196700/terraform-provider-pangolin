---
page_title: "pangolin_user Data Source - pangolin"
description: |-
  Reads global user metadata by ID. This server-level API requires a root API key with getUser permission.
---

# pangolin_user (Data Source)

Reads global user metadata by ID. This server-level API requires a root API key with getUser permission.

## Example Usage

```terraform
data "pangolin_user" "example" {
  user_id = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `date_created` | String | Read-only |  |
| `email` | String | Read-only |  |
| `email_verified` | Bool | Read-only |  |
| `idp_id` | String (numeric API ID) | Read-only |  |
| `idp_name` | String | Read-only |  |
| `name` | String | Read-only |  |
| `server_admin` | Bool | Read-only |  |
| `two_factor_enabled` | Bool | Read-only |  |
| `two_factor_setup_requested` | Bool | Read-only |  |
| `type` | String | Read-only |  |
| `user_id` | String | Required |  Retained in state when the API does not return it. |
| `username` | String | Read-only |  |
