---
page_title: "pangolin_event_streaming_destination Data Source - pangolin"
description: |-
  Reads an existing event streaming destination.
---

# pangolin_event_streaming_destination (Data Source)

Reads an existing event streaming destination.

## Example Usage

```terraform
data "pangolin_event_streaming_destination" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `config` | String | Required, Sensitive |  |
| `created_at` | Int | Read-only |  |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `id` | String | Required |  |
| `last_error` | String | Read-only |  |
| `last_error_at` | Int | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `send_access_logs` | Bool | Read-only |  Default: `false`. |
| `send_action_logs` | Bool | Read-only |  Default: `false`. |
| `send_ai_session_logs` | Bool | Read-only |  Default: `false`. |
| `send_connection_logs` | Bool | Read-only |  Default: `false`. |
| `send_request_logs` | Bool | Read-only |  Default: `false`. |
| `type` | String | Required |  |
| `updated_at` | Int | Read-only |  |
