---
page_title: "pangolin_identity_provider_organization Data Source - pangolin"
description: |-
  Reads an existing identity provider organization.
---

# pangolin_identity_provider_organization (Data Source)

Reads an existing identity provider organization.

## Example Usage

```terraform
data "pangolin_identity_provider_organization" "example" {
  idp_id = "1"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `id` | String | Read-only |  |
| `idp_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `org_mapping` | String | Read-only |  |
| `role_mapping` | String | Read-only |  |
