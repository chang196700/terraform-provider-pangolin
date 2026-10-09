---
page_title: "pangolin_site_provisioning_key Resource - pangolin"
description: |-
  An organization site provisioning key. The secret is returned only during creation, retained in sensitive state, and cannot be recovered by import.
---

# pangolin_site_provisioning_key (Resource)

An organization site provisioning key. The secret is returned only during creation, retained in sensitive state, and cannot be recovered by import.

## Example Usage

```terraform
resource "pangolin_site_provisioning_key" "example" {
  name   = "example"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `approve_new_sites` | Bool | Optional |  Default: `true`. |
| `created_at` | String | Read-only |  |
| `id` | String | Read-only |  |
| `last_chars` | String | Read-only |  |
| `last_used` | String | Read-only |  |
| `max_batch_size` | Int | Optional | Maximum sites per batch, from 1 through 1000000. Null means no limit. |
| `name` | String | Required, replacement |  |
| `num_used` | Int | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `site_provisioning_key` | String | Read-only, Sensitive |  Retained in state when the API does not return it. |
| `valid_until` | String | Optional | RFC3339 expiry timestamp. Empty string removes expiration. Default: `""`. |

## Import

```shell
terraform import pangolin_site_provisioning_key.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
