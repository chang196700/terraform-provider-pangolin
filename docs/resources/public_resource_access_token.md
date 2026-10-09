---
page_title: "pangolin_public_resource_access_token Resource - pangolin"
description: |-
  A public resource access token. Settings are immutable, so changing them replaces the token. The one-time token is retained in sensitive state; import cannot recover it. Expired tokens disappear from the API list and will be recreated on the next apply.
---

# pangolin_public_resource_access_token (Resource)

A public resource access token. Settings are immutable, so changing them replaces the token. The one-time token is retained in sensitive state; import cannot recover it. Expired tokens disappear from the API list and will be recreated on the next apply.

## Example Usage

```terraform
resource "pangolin_public_resource_access_token" "example" {
  org_id      = "example"
  resource_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access_token` | String | Read-only, Sensitive |  Retained in state when the API does not return it. |
| `created_at` | Int | Read-only |  |
| `description` | String | Optional, replacement |  |
| `expires_at` | Int | Read-only |  |
| `id` | String | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `path` | String | Optional, replacement | Optional request path restriction. The token list API omits this value, so it is retained locally and cannot be recovered on import. Retained in state when the API does not return it. |
| `persist_session` | Bool | Optional, replacement |  Default: `false`. |
| `resource_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `title` | String | Optional, replacement |  |
| `user_id` | String | Optional, replacement |  |
| `valid_for_seconds` | Int | Optional, replacement |  Retained in state when the API does not return it. |

## Import

```shell
terraform import pangolin_public_resource_access_token.example 'example/1/2'
```

Import identity: `org_id/resource_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
