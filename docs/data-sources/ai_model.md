---
page_title: "pangolin_ai_model Data Source - pangolin"
description: |-
  Reads an existing ai model.
---

# pangolin_ai_model (Data Source)

Reads an existing ai model.

## Example Usage

```terraform
data "pangolin_ai_model" "example" {
  id = "example-id"
  provider_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `id` | String | Required |  |
| `list_type` | String | Read-only |  Default: `"allow"`. Allowed: `allow`, `block`. |
| `model_key` | String | Required |  |
| `name` | String | Required |  |
| `provider_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `updated_at` | Int | Read-only |  |
