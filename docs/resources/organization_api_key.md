---
page_title: "pangolin_organization_api_key Resource - pangolin"
description: |-
  An organization API key with authoritative action assignments. The API key value is returned only during creation and is stored as sensitive state. Imported keys cannot recover their original secret.
---

# pangolin_organization_api_key (Resource)

An organization API key with authoritative action assignments. The API key value is returned only during creation and is stored as sensitive state. Imported keys cannot recover their original secret.

## Example Usage

```terraform
resource "pangolin_organization_api_key" "example" {
  action_ids = ["createSite", "getSite", "listSites", "updateSite", "deleteSite"]
  name       = "example"
  org_id     = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `action_ids` | Set of String | Required | Nonempty authoritative set of API action identifiers. The server cannot replace this set with an empty set. |
| `api_key` | String | Read-only, Sensitive | One-time API key secret. The server returns only the secret component; form the Bearer credential as id.api_key (for example, ${pangolin_organization_api_key.example.id}.${pangolin_organization_api_key.example.api_key}). Import cannot recover this secret. Retained in state when the API does not return it. |
| `created_at` | String | Read-only |  |
| `id` | String | Read-only |  |
| `last_chars` | String | Read-only |  |
| `name` | String | Required, replacement |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |

## Import

```shell
terraform import pangolin_organization_api_key.example 'example/object-id'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
