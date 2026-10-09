---
page_title: "pangolin_alert_rule Data Source - pangolin"
description: |-
  Reads an existing alert rule.
---

# pangolin_alert_rule (Data Source)

Reads an existing alert rule.

## Example Usage

```terraform
data "pangolin_alert_rule" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `all_health_checks` | Bool | Read-only |  Default: `false`. |
| `all_resources` | Bool | Read-only |  Default: `false`. |
| `all_sites` | Bool | Read-only |  Default: `false`. |
| `cooldown_seconds` | Int | Read-only |  Default: `0`. |
| `created_at` | Int | Read-only |  |
| `emails` | Set of String | Read-only |  Default: `[]`. |
| `enabled` | Bool | Read-only |  Default: `true`. |
| `event_type` | String | Required |  Allowed: `health_check_healthy`, `health_check_unhealthy`, `health_check_toggle`, `site_online`, `site_offline`, `site_toggle`, `resource_healthy`, `resource_unhealthy`, `resource_degraded`, `resource_toggle`. |
| `health_check_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `id` | String | Required |  |
| `last_triggered_at` | Int | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `resource_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `role_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `site_ids` | Set of String (numeric API ID) | Read-only |  Default: `[]`. |
| `updated_at` | Int | Read-only |  |
| `user_ids` | Set of String | Read-only |  Default: `[]`. |
| `webhook_actions` | List of Object | Read-only, Sensitive |  Default: `[]`. |

### Nested webhook_actions

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `webhook_actions[].config` | String | Optional, Sensitive |  Default: `""`. |
| `webhook_actions[].enabled` | Bool | Optional |  Default: `true`. |
| `webhook_actions[].webhook_url` | String | Required |  |
