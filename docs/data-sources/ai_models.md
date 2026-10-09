---
page_title: "pangolin_ai_models Data Source - pangolin"
description: |-
  Lists ai model objects across all API pages. Secret fields are excluded.
---

# pangolin_ai_models (Data Source)

Lists ai model objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_ai_models" "example" {
  provider_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `provider_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].created_at` | Int | Read-only |  |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].id` | String | Read-only |  |
| `items[].list_type` | String | Optional |  Default: `"allow"`. Allowed: `allow`, `block`. |
| `items[].model_key` | String | Required |  |
| `items[].name` | String | Required |  |
| `items[].provider_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `items[].updated_at` | Int | Read-only |  |
