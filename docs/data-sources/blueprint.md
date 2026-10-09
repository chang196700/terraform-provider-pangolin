---
page_title: "pangolin_blueprint Data Source - pangolin"
description: |-
  Reads an existing blueprint execution. This data source never applies the blueprint. Contents may include credentials and are stored as sensitive state.
---

# pangolin_blueprint (Data Source)

Reads an existing blueprint execution. This data source never applies the blueprint. Contents may include credentials and are stored as sensitive state.

## Example Usage

```terraform
data "pangolin_blueprint" "example" {
  blueprint_id = "1"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `blueprint_id` | String (numeric API ID) | Required |  Retained in state when the API does not return it. |
| `contents` | String | Read-only, Sensitive |  |
| `created_at` | Int | Read-only |  |
| `message` | String | Read-only |  |
| `name` | String | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |
| `source` | String | Read-only |  |
| `succeeded` | Bool | Read-only |  |
