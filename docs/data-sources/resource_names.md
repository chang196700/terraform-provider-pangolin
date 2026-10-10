---
page_title: "pangolin_resource_names Data Source - pangolin"
description: |-
  Lists public resource IDs and names for an organization.
---

# pangolin_resource_names (Data Source)

Lists public resource IDs and names for an organization.

## Example Usage

```terraform
data "pangolin_resource_names" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].name` | String | Optional |  |
| `items[].resource_id` | String (numeric API ID) | Optional |  |
