---
page_title: "pangolin_identity_provider_organizations Data Source - pangolin"
description: |-
  Lists identity provider organization objects across all API pages. Secret fields are excluded.
---

# pangolin_identity_provider_organizations (Data Source)

Lists identity provider organization objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_identity_provider_organizations" "example" {
  idp_id = "1"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `idp_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `items` | List of Object | Read-only |  |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].id` | String | Read-only |  |
| `items[].idp_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].org_mapping` | String | Optional |  |
| `items[].role_mapping` | String | Optional |  |
