---
page_title: "pangolin_organizations Data Source - pangolin"
description: |-
  Lists organization objects across all API pages. Secret fields are excluded.
---

# pangolin_organizations (Data Source)

Lists organization objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_organizations" "example" {
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].id` | String | Read-only |  |
| `items[].max_session_length_hours` | Number | Optional |  |
| `items[].name` | String | Required |  |
| `items[].org_id` | String | Required, replacement |  |
| `items[].require_two_factor` | Bool | Optional |  |
| `items[].settings_enable_global_newt_auto_update` | Bool | Optional |  |
| `items[].settings_log_retention_days_access` | Number | Optional |  |
| `items[].settings_log_retention_days_action` | Number | Optional |  |
| `items[].settings_log_retention_days_ai_sessions` | Number | Optional |  |
| `items[].settings_log_retention_days_connection` | Number | Optional |  |
| `items[].settings_log_retention_days_request` | Number | Optional |  |
| `items[].subnet` | String | Required, replacement |  |
| `items[].utility_subnet` | String | Required, replacement |  |
