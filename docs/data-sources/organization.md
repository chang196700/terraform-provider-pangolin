---
page_title: "pangolin_organization Data Source - pangolin"
description: |-
  Reads an existing organization.
---

# pangolin_organization (Data Source)

Reads an existing organization.

## Example Usage

```terraform
data "pangolin_organization" "example" {
  id = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `id` | String | Required |  |
| `max_session_length_hours` | Number | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  |
| `password_expiry_days` | Number | Read-only, Sensitive |  |
| `require_two_factor` | Bool | Read-only |  |
| `settings_enable_global_newt_auto_update` | Bool | Read-only |  |
| `settings_log_retention_days_access` | Number | Read-only |  |
| `settings_log_retention_days_action` | Number | Read-only |  |
| `settings_log_retention_days_ai_sessions` | Number | Read-only |  |
| `settings_log_retention_days_connection` | Number | Read-only |  |
| `settings_log_retention_days_request` | Number | Read-only |  |
| `subnet` | String | Required, replacement |  |
| `utility_subnet` | String | Required, replacement |  |
