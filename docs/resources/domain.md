---
page_title: "pangolin_domain Resource - pangolin"
description: |-
  Manages an organization domain. Domain type and base domain are immutable; wildcard domain settings can be updated. Import as org_id/domain_id.
---

# pangolin_domain (Resource)

Manages an organization domain. Domain type and base domain are immutable; wildcard domain settings can be updated. Import as org_id/domain_id.

## Example Usage

```terraform
resource "pangolin_domain" "example" {
  base_domain = "example"
  org_id      = "example"
  type        = "ns"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `base_domain` | String | Required, replacement | The base domain to register. Changes replace the domain. |
| `cert_resolver` | String | Optional | Optional certificate resolver. Only wildcard domains accept updates to this setting. |
| `config_managed` | Bool | Read-only |  |
| `custom_cert_resolver` | String | Read-only |  |
| `error_message` | String | Read-only |  |
| `failed` | Bool | Read-only |  |
| `id` | String | Read-only |  |
| `last_checked_at` | Int | Read-only |  |
| `org_id` | String | Required, replacement |  Retained in state when the API does not return it. |
| `prefer_wildcard_cert` | Bool | Optional | Whether wildcard certificate issuance is preferred for wildcard domains. Default: `false`. |
| `tries` | Int | Read-only |  |
| `type` | String | Required, replacement | Domain integration type. OSS supports wildcard; SaaS supports ns and cname; Enterprise supports all three. Changes replace the domain. Allowed: `ns`, `cname`, `wildcard`. |
| `verified` | Bool | Read-only |  |

## Import

```shell
terraform import pangolin_domain.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
