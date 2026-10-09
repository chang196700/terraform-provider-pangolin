---
page_title: "pangolin_ai_provider Data Source - pangolin"
description: |-
  Reads an existing ai provider.
---

# pangolin_ai_provider (Data Source)

Reads an existing ai provider.

## Example Usage

```terraform
data "pangolin_ai_provider" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `api_key` | String | Read-only, Sensitive |  |
| `api_key_last_chars` | String | Read-only |  |
| `auth_type` | String | Read-only |  Allowed: `bearer`, `x-api-key`, `x-goog-api-key`, `hec`, `cf-aig-authorization`, `none`, `passthrough`. |
| `capabilities` | Set of String | Read-only |  |
| `created_at` | Int | Read-only |  |
| `effective_auth_type` | String | Read-only |  |
| `effective_upstream_url` | String | Read-only |  |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `headers` | List of Object | Read-only, Sensitive |  Default: `[]`. |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `nice_id` | String | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `routing_mode` | String | Read-only |  Default: `"url"`. Allowed: `url`, `target`. |
| `skip_tls_verification` | Bool | Read-only |  Default: `false`. |
| `type` | String | Required, replacement |  Allowed: `openai`, `anthropic`, `googleGemini`, `vertexAi`, `bedrock`, `microsoftFoundry`, `openRouter`, `vercelAiGateway`, `custom`. |
| `updated_at` | Int | Read-only |  |
| `upstream_url` | String | Read-only |  |

### Nested headers

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `headers[].name` | String | Required |  |
| `headers[].value` | String | Required |  |
