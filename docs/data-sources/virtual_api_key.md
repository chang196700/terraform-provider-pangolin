---
page_title: "pangolin_virtual_api_key Data Source - pangolin"
description: |-
  Reads an existing virtual api key.
---

# pangolin_virtual_api_key (Data Source)

Reads an existing virtual api key.

## Example Usage

```terraform
data "pangolin_virtual_api_key" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `all_resources` | Bool | Read-only |  Default: `false`. |
| `created_at` | Int | Read-only |  |
| `created_by_user_id` | String | Read-only |  |
| `description` | String | Read-only |  |
| `expires_at` | Int | Read-only |  |
| `id` | String | Required |  |
| `kind` | String | Read-only |  |
| `last_chars` | String | Read-only |  |
| `last_used_at` | Int | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `resource_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `secret` | String | Read-only, Sensitive |  |
| `user_id` | String | Read-only |  |
| `valid_for_seconds` | Int | Read-only | Lifetime from creation, or from the most recent change to this value. Removing it clears the expiry. This input is not returned by the API. Retained in state when the API does not return it. |
