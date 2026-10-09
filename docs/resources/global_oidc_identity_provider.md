---
page_title: "pangolin_global_oidc_identity_provider Resource - pangolin"
description: |-
  A global OIDC identity provider. Requires Pangolin global identity provider mode.
---

# pangolin_global_oidc_identity_provider (Resource)

A global OIDC identity provider. Requires Pangolin global identity provider mode.

## Example Usage

```terraform
resource "pangolin_global_oidc_identity_provider" "example" {
  auth_url        = "https://idp.example.com/"
  client_id       = "example"
  client_secret   = "replace-with-a-secure-value"
  identifier_path = "example"
  name            = "example"
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
| `default_org_mapping` | String | Optional |  |
| `default_role_mapping` | String | Optional |  |
| `email_path` | String | Optional |  |
| `id` | String | Read-only |  |
| `identifier_path` | String | Required |  |
| `name` | String | Required |  |
| `name_path` | String | Optional |  |
| `redirect_url` | String | Read-only |  Retained in state when the API does not return it. |
| `scopes` | String | Required |  |
| `tags` | String | Optional |  |
| `token_url` | String | Required |  |
| `variant` | String | Optional |  Default: `"oidc"`. Allowed: `oidc`, `google`, `azure`. |

## Import

```shell
terraform import pangolin_global_oidc_identity_provider.example '2'
```

Import identity: `id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
