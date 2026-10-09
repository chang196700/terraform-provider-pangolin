---
page_title: "pangolin_site_provisioning_keys Data Source - pangolin"
description: |-
  Lists site provisioning key objects across all API pages. Secret fields are excluded.
---

# pangolin_site_provisioning_keys (Data Source)

Lists site provisioning key objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_site_provisioning_keys" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].approve_new_sites` | Bool | Optional |  Default: `true`. |
| `items[].created_at` | String | Read-only |  |
| `items[].id` | String | Read-only |  |
| `items[].last_chars` | String | Read-only |  |
| `items[].last_used` | String | Read-only |  |
| `items[].max_batch_size` | Int | Optional | Maximum sites per batch, from 1 through 1000000. Null means no limit. |
| `items[].name` | String | Required, replacement |  |
| `items[].num_used` | Int | Read-only |  |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].valid_until` | String | Optional | RFC3339 expiry timestamp. Empty string removes expiration. Default: `""`. |
