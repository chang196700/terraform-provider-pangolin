---
page_title: "pangolin_event_streaming_destinations Data Source - pangolin"
description: |-
  Lists event streaming destination objects across all API pages. Secret fields are excluded.
---

# pangolin_event_streaming_destinations (Data Source)

Lists event streaming destination objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_event_streaming_destinations" "example" {
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
| `items[].created_at` | Int | Read-only |  |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].id` | String | Read-only |  |
| `items[].last_error` | String | Read-only |  |
| `items[].last_error_at` | Int | Read-only |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].send_access_logs` | Bool | Optional |  Default: `false`. |
| `items[].send_action_logs` | Bool | Optional |  Default: `false`. |
| `items[].send_ai_session_logs` | Bool | Optional |  Default: `false`. |
| `items[].send_connection_logs` | Bool | Optional |  Default: `false`. |
| `items[].send_request_logs` | Bool | Optional |  Default: `false`. |
| `items[].type` | String | Required |  |
| `items[].updated_at` | Int | Read-only |  |
