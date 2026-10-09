---
page_title: "pangolin_ai_budgets Data Source - pangolin"
description: |-
  Lists ai budget objects across all API pages. Secret fields are excluded.
---

# pangolin_ai_budgets (Data Source)

Lists ai budget objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_ai_budgets" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].amount` | Number | Required |  |
| `items[].created_at` | Int | Read-only |  |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].enforcement` | String | Optional |  Default: `"hard"`. Allowed: `hard`, `soft`. |
| `items[].id` | String | Read-only |  |
| `items[].model_id` | String (numeric API ID) | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].period` | String | Optional |  Default: `"monthly"`. Allowed: `monthly`, `yearly`, `lifetime`, `daily`, `hourly`, `weekly`. |
| `items[].provider_id` | String (numeric API ID) | Optional |  |
| `items[].resource_id` | String (numeric API ID) | Optional |  |
| `items[].role_id` | String (numeric API ID) | Optional |  |
| `items[].site_resource_id` | String (numeric API ID) | Optional |  |
| `items[].unit` | String | Required |  Allowed: `usd`, `tokens`. |
| `items[].updated_at` | Int | Read-only |  |
| `items[].virtual_api_key_id` | String | Optional |  |
