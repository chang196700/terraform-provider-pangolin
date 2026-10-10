# Terraform Provider for Pangolin

Manage Pangolin organizations, sites, public/private resources, access policies, identities, credentials, monitoring and AI gateway configuration with Terraform.

Provider address: **`chang196700/pangolin`**. API baseline: **Pangolin 1.24.0**, commit `a235635d9e5a6431ea426431e3e862e965928cd3`. Requires Terraform **1.11+**; developing the provider requires Go **1.25+**.

[中文说明](README.zh-CN.md) · [Provider configuration](docs/index.md) · [Resources](docs/resources) · [Data sources](docs/data-sources) · [API coverage](docs/api-coverage.md) · [Compatibility](docs/compatibility.md) · [Permissions](docs/permissions.md) · [Release procedure](docs/releasing.md)

**Pangolin 1.24.0 has documented API limitations.** Some operations declared in its OpenAPI are not mounted on the Integration API, and some mounted handlers require session-only context. All 24 resource schemas are included, but not every lifecycle is usable on an unmodified 1.24.0 server. Read the [compatibility table](docs/compatibility.md) before adopting a resource. Mock tests verify the adapters; they do not prove availability of missing upstream routes.

## Quick start

Install the published provider directly from the Terraform Registry:

```hcl
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
  endpoint = "https://api.example.com/v1"
  # Set PANGOLIN_API_KEY in the environment.
}

resource "pangolin_site" "office" {
  org_id = "your-org-id"
  name   = "Office"
  type   = "newt"
}

output "newt_id" {
  value = pangolin_site.office.newt_id
}

output "newt_secret" {
  value     = pangolin_site.office.secret
  sensitive = true
}
```

The `v0.1.1` provider is published at [registry.terraform.io/chang196700/pangolin](https://registry.terraform.io/providers/chang196700/pangolin). Enable the Integration API on self-hosted Pangolin and grant the API key the [actions required by your resources](docs/permissions.md). Organizations, global IdPs and organization API-key management require a root key. Commercial features require the corresponding Pangolin edition. A root key cannot enable an upstream route that is absent or requires a browser session.

`PANGOLIN_API_KEY` must contain the complete `apiKeyId.apiKeySecret` credential. The `pangolin_organization_api_key.api_key` resource attribute contains only the secret portion; combine it with that resource's `id` when configuring a client.

## Coverage

The provider exposes **25 managed resource types**, including `pangolin_domain`, plus single-object data sources, paginated list data sources, and additional domain, DNS, resource-name, certificate, Blueprint, invitation, device and AI catalog queries. HCL attributes are typed and use snake_case. IDs are strings; numeric API IDs are converted at the request boundary.

The supplied `openapi.yaml` remains unchanged. [The inventory](docs/api-coverage.md) classifies all **312 operations**, including aliases and intentionally excluded workflows. API descriptions omit most response schemas, so adapters and fixtures also use the pinned upstream implementation. Upstream source is a development reference in ignored `.cache/`, not distributed with this provider.

## State ownership

- Role/user/client and AI association collections are authoritative. Empty collections remove the provider-managed assignments; system-mandated administrator membership is preserved.
- Use either `public_resource.access` for inline authentication or `resource_policy_id` for a shared policy. Pangolin 1.24.0 routes can modify the shared policy instead of a local override, so configuring both is rejected. Authentication functionality is available on `public_resource_policy`.
- `organization_user` owns organization membership. Creation supports OIDC users; destruction does not delete the global user. Internal users can be imported. Removing a user's last role or an owner's mandatory role is rejected by Pangolin.
- Deleting a role requires `reassign_role_id`, the surviving role that will receive its users. Deleting `user_server_admin` demotes the existing user.
- Newt credentials and ordinary API-key values are returned only at creation. Import cannot recover them. Passwords that are not readable are retained locally and cannot provide complete external drift detection.
- Sensitive values are hidden in normal Terraform output **but are stored in state**. Store state securely. The provider does not log bearer tokens or request/response bodies.
- A virtual key's `valid_for_seconds` is applied at creation or when that attribute changes; unrelated updates do not extend expiry. No resource sends notification emails.
- Objects use explicit import identities documented per resource. Sites are deleted with `deleteResources=false` so Terraform controls dependency ordering.

## Develop and test

```shell
go test ./...
go vet ./...
go run ./cmd/docs -check
```

Run real Terraform against the local HTTP fixtures:

```shell
export TF_ACC=1
export TF_ACC_TERRAFORM_PATH="$(command -v terraform)"
go test ./internal/provider -run TestAcc -count=1 -timeout 30m
```

PowerShell:

```powershell
. ./scripts/env.ps1
$env:TF_ACC = '1'
$env:TF_ACC_TERRAFORM_PATH = (Get-Command terraform).Source
go test ./internal/provider -run TestAcc -count=1 -timeout 30m
```

The fixtures exercise create, empty plans, update, import, external deletion and destroy. Additional tests cover pagination, permissions versus absence, secret preservation, partial failures, policy transitions, budget scope changes and expiry preservation. **These are protocol/contract tests, not a live Pangolin deployment test.**

For a live, disposable test organization, explicitly set `PANGOLIN_LIVE_TEST=1`, `PANGOLIN_LIVE_ORG_ID`, `PANGOLIN_ENDPOINT`, `PANGOLIN_API_KEY` and `TF_ACC=1`, then run `go test ./internal/provider -run TestAccLiveSite -count=1`. This creates, updates and deletes a local site in that organization. Never enable it implicitly in normal CI.

## Local build and installation

```shell
go build -o terraform-provider-pangolin .
go run ./cmd/release -version 0.1.1
```

The release builder emits six platform ZIP archives, a protocol manifest and SHA-256 checksums under `dist/`. The v0.1.1 release is signed with the Registry key and published on [GitHub Releases](https://github.com/chang196700/terraform-provider-pangolin/releases/tag/v0.1.1). `scripts/verify-mirror.ps1` installs the Windows build into an isolated filesystem mirror, runs `terraform init`, checks provider schema counts, and validates all generated examples.

For development, configure a Terraform CLI `dev_overrides` entry for `registry.terraform.io/chang196700/pangolin` pointing to the folder containing your binary. Alternatively put a release ZIP in `MIRROR/registry.terraform.io/chang196700/pangolin/` and configure `provider_installation.filesystem_mirror` with that directory. A mirror installation does not prove Registry availability.

Regenerate schema-derived documentation and verify specification transcription:

```shell
go run ./cmd/specgen
go run ./cmd/docs
```

## License

MPL-2.0. This is an independently implemented community provider, not an official Fossorial product.
