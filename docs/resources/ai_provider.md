---
page_title: "pangolin_ai_provider Resource - pangolin"
description: |-
  Manages an AI gateway upstream provider. API keys and custom headers are stored as sensitive Terraform state.
---

# pangolin_ai_provider (Resource)

Manages an AI gateway upstream provider. API keys and custom headers are stored as sensitive Terraform state.

## Example Usage

```terraform
resource "pangolin_ai_provider" "example" {
  api_key      = "replace-with-your-upstream-key"
  name         = "example"
  org_id       = "example"
  type         = "openai"
  upstream_url = "https://api.openai.com/v1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `api_key` | String | Optional, Sensitive |  |
| `api_key_last_chars` | String | Read-only |  |
| `auth_type` | String | Optional |  Allowed: `bearer`, `x-api-key`, `x-goog-api-key`, `hec`, `cf-aig-authorization`, `none`, `passthrough`. |
| `capabilities` | Set of String | Optional |  |
| `created_at` | Int | Read-only |  |
| `effective_auth_type` | String | Read-only |  |
| `effective_upstream_url` | String | Read-only |  |
| `enabled` | Bool | Optional |  Default: `true`. |
| `headers` | List of Object | Optional, Sensitive |  Default: `[]`. |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `nice_id` | String | Optional |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `routing_mode` | String | Optional |  Default: `"url"`. Allowed: `url`, `target`. |
| `skip_tls_verification` | Bool | Optional |  Default: `false`. |
| `type` | String | Required, replacement |  Allowed: `openai`, `anthropic`, `googleGemini`, `vertexAi`, `bedrock`, `microsoftFoundry`, `openRouter`, `vercelAiGateway`, `custom`. |
| `updated_at` | Int | Read-only |  |
| `upstream_url` | String | Optional |  |

### Nested headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `headers[].name` | String | Required |  |
| `headers[].value` | String | Required |  |

## Import

```shell
terraform import pangolin_ai_provider.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
