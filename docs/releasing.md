# Release procedure

Repository: `chang196700/terraform-provider-pangolin`. Provider: `registry.terraform.io/chang196700/pangolin`.

1. Run unit tests, all fixture acceptance tests, static checks, generated-document checks and the local mirror installation check. Record the Terraform and Go versions. Live Pangolin tests are separate and opt-in.
2. Build `go run ./cmd/release -version 0.1.1`. Verify all six archives and the manifest against `dist/terraform-provider-pangolin_0.1.1_SHA256SUMS`.
3. Sign the checksum file with the owner's Registry signing key. Use a Registry-supported key type and upload its public key in the `chang196700` namespace. Keep the private key outside the repository.
4. Publish the reviewed source to the public `terraform-provider-pangolin` repository and tag `v0.1.1`. Release assets must include six ZIPs, the manifest, SHA256SUMS and its detached binary `.sig` signature.
5. Sign into Terraform Registry with GitHub, select **Publish → Provider**, then select this repository. Complete any account or organization authorization in the browser.
6. From a clean directory with no filesystem mirror or development overrides, run `terraform init` using `source = "chang196700/pangolin"` and `version = "0.1.1"`. Confirm the signing-key fingerprint and install success. This is the final Registry verification.

The Actions release workflow uses `GPG_PRIVATE_KEY` and `PASSPHRASE` repository secrets and runs only when the repository variable `RELEASE_AUTOMATION_ENABLED` is `true`. Enable it after supplying a key explicitly intended for automated release signing. Do not export a personal key or place it in CI merely to enable this workflow. For local signing, a key can remain in GnuPG or hardware-backed storage; upload only public release assets.

Manual signing and verification:

```shell
gpg --local-user SIGNING_FINGERPRINT --output dist/terraform-provider-pangolin_0.1.1_SHA256SUMS.sig --detach-sign dist/terraform-provider-pangolin_0.1.1_SHA256SUMS
gpg --verify dist/terraform-provider-pangolin_0.1.1_SHA256SUMS.sig dist/terraform-provider-pangolin_0.1.1_SHA256SUMS
```

The GoReleaser workflow provides the same archive layout and protocol 6 manifest. Never replace assets for an already published version; issue a new version for fixes.
