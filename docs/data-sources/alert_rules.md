---
page_title: "pangolin_alert_rules Data Source - pangolin"
description: |-
  Lists alert rule objects across all API pages. Secret fields are excluded.
---

# pangolin_alert_rules (Data Source)

Lists alert rule objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_alert_rules" "example" {
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
| `items[].all_health_checks` | Bool | Optional |  Default: `false`. |
| `items[].all_resources` | Bool | Optional |  Default: `false`. |
| `items[].all_sites` | Bool | Optional |  Default: `false`. |
| `items[].cooldown_seconds` | Int | Optional |  Default: `0`. |
| `items[].created_at` | Int | Read-only |  |
| `items[].emails` | Set of String | Optional |  Default: `[]`. |
| `items[].enabled` | Bool | Optional |  Default: `true`. |
| `items[].event_type` | String | Required |  Allowed: `health_check_healthy`, `health_check_unhealthy`, `health_check_toggle`, `site_online`, `site_offline`, `site_toggle`, `resource_healthy`, `resource_unhealthy`, `resource_degraded`, `resource_toggle`. |
| `items[].health_check_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].id` | String | Read-only |  |
| `items[].last_triggered_at` | Int | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `items[].resource_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].role_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].site_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `items[].updated_at` | Int | Read-only |  |
| `items[].user_ids` | Set of String | Optional |  Default: `[]`. |
