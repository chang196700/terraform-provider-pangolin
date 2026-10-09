---
page_title: "pangolin_client Resource - pangolin"
description: |-
  A Pangolin OLM client. Credentials and address defaults are allocated only during creation. subnet is a host IP address, without a CIDR prefix.
---

# pangolin_client (Resource)

A Pangolin OLM client. Credentials and address defaults are allocated only during creation. subnet is a host IP address, without a CIDR prefix.

## Example Usage

```terraform
resource "pangolin_client" "example" {
  name   = "example"
  org_id = "example"
  type   = "olm"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `archived` | Bool | Optional | Whether the client is archived. Defaults to false. Default: `false`. |
| `blocked` | Bool | Optional | Whether the client is blocked. Blocking terminates connected OLM sessions. Defaults to false. Default: `false`. |
| `id` | String | Read-only |  |
| `name` | String | Required |  |
| `nice_id` | String | Optional |  |
| `olm_id` | String | Optional, replacement |  Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `secret` | String | Optional, Sensitive, replacement |  Retained in state when the API does not return it. |
| `subnet` | String | Optional, replacement | Host IP address without CIDR prefix. Allocated from the organization only on create when omitted. |
| `type` | String | Optional, replacement |  Default: `"olm"`. Allowed: `olm`. |

## Import

```shell
terraform import pangolin_client.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
