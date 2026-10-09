---
page_title: "pangolin_event_streaming_destination Resource - pangolin"
description: |-
  Manages an organization event streaming destination. The opaque destination configuration is sensitive because it may contain authentication credentials.
---

# pangolin_event_streaming_destination (Resource)

Manages an organization event streaming destination. The opaque destination configuration is sensitive because it may contain authentication credentials.

## Example Usage

```terraform
resource "pangolin_event_streaming_destination" "example" {
  config = "{}"
  org_id = "example"
  type   = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `config` | String | Required, Sensitive |  |
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Optional |  Default: `true`. |
| `id` | String | Read-only |  |
| `last_error` | String | Read-only |  |
| `last_error_at` | Int | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `send_access_logs` | Bool | Optional |  Default: `false`. |
| `send_action_logs` | Bool | Optional |  Default: `false`. |
| `send_ai_session_logs` | Bool | Optional |  Default: `false`. |
| `send_connection_logs` | Bool | Optional |  Default: `false`. |
| `send_request_logs` | Bool | Optional |  Default: `false`. |
| `type` | String | Required |  |
| `updated_at` | Int | Read-only |  |

## Import

```shell
terraform import pangolin_event_streaming_destination.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
