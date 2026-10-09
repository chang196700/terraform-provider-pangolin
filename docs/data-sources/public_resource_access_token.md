---
page_title: "pangolin_public_resource_access_token Data Source - pangolin"
description: |-
  Reads an existing public resource access token.
---

# pangolin_public_resource_access_token (Data Source)

Reads an existing public resource access token.

## Example Usage

```terraform
data "pangolin_public_resource_access_token" "example" {
  id = "example-id"
  org_id = "example"
  resource_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `access_token` | String | Read-only, Sensitive |  Retained in state when the API does not return it. |
| `created_at` | Int | Read-only |  |
| `description` | String | Read-only, replacement |  |
| `expires_at` | Int | Read-only |  |
| `id` | String | Required |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `path` | String | Read-only, replacement | Optional request path restriction. The token list API omits this value, so it is retained locally and cannot be recovered on import. Retained in state when the API does not return it. |
| `persist_session` | Bool | Read-only, replacement |  Default: `false`. |
| `resource_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `title` | String | Read-only, replacement |  |
| `user_id` | String | Read-only, replacement |  |
| `valid_for_seconds` | Int | Read-only, replacement |  Retained in state when the API does not return it. |
