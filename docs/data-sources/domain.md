---
page_title: "pangolin_domain Data Source - pangolin"
description: |-
  Reads a domain available to an organization.
---

# pangolin_domain (Data Source)

Reads a domain available to an organization.

## Example Usage

```terraform
data "pangolin_domain" "example" {
  domain_id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `base_domain` | String | Read-only |  |
| `cert_resolver` | String | Read-only |  |
| `config_managed` | Bool | Read-only |  |
| `custom_cert_resolver` | String | Read-only |  |
| `domain_id` | String | Required |  Retained in state when the API does not return it. |
| `error_message` | String | Read-only |  |
| `failed` | Bool | Read-only |  |
| `last_checked_at` | Int | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |
| `prefer_wildcard_cert` | Bool | Read-only |  |
| `tries` | Int | Read-only |  |
| `type` | String | Read-only |  |
| `verified` | Bool | Read-only |  |
