---
page_title: "pangolin_ai_providers Data Source - pangolin"
description: |-
  Lists ai provider objects across all API pages. Secret fields are excluded.
---

# pangolin_ai_providers (Data Source)

Lists ai provider objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_ai_providers" "example" {
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
| `items[].api_key_last_chars` | String | Read-only |  |
| `items[].auth_type` | String | Optional |  Allowed: `bearer`, `x-api-key`, `x-goog-api-key`, `hec`, `cf-aig-authorization`, `none`, `passthrough`. |
| `items[].capabilities` | Set of String | Optional |  |
| `items[].created_at` | Int | Read-only |  |
| `items[].effective_auth_type` | String | Read-only |  |
| `items[].effective_upstream_url` | String | Read-only |  |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].nice_id` | String | Optional |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].routing_mode` | String | Optional |  Default: `"url"`. Allowed: `url`, `target`. |
| `items[].skip_tls_verification` | Bool | Optional |  Default: `false`. |
| `items[].type` | String | Required, replacement |  Allowed: `openai`, `anthropic`, `googleGemini`, `vertexAi`, `bedrock`, `microsoftFoundry`, `openRouter`, `vercelAiGateway`, `custom`. |
| `items[].updated_at` | Int | Read-only |  |
| `items[].upstream_url` | String | Optional |  |
