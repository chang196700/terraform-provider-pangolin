---
page_title: "pangolin_ai_catalog_models Data Source - pangolin"
description: |-
  Lists known catalog model keys. Supply provider_id, or supply both org_id and type to query before creating a provider.
---

# pangolin_ai_catalog_models (Data Source)

Lists known catalog model keys. Supply provider_id, or supply both org_id and type to query before creating a provider.

## Example Usage

```terraform
data "pangolin_ai_catalog_models" "example" {
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items` | List of Object | Read-only |  |
| `org_id` | String | Optional |  Retained in state when the API does not return it. |
| `provider_id` | String (numeric API ID) | Optional |  Retained in state when the API does not return it. |
| `query` | String | Optional |  Retained in state when the API does not return it. |
| `type` | String | Optional |  Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].model` | String | Optional |  |
