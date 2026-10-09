---
page_title: "pangolin_user_devices Data Source - pangolin"
description: |-
  Lists organization user devices across all pages, including identity, network and device metadata.
---

# pangolin_user_devices (Data Source)

Lists organization user devices across all pages, including identity, network and device metadata.

## Example Usage

```terraform
data "pangolin_user_devices" "example" {
  org_id = "example"
}
```

## Schema

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `agent` | String | Optional |  Retained in state when the API does not return it. |
| `items` | List of Object | Read-only |  |
| `online` | Bool | Optional |  Retained in state when the API does not return it. |
| `org_id` | String | Required |  Retained in state when the API does not return it. |
| `query` | String | Optional |  Retained in state when the API does not return it. |
| `status` | Set of String | Optional | Optional device status filter: active, pending, denied, blocked or archived. Pangolin defaults to active and pending when omitted. Retained in state when the API does not return it. |

### Nested items

| Attribute | Type | Usage | Description |
|---|---|---|---|
| `items[].agent` | String | Optional |  |
| `items[].approval_state` | String | Optional |  |
| `items[].archived` | Bool | Optional |  |
| `items[].blocked` | Bool | Optional |  |
| `items[].client_id` | String (numeric API ID) | Optional |  |
| `items[].device_model` | String | Optional |  |
| `items[].fingerprint_arch` | String | Optional |  |
| `items[].fingerprint_hostname` | String | Optional |  |
| `items[].fingerprint_kernel_version` | String | Optional |  |
| `items[].fingerprint_os_version` | String | Optional |  |
| `items[].fingerprint_platform` | String | Optional |  |
| `items[].fingerprint_serial_number` | String | Optional |  |
| `items[].fingerprint_username` | String | Optional |  |
| `items[].first_seen` | Int | Optional |  |
| `items[].idp_name` | String | Optional |  |
| `items[].idp_variant` | String | Optional |  |
| `items[].last_seen` | Int | Optional |  |
| `items[].megabytes_in` | Number | Optional |  |
| `items[].megabytes_out` | Number | Optional |  |
| `items[].name` | String | Optional |  |
| `items[].nice_id` | String | Optional |  |
| `items[].olm_archived` | Bool | Optional |  |
| `items[].olm_update_available` | Bool | Optional |  |
| `items[].olm_version` | String | Optional |  |
| `items[].online` | Bool | Optional |  |
| `items[].org_id` | String | Optional |  |
| `items[].org_name` | String | Optional |  |
| `items[].pub_key` | String | Optional |  |
| `items[].subnet` | String | Optional |  |
| `items[].type` | String | Optional |  |
| `items[].user_email` | String | Optional |  |
| `items[].user_id` | String | Optional |  |
| `items[].user_type` | String | Optional |  |
| `items[].username` | String | Optional |  |
