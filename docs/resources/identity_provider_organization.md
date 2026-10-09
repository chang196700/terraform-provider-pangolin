---
page_title: "pangolin_identity_provider_organization Resource - pangolin"
description: |-
  A global identity provider's organization association and mappings. Do not manage the same association through an organization-scoped identity provider resource.
---

# pangolin_identity_provider_organization (Resource)

A global identity provider's organization association and mappings. Do not manage the same association through an organization-scoped identity provider resource.

## Example Usage

```terraform
resource "pangolin_identity_provider_organization" "example" {
  idp_id = "1"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `id` | String | Read-only |  |
| `idp_id` | String (numeric API ID) | Required, replacement |  Retained in state when the API does not return it. |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `org_mapping` | String | Optional |  |
| `role_mapping` | String | Optional |  |

## Import

```shell
terraform import pangolin_identity_provider_organization.example '1/example'
```

Import identity: `idp_id/org_id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
