---
page_title: "pangolin_ai_budget Resource - pangolin"
description: |-
  Manages an AI usage budget. Set at most one scope ID; no scope ID selects the organization. A scope, unit and period combination must be unique.
---

# pangolin_ai_budget (Resource)

Manages an AI usage budget. Set at most one scope ID; no scope ID selects the organization. A scope, unit and period combination must be unique.

## Example Usage

```terraform
resource "pangolin_ai_budget" "example" {
  amount = 1000000
  org_id = "example"
  period = "monthly"
  unit   = "tokens"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `amount` | Number | Required |  |
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Optional |  Default: `true`. |
| `enforcement` | String | Optional |  Default: `"hard"`. Allowed: `hard`, `soft`. |
| `id` | String | Read-only |  |
| `model_id` | String (numeric API ID) | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `period` | String | Optional |  Default: `"monthly"`. Allowed: `monthly`, `yearly`, `lifetime`, `daily`, `hourly`, `weekly`. |
| `provider_id` | String (numeric API ID) | Optional |  |
| `resource_id` | String (numeric API ID) | Optional |  |
| `role_id` | String (numeric API ID) | Optional |  |
| `site_resource_id` | String (numeric API ID) | Optional |  |
| `unit` | String | Required |  Allowed: `usd`, `tokens`. |
| `updated_at` | Int | Read-only |  |
| `virtual_api_key_id` | String | Optional |  |

## Import

```shell
terraform import pangolin_ai_budget.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
