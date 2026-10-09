---
page_title: "pangolin_organization Resource - pangolin"
description: |-
  An organization and its network and security settings. Pangolin 1.24.0 has an upstream Integration API defect: organization deletion requires a session user and fails for API keys. Destroy retains state when the server rejects deletion; deleting organizations requires an upstream fix.
---

# pangolin_organization (Resource)

An organization and its network and security settings. Pangolin 1.24.0 has an upstream Integration API defect: organization deletion requires a session user and fails for API keys. Destroy retains state when the server rejects deletion; deleting organizations requires an upstream fix.

## Example Usage

```terraform
resource "pangolin_organization" "example" {
  name           = "example"
  org_id         = "example"
  subnet         = "100.90.0.0/16"
  utility_subnet = "100.91.0.0/16"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `id` | String | Read-only |  |
| `max_session_length_hours` | Number | Optional |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  |
| `password_expiry_days` | Number | Optional, Sensitive |  |
| `require_two_factor` | Bool | Optional |  |
| `settings_enable_global_newt_auto_update` | Bool | Optional |  |
| `settings_log_retention_days_access` | Number | Optional |  |
| `settings_log_retention_days_action` | Number | Optional |  |
| `settings_log_retention_days_ai_sessions` | Number | Optional |  |
| `settings_log_retention_days_connection` | Number | Optional |  |
| `settings_log_retention_days_request` | Number | Optional |  |
| `subnet` | String | Required, replacement |  |
| `utility_subnet` | String | Required, replacement |  |

## Import

```shell
terraform import pangolin_organization.example 'example'
```

Import identity: `id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
