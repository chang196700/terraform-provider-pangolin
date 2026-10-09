---
page_title: "pangolin_domains Data Source - pangolin"
description: |-
  Lists all domains available to an organization.
---

# pangolin_domains (Data Source)

Lists all domains available to an organization.

## Example Usage

```terraform
data "pangolin_domains" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].base_domain` | String | Optional |  |
| `items[].cert_resolver` | String | Optional |  |
| `items[].config_managed` | Bool | Optional |  |
| `items[].domain_id` | String | Optional |  |
| `items[].error_message` | String | Optional |  |
| `items[].failed` | Bool | Optional |  |
| `items[].prefer_wildcard_cert` | Bool | Optional |  |
| `items[].tries` | Int | Optional |  |
| `items[].type` | String | Optional |  |
| `items[].verified` | Bool | Optional |  |
