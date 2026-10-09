---
page_title: "pangolin_alert_rule Resource - pangolin"
description: |-
  Manages an alert rule, event sources, email recipients and webhook actions. Source and recipient collections are authoritative. The provider reconciles the API's recipient rows and webhook configuration objects back into configuration fields.
---

# pangolin_alert_rule (Resource)

Manages an alert rule, event sources, email recipients and webhook actions. Source and recipient collections are authoritative. The provider reconciles the API's recipient rows and webhook configuration objects back into configuration fields.

## Example Usage

```terraform
resource "pangolin_alert_rule" "example" {
  all_sites  = true
  event_type = "site_offline"
  name       = "example"
  org_id     = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `all_health_checks` | Bool | Optional |  Default: `false`. |
| `all_resources` | Bool | Optional |  Default: `false`. |
| `all_sites` | Bool | Optional |  Default: `false`. |
| `cooldown_seconds` | Int | Optional |  Default: `0`. |
| `created_at` | Int | Read-only |  |
| `emails` | Set of String | Optional |  Default: `[]`. |
| `enabled` | Bool | Optional |  Default: `true`. |
| `event_type` | String | Required |  Allowed: `health_check_healthy`, `health_check_unhealthy`, `health_check_toggle`, `site_online`, `site_offline`, `site_toggle`, `resource_healthy`, `resource_unhealthy`, `resource_degraded`, `resource_toggle`. |
| `health_check_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `id` | String | Read-only |  |
| `last_triggered_at` | Int | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `resource_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `site_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `updated_at` | Int | Read-only |  |
| `user_ids` | Set of String | Optional |  Default: `[]`. |
| `webhook_actions` | List of Object | Optional, Sensitive |  Default: `[]`. |

### Nested webhook_actions

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `webhook_actions[].config` | String | Optional, Sensitive |  Default: `""`. |
| `webhook_actions[].enabled` | Bool | Optional |  Default: `true`. |
| `webhook_actions[].webhook_url` | String | Required |  |

## Import

```shell
terraform import pangolin_alert_rule.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
