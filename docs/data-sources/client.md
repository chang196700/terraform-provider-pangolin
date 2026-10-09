---
page_title: "pangolin_client Data Source - pangolin"
description: |-
  Reads an existing client.
---

# pangolin_client (Data Source)

Reads an existing client.

## Example Usage

```terraform
data "pangolin_client" "example" {
  id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `archived` | Bool | Read-only | Whether the client is archived. Defaults to false. Default: `false`. |
| `blocked` | Bool | Read-only | Whether the client is blocked. Blocking terminates connected OLM sessions. Defaults to false. Default: `false`. |
| `id` | String | Required |  |
| `name` | String | Required |  |
| `nice_id` | String | Read-only |  |
| `olm_id` | String | Read-only, replacement |  Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `secret` | String | Read-only, Sensitive, replacement |  Retained in state when the API does not return it. |
| `subnet` | String | Read-only, replacement | Host IP address without CIDR prefix. Allocated from the organization only on create when omitted. |
| `type` | String | Read-only, replacement |  Default: `"olm"`. Allowed: `olm`. |
