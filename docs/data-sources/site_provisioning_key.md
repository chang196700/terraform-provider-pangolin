---
page_title: "pangolin_site_provisioning_key Data Source - pangolin"
description: |-
  Reads an existing site provisioning key.
---

# pangolin_site_provisioning_key (Data Source)

Reads an existing site provisioning key.

## Example Usage

```terraform
data "pangolin_site_provisioning_key" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `approve_new_sites` | Bool | Read-only |  Default: `true`. |
| `created_at` | String | Read-only |  |
| `id` | String | Required |  |
| `last_chars` | String | Read-only |  |
| `last_used` | String | Read-only |  |
| `max_batch_size` | Int | Read-only | Maximum sites per batch, from 1 through 1000000. Null means no limit. |
| `name` | String | Required, replacement |  |
| `num_used` | Int | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `site_provisioning_key` | String | Read-only, Sensitive |  Retained in state when the API does not return it. |
| `valid_until` | String | Read-only | RFC3339 expiry timestamp. Empty string removes expiration. Default: `""`. |
