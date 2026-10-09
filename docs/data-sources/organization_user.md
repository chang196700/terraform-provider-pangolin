---
page_title: "pangolin_organization_user Data Source - pangolin"
description: |-
  Reads an existing organization user.
---

# pangolin_organization_user (Data Source)

Reads an existing organization user.

## Example Usage

```terraform
data "pangolin_organization_user" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `auto_provisioned` | Bool | Read-only |  |
| `email` | String | Read-only, replacement |  |
| `id` | String | Required |  |
| `idp_id` | String (numeric API ID) | Required, replacement |  |
| `name` | String | Read-only, replacement |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `role_ids` | Set of String (numeric API ID) | Required | Nonempty authoritative role ID set. New assignments are added before obsolete ones are removed. |
| `type` | String | Read-only, replacement |  Default: `"oidc"`. Allowed: `oidc`. |
| `username` | String | Required, replacement |  |
