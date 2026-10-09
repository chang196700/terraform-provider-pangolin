---
page_title: "pangolin_user_client Resource - pangolin"
description: |-
  A Pangolin OLM client. Credentials and address defaults are allocated only during creation. subnet is a host IP address, without a CIDR prefix. In Pangolin 1.24.0, the user-client CREATE route is commented out in the Integration API, and the shared DELETE handler rejects clients owned by users. Existing user clients can be imported, read, updated, archived, and blocked through the shared client routes, but create, replacement, and destroy require a server that implements these operations. There is no Bearer-key alternative for the missing lifecycle operations. Mock create/destroy tests verify specification adapter behavior only.
---

# pangolin_user_client (Resource)

A Pangolin OLM client. Credentials and address defaults are allocated only during creation. subnet is a host IP address, without a CIDR prefix. In Pangolin 1.24.0, the user-client CREATE route is commented out in the Integration API, and the shared DELETE handler rejects clients owned by users. Existing user clients can be imported, read, updated, archived, and blocked through the shared client routes, but create, replacement, and destroy require a server that implements these operations. There is no Bearer-key alternative for the missing lifecycle operations. Mock create/destroy tests verify specification adapter behavior only.

## Example Usage

```terraform
resource "pangolin_user_client" "example" {
  name    = "example"
  olm_id  = "existing-user-olm-id"
  org_id  = "example"
  type    = "olm"
  user_id = "example-id"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `archived` | Bool | Optional | Whether the client is archived. Defaults to false. Default: `false`. |
| `blocked` | Bool | Optional | Whether the client is blocked. Blocking terminates connected OLM sessions. Defaults to false. Default: `false`. |
| `id` | String | Read-only |  |
| `name` | String | Required | Client display name. Pangolin 1.24.0 replaces this with a recognized Mac/iOS model name for connected user devices; use that returned display name when importing those devices because custom names cannot be read back reliably. |
| `nice_id` | String | Optional |  |
| `olm_id` | String | Required, replacement | Existing OLM ID belonging to user_id. The provider does not generate an OLM for user clients. Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `subnet` | String | Optional, replacement | Host IP address without CIDR prefix. Allocated from the organization only on create when omitted. |
| `type` | String | Optional, replacement |  Default: `"olm"`. Allowed: `olm`. |
| `user_id` | String | Required, replacement |  Retained in state when the API does not return it. |

## Import

```shell
terraform import pangolin_user_client.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
