---
page_title: "pangolin_organization_user Resource - pangolin"
description: |-
  An OIDC user's membership and authoritative role assignments in an organization. Destroy removes membership, not the global user. The API does not support creating internal users. In Pangolin 1.24.0, role reassignment endpoints require the enterprise/SaaS integration router. Reads use the organization list and username lookup because the individual membership GET incorrectly requires a browser session.
---

# pangolin_organization_user (Resource)

An OIDC user's membership and authoritative role assignments in an organization. Destroy removes membership, not the global user. The API does not support creating internal users. In Pangolin 1.24.0, role reassignment endpoints require the enterprise/SaaS integration router. Reads use the organization list and username lookup because the individual membership GET incorrectly requires a browser session.

## Example Usage

```terraform
resource "pangolin_organization_user" "example" {
  idp_id   = "1"
  org_id   = "example"
  role_ids = ["1"]
  type     = "oidc"
  username = "user@example.com"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `auto_provisioned` | Bool | Optional |  |
| `email` | String | Optional, replacement |  |
| `id` | String | Read-only |  |
| `idp_id` | String (numeric API ID) | Required, replacement |  |
| `name` | String | Optional, replacement |  |
| `org_id` | String | Required, replacement | Organization identifier. Retained in state when the API does not return it. |
| `role_ids` | Set of String (numeric API ID) | Required | Nonempty authoritative role ID set. New assignments are added before obsolete ones are removed. |
| `type` | String | Optional, replacement |  Default: `"oidc"`. Allowed: `oidc`. |
| `username` | String | Required, replacement |  |

## Import

```shell
terraform import pangolin_organization_user.example 'example/object-id'
```

Import identity: `org_id/id`. Import recovers readable metadata; one-time secrets and write-only input values cannot be recovered.

## API compatibility and permissions

Contract checked against Pangolin 1.24.0 (`a235635d9e5a6431ea426431e3e862e965928cd3`). See [known route limitations](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/compatibility.md) and [required API-key actions](https://github.com/chang196700/terraform-provider-pangolin/blob/main/docs/permissions.md). Automated acceptance tests use an HTTP fixture server, not a live Pangolin deployment. Edition-specific operations require the corresponding server capabilities.
