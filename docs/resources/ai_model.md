---
page_title: "pangolin_ai_model Resource - pangolin"
description: |-
  Manages a model allow or block entry on an AI provider. Model keys are unique within a provider.
---

# pangolin_ai_model (Resource)

Manages a model allow or block entry on an AI provider. Model keys are unique within a provider.

## Example Usage

```terraform
resource "pangolin_ai_model" "example" {
  model_key   = "example"
  name        = "example"
  provider_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Optional |  Default: `true`. |
| `id` | String | Read-only |  |
| `list_type` | String | Optional |  Default: `"allow"`. Allowed: `allow`, `block`. |
| `model_key` | String | Required |  |
| `name` | String | Required |  |
| `provider_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `updated_at` | Int | Read-only |  |

## Import

```shell
terraform import pangolin_ai_model.example '1/2'
```

Import identity: `provider_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
