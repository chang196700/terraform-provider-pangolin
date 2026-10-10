---
page_title: "Pangolin Provider"
description: |-
  Manage Pangolin networks, access configuration, identities, monitoring and AI gateway infrastructure.
---

# Pangolin Provider

The Pangolin provider manages the Integration API. It targets Pangolin 1.24.0 and Terraform 1.11+. Configure a bearer API key with permission to perform the operations required by your resources.

## Example Usage

```terraform
terraform {
  required_version = ">= 1.11.0"
  required_providers {
    pangolin = {
      source  = "chang196700/pangolin"
      version = "~> 0.1.1"
    }
  }
}

provider "pangolin" {
  endpoint        = "https://api.example.com/v1"
  request_timeout = 30
  # Read the API key from PANGOLIN_API_KEY.
}
```

## Schema

| Attribute | Type | Description |
|---|---|---|
| `endpoint` | String | Absolute Integration API URL including `/v1`. Falls back to `PANGOLIN_ENDPOINT`. Required through configuration or environment. |
| `api_key` | String, Sensitive | Complete bearer credential in `apiKeyId.apiKeySecret` form. Falls back to `PANGOLIN_API_KEY`. Required through configuration or environment. |
| `request_timeout` | Number | Positive HTTP timeout in seconds; defaults to 30. |

Organization IDs are explicit resource inputs. Use provider aliases for different servers or API credentials. The client verifies TLS using the operating-system trust store.

## Permissions and compatibility

Organization API keys can manage their permitted organization operations. Organization API-key management, organization creation and global identity-provider operations require root keys. See the [literal action IDs](permissions.md) and [known 1.24.0 route limitations](compatibility.md), including unavailable administrator mutations and defective organization deletion. Commercial operations require an edition that exposes those routes. A permission failure is reported without deleting Terraform state.

Sensitive values are still stored in state. One-time generated credentials cannot be recovered on import. Resources that retain unreadable input values document their drift-detection limitations.

Use either inline public-resource access or a shared public-resource policy. Pangolin 1.24.0 does not provide safe independent overrides for all authentication setters, so the provider rejects combining them.
