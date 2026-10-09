---
page_title: "pangolin_clients Data Source - pangolin"
description: |-
  Lists client objects across all API pages. Secret fields are excluded.
---

# pangolin_clients (Data Source)

Lists client objects across all API pages. Secret fields are excluded.

## Example Usage

```terraform
data "pangolin_clients" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].archived` | Bool | Optional | Whether the client is archived. Defaults to false. Default: `false`. |
| `items[].blocked` | Bool | Optional | Whether the client is blocked. Blocking terminates connected OLM sessions. Defaults to false. Default: `false`. |
| `items[].id` | String | Read-only |  |
| `items[].name` | String | Required |  |
| `items[].nice_id` | String | Optional |  |
| `items[].olm_id` | String | Optional, replacement |  Retained in state when the API does not return it. |
| `items[].org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `items[].subnet` | String | Optional, replacement | Host IP address without CIDR prefix. Allocated from the organization only on create when omitted. |
| `items[].type` | String | Optional, replacement |  Default: `"olm"`. Allowed: `olm`. |
