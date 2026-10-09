---
page_title: "pangolin_organization_api_key Data Source - pangolin"
description: |-
  Reads an existing organization api key.
---

# pangolin_organization_api_key (Data Source)

Reads an existing organization api key.

## Example Usage

```terraform
data "pangolin_organization_api_key" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `action_ids` | Set of String | Required | Nonempty authoritative set of API action identifiers. The server cannot replace this set with an empty set. |
| `api_key` | String | Read-only, Sensitive | One-time API key secret. The server returns only the secret component; form the Bearer credential as id.api_key (for example, ${pangolin_organization_api_key.example.id}.${pangolin_organization_api_key.example.api_key}). Import cannot recover this secret. Retained in state when the API does not return it. |
| `created_at` | String | Read-only |  |
| `id` | String | Required |  |
| `last_chars` | String | Read-only |  |
| `name` | String | Required, replacement |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
