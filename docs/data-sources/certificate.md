---
page_title: "pangolin_certificate Data Source - pangolin"
description: |-
  OpenAPI compatibility schema for certificate issuance and renewal metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires session authentication. Requires a server version or extension that mounts this route for API keys. Certificate and private-key material are not returned.
---

# pangolin_certificate (Data Source)

OpenAPI compatibility schema for certificate issuance and renewal metadata. Unavailable through the stock Pangolin 1.24.0 Integration API: its route is registered only on the dashboard service, which requires session authentication. Requires a server version or extension that mounts this route for API keys. Certificate and private-key material are not returned.

## Example Usage

```terraform
data "pangolin_certificate" "example" {
  domain = "example"
  domain_id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `cert_id` | String (numeric API ID) | Read-only |  |
| `certificate_domain` | String | Read-only |  |
| `created_at` | Int | Read-only |  |
| `domain` | String | Required |  Retained in state when the API does not return it. |
| `domain_id` | String | Required |  Retained in state when the API does not return it. |
| `domain_type` | String | Read-only |  |
| `error_message` | String | Read-only |  |
| `expires_at` | Int | Read-only |  |
| `last_renewal_attempt` | Int | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |
| `renewal_count` | Int | Read-only |  |
| `status` | String | Read-only |  |
| `updated_at` | Int | Read-only |  |
| `wildcard` | Bool | Read-only |  |
