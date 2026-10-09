---
page_title: "pangolin_user_client Data Source - pangolin"
description: |-
  Reads an existing user client.
---

# pangolin_user_client (Data Source)

Reads an existing user client.

## Example Usage

```terraform
data "pangolin_user_client" "example" {
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
| `name` | String | Required | Client display name. Pangolin 1.24.0 replaces this with a recognized Mac/iOS model name for connected user devices; use that returned display name when importing those devices because custom names cannot be read back reliably. |
| `nice_id` | String | Read-only |  |
| `olm_id` | String | Required, replacement | Existing OLM ID belonging to user_id. The provider does not generate an OLM for user clients. Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `subnet` | String | Read-only, replacement | Host IP address without CIDR prefix. Allocated from the organization only on create when omitted. |
| `type` | String | Read-only, replacement |  Default: `"olm"`. Allowed: `olm`. |
| `user_id` | String | Required, replacement |  Retained in state when the API does not return it. |
