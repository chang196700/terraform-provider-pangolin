---
page_title: "pangolin_organization_oidc_identity_provider Resource - pangolin"
description: |-
  An organization-scoped OIDC identity provider. Requires the corresponding Pangolin edition and organization identity provider mode.
---

# pangolin_organization_oidc_identity_provider (Resource)

An organization-scoped OIDC identity provider. Requires the corresponding Pangolin edition and organization identity provider mode.

## Example Usage

```terraform
resource "pangolin_organization_oidc_identity_provider" "example" {
  auth_url        = "https://idp.example.com/"
  client_id       = "example"
  client_secret   = "replace-with-a-secure-value"
  identifier_path = "example"
  name            = "example"
  org_id          = "example"
  scopes          = "example"
  token_url       = "https://idp.example.com/"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `auth_url` | String | Required |  |
| `auto_provision` | Bool | Optional |  |
| `client_id` | String | Required |  |
| `client_secret` | String | Required, Sensitive |  |
| `email_path` | String | Optional |  |
| `id` | String | Read-only |  |
| `identifier_path` | String | Required |  |
| `name` | String | Required |  |
| `name_path` | String | Optional |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `org_mapping` | String | Optional |  |
| `redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `role_mapping` | String | Optional |  |
| `scopes` | String | Required |  |
| `tags` | String | Optional |  |
| `token_url` | String | Required |  |
| `variant` | String | Optional, replacement |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |

## Import

```shell
terraform import pangolin_organization_oidc_identity_provider.example 'example/2'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
