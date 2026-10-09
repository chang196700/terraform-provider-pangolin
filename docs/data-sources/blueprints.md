---
page_title: "pangolin_blueprints Data Source - pangolin"
description: |-
  Lists previous blueprint executions without retrieving their potentially sensitive contents.
---

# pangolin_blueprints (Data Source)

Lists previous blueprint executions without retrieving their potentially sensitive contents.

## Example Usage

```terraform
data "pangolin_blueprints" "example" {
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
| `items[].blueprint_id` | String (numeric API ID) | Optional |  |
| `items[].created_at` | Int | Optional |  |
| `items[].name` | String | Optional |  |
| `items[].org_id` | String | Optional |  |
| `items[].source` | String | Optional |  |
| `items[].succeeded` | Bool | Optional |  |
