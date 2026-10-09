---
page_title: "pangolin_domain_dns_records Data Source - pangolin"
description: |-
  Lists DNS records required for a Pangolin domain, including verification status and the server's effective record values.
---

# pangolin_domain_dns_records (Data Source)

Lists DNS records required for a Pangolin domain, including verification status and the server's effective record values.

## Example Usage

```terraform
data "pangolin_domain_dns_records" "example" {
  domain_id = "example-id"
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `domain_id` | String | Required |  Retained in state when the API does not return it. |
| `items` | List of Object | Read-only |  |
| `org_id` | String | Required |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].base_domain` | String | Optional |  |
| `items[].domain_id` | String | Optional |  |
| `items[].id` | String (numeric API ID) | Optional |  |
| `items[].record_type` | String | Optional |  |
| `items[].value` | String | Optional |  |
| `items[].verified` | Bool | Optional |  |
