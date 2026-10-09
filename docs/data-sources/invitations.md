---
page_title: "pangolin_invitations Data Source - pangolin"
description: |-
  Lists pending organization invitations and their assigned roles. This data source never sends invitations or email.
---

# pangolin_invitations (Data Source)

Lists pending organization invitations and their assigned roles. This data source never sends invitations or email.

## Example Usage

```terraform
data "pangolin_invitations" "example" {
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
| `items[].email` | String | Optional |  |
| `items[].expires_at` | Int | Optional |  |
| `items[].invite_id` | String | Optional |  |
| `items[].roles` | List of Object | Optional |  |

### Nested items[].roles

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].roles[].role_id` | String (numeric API ID) | Optional |  |
| `items[].roles[].role_name` | String | Optional |  |
