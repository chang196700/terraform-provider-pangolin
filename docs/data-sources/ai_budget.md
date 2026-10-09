---
page_title: "pangolin_ai_budget Data Source - pangolin"
description: |-
  Reads an existing ai budget.
---

# pangolin_ai_budget (Data Source)

Reads an existing ai budget.

## Example Usage

```terraform
data "pangolin_ai_budget" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `amount` | Number | Required |  |
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `enforcement` | String | Read-only |  Default: `"hard"`. Allowed: `hard`, `soft`. |
| `id` | String | Required |  |
| `model_id` | String (numeric API ID) | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `period` | String | Read-only |  Default: `"monthly"`. Allowed: `monthly`, `yearly`, `lifetime`, `daily`, `hourly`, `weekly`. |
| `provider_id` | String (numeric API ID) | Read-only |  |
| `resource_id` | String (numeric API ID) | Read-only |  |
| `role_id` | String (numeric API ID) | Read-only |  |
| `site_resource_id` | String (numeric API ID) | Read-only |  |
| `unit` | String | Required |  Allowed: `usd`, `tokens`. |
| `updated_at` | Int | Read-only |  |
| `virtual_api_key_id` | String | Read-only |  |
