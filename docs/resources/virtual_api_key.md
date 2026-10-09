---
page_title: "pangolin_virtual_api_key Resource - pangolin"
description: |-
  Manages a manual virtual API key for public inference resources. The returned secret is sensitive and stored in Terraform state. Updating other attributes does not renew its expiry; changing valid_for_seconds does.
---

# pangolin_virtual_api_key (Resource)

Manages a manual virtual API key for public inference resources. The returned secret is sensitive and stored in Terraform state. Updating other attributes does not renew its expiry; changing valid_for_seconds does.

## Example Usage

```terraform
resource "pangolin_virtual_api_key" "example" {
  all_resources = true
  name          = "example"
  org_id        = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `all_resources` | Bool | Optional |  Default: `false`. |
| `created_at` | Int | Read-only |  |
| `created_by_user_id` | String | Read-only |  |
| `description` | String | Optional |  |
| `expires_at` | Int | Read-only |  |
| `id` | String | Read-only |  |
| `kind` | String | Read-only |  |
| `last_chars` | String | Read-only |  |
| `last_used_at` | Int | Read-only |  |
| `name` | String | Required |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `resource_ids` | Set of String (numeric API ID) | Optional |  Default: `[]`. |
| `secret` | String | Read-only, Sensitive |  |
| `user_id` | String | Optional |  |
| `valid_for_seconds` | Int | Optional | Lifetime from creation, or from the most recent change to this value. Removing it clears the expiry. This input is not returned by the API. Retained in state when the API does not return it. |

## Import

```shell
terraform import pangolin_virtual_api_key.example 'example/object-id'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
