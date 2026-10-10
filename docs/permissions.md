# API Key permissions

These are the **literal action IDs** checked by Pangolin 1.24.0, commit [`a235635d9e5a6431ea426431e3e862e965928cd3`](https://github.com/fosrl/pangolin/tree/a235635d9e5a6431ea426431e3e862e965928cd3). The upstream enum is named [`ActionsEnum`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/auth/actions.ts), not `APIKeyActions`. [Startup synchronization](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/setup/ensureActions.ts) populates the action table from its string values. Action names are case-sensitive.

Use an API key scoped to the target organization, or a root key where required. Root status does **not** bypass action checks: [`verifyApiKeyHasAction`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/middlewares/integration/verifyApiKeyHasAction.ts) always checks the key's action assignments. Root-only routes separately invoke [`verifyApiKeyIsRoot`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/middlewares/integration/verifyApiKeyIsRoot.ts). Organization and object access middleware additionally restricts which objects the key can reach.

Read permissions are required after create and update, during planning, and for import. A resource that refreshes through a list needs the list action even if no list data source is configured. For the unavailable declarations, see [compatibility limits](compatibility.md); do not invent an action to work around an unregistered route.

## Resource operation matrix

Names below omit the `pangolin_` prefix. `—` means no separate operation; replacement uses the create and delete permissions. `Unavailable` refers to the stock 1.24.0 Integration API. The table covers all 25 declared resource types; this count includes retained compatibility schemas and does not assert that all 25 can complete a lifecycle on the baseline.

| Resource | Create | Refresh / single read | List data source | Update | Delete |
|---|---|---|---|---|---|
| `organization` | `createOrg` **root** | `getOrg` | `listOrgs` **root** | `updateOrg` | `deleteOrg` **root**, but baseline handler fails for API keys; see compatibility |
| `site` | `createSite` | `getSite` | `listSites` | `updateSite` | `deleteSite` |
| `public_resource` | `createResource` | `getResource` | `listResources` | `updateResource` | `deleteResource` |
| `private_resource` | `createSiteResource` | `getSiteResource` + `listSiteResources` + `listResourceRoles` + `listResourceUsers` | `listSiteResources` plus association reads | `updateSiteResource` | `deleteSiteResource` |
| `resource_target` | `createTarget` | `getTarget` | `listTargets` | `updateTarget` | `deleteTarget` |
| `public_resource_policy` | `createResourcePolicy` | `getResourcePolicy` | `listResourcePolicies` | `updateResourcePolicy` | `deleteResourcePolicy` |
| `domain` | `createOrgDomain` | `getDomain` | `listOrgDomains` | `updateOrgDomain` (wildcard settings) | `deleteOrgDomain` |
| `role` | `createRole` | `getRole` + `listRoles` | `listRoles` | `updateRole` | `deleteRole` + `getRole` |
| `organization_user` | `createOrgUser` + `getOrgUser` + `listUsers` | `listUsers` + `getOrgUser` | `listUsers` | `updateOrgUser` + role actions below | `removeUser` |
| `client` | `createClient` | `getClient` | `listClients` | `updateClient` + state actions below | `deleteClient` |
| `user_client` | **Unavailable** | `getClient` | — | `updateClient` + state actions below | **Unavailable**: shared handler rejects user-associated clients |
| `user_server_admin` | **Unavailable** | `getUser` **root** | — | **Unavailable** | **Unavailable** |
| `global_oidc_identity_provider` | `createIdp` **root** | `getIdp` **root** | `listIdps` **root** | `updateIdp` **root** | `deleteIdp` **root** |
| `organization_oidc_identity_provider` | `createIdp` | `getIdp` | `listIdps` | `updateIdp` | `deleteIdp` |
| `identity_provider_organization` | `createIdpOrg` **root** | `listIdpOrgs` **root** | `listIdpOrgs` **root** | `updateIdpOrg` **root** | `deleteIdpOrg` **root** |
| `organization_api_key` | `createApiKey` + `setApiKeyActions` **root** | `listApiKeys` + `listApiKeyActions` **root** | `listApiKeys` + `listApiKeyActions` **root** | `setApiKeyActions` **root** | `deleteApiKey` **root** |
| `site_provisioning_key` | `createSiteProvisioningKey` | `listSiteProvisioningKeys` | `listSiteProvisioningKeys` | `updateSiteProvisioningKey` | `deleteSiteProvisioningKey` |
| `public_resource_access_token` | `generateAccessToken` | `listAccessTokens` | `listAccessTokens` | — | `deleteAcessToken` |
| `ai_provider` | `createAiProvider` | `getAiProvider` | `listAiProviders` | `updateAiProvider` | `deleteAiProvider` |
| `ai_model` | `createAiModel` | `getAiModel` | `listAiModels` | `updateAiModel` | `deleteAiModel` |
| `ai_budget` | `createAiBudget` | `getAiBudget` | `listAiBudgets` | `updateAiBudget` | `deleteAiBudget` |
| `virtual_api_key` | `createVirtualApiKey` | `getVirtualApiKey` | `listVirtualApiKeys` | `updateVirtualApiKey` | `deleteVirtualApiKey` |
| `event_streaming_destination` | `createEventStreamingDestination` | `listEventStreamingDestinations` | `listEventStreamingDestinations` | `updateEventStreamingDestination` | `deleteEventStreamingDestination` |
| `alert_rule` | `createAlertRule` | `getAlertRule` | `listAlertRules` | `updateAlertRule` | `deleteAlertRule` |
| `health_check` | `createHealthCheck` | `listHealthChecks` | `listHealthChecks` | `updateHealthCheck` | `deleteHealthCheck` |

The spelling **`deleteAcessToken`** is intentional: the upstream enum and route both omit the second `c`. There is no `getHealthCheck`, `getEventStreamingDestination`, or `getSiteProvisioningKey` action for these refreshes.

The base actions come from the [base integration router](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts). Shared-policy create/list/delete, organization OIDC, provisioning keys, membership role removal, and monitoring actions are in its [private extension](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/routers/integration.ts).

## Additional calls made by resource fields

| Configuration or operation | Additional actions |
|---|---|
| `public_resource.access` read | `getResourcePolicy`, through `GET /public-resource/{id}/policies`. |
| `public_resource.access` apply or removal | `updateResource`, `setResourceRoles`, `setResourceUsers`, `setResourcePassword`, `setResourcePincode`, `setResourceHeaderAuth`, and `listResourceRules`. Whitelist changes also require `setResourceWhitelist`. Rules require `createResourceRule`, `updateResourceRule`, and `deleteResourceRule` as rows are added, changed, or removed. Setters clear fields even when the configured collection is empty. |
| `public_resource_policy` access control | **Both** `setResourcePolicyUsers` **and** `setResourcePolicyRoles` for the same `/access-control` request; the route contains two consecutive checks. |
| `public_resource_policy` credential/whitelist/rule fields | `setResourcePolicyPassword`, `setResourcePolicyPincode`, `setResourcePolicyHeaderAuth`, `setResourcePolicyWhitelist`, and `setResourcePolicyRules` for their corresponding setters. |
| Public/private inference resource associations | `listResourceAiModels` for both provider and model association reads; `setResourceAiModels` for both provider and model setters. There are no separate `listResourceAiProviders` or `setResourceAiProviders` action IDs. |
| Private resource association refresh | Roles use `listResourceRoles`; users **and clients** use `listResourceUsers`. The direct private-resource update applies its configured associations under `updateSiteResource`. |
| Client `archived` changes | `archiveClient` to archive, `unarchiveClient` to restore. |
| Client `blocked` changes | `blockClient` to block, `unblockClient` to restore. |
| Client automatically allocated connection defaults | `createClient` also protects `GET /org/{orgId}/pick-client-defaults`; no additional allocation action is required. |
| `organization_user.role_ids` reconciliation | `addUserRole` and `removeUserRole`; additions precede removals. `setUserOrgRoles` is not the action checked by these two routes. |
| User lookup after membership creation and during refresh | `getOrgUser` for `GET /org/{orgId}/user-by-username`. Normal refresh first uses `listUsers` to find the immutable user ID and username/IdP pair, then the username route returns role assignments and `autoProvisioned`. This avoids the baseline's session-dependent single-user handler. Grant both actions. |
| Public resource access-token refresh | Uses `GET /org/{orgId}/access-tokens` with `listAccessTokens` and filters by the configured resource. The baseline's resource-scoped token list depends on organization context that is missing for a root key. |
| `role` destroy reassignment validation | `getRole` for the managed role and reassignment role, before `deleteRole`. Every normal refresh also uses `listRoles` to recover the omitted `allowSsh` field. |
| Post-create `nice_id`, or update-only configuration | The owning resource's update action is required during create when these fields trigger a post-create setter. Global OIDC default mappings likewise need `updateIdp`. |

Source anchors: [policy access-control checks](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts#L856), [resource authentication setters](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts#L1108), [client state actions](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts#L1452), and [membership role mutation](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/routers/integration.ts#L229).

## Other data sources

| Data source | API Key action |
|---|---|
| `domain` | `getDomain` |
| `domains` | `listOrgDomains` |
| `domain_dns_records` | `getDNSRecords` |
| `resource_names` | `listResources` |
| `blueprint` | `getBlueprint` |
| `blueprints` | `listBlueprints` |
| `invitations` | `listInvitations` |
| `user_devices` | `listClients` |
| `ai_catalog_models` | `listAiModels`, with either supported lookup route |
| `ai_provider_targets` | `listTargets` |
| `user` | `getUser`, **root key required** |
| `users`, `certificate` | **Unavailable on the baseline**; dashboard permissions are not API Key permissions. |

Derived single/list data sources use the resource reads above, including enrichment calls where present.

## Edition and configuration requirements

The [build resolver](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/esbuild.mjs#L185) selects the public router for OSS and the private extension for Enterprise/SaaS. Thus the private routes below are absent from an OSS build even though a combined OpenAPI document may describe them.

| Capability | Source-enforced condition in 1.24.0 |
|---|---|
| Organization OIDC providers | Private router. All CRUD/list routes use `verifyValidLicense`; create/update additionally use `verifyValidSubscription(tierMatrix.orgOidc)`. The pinned SaaS tier set is `tier1`, `tier2`, `tier3`, `enterprise`. |
| Site provisioning keys | Private router. All CRUD/list routes check the license and subscription set `tierMatrix.siteProvisioningKeys`: `tier3`, `enterprise`. |
| Shared resource-policy create/list/delete | Private router, license check, and subscription set `tierMatrix.resourcePolicies`: `tier3`, `enterprise`. Their GET/update and access setters are in the base router. |
| Event streaming, alert rules, standalone health checks | Private router required. These Integration route declarations have action/organization checks and some usage-limit checks, but do **not** invoke `verifyValidLicense` or `verifyValidSubscription`. Do not infer an enforced tier from the feature name alone. |
| Organization membership role removal | `DELETE /user/{userId}/remove-role/{roleId}` is registered only in the private router. Its route has no license/subscription middleware. A complete mutable membership lifecycle therefore needs that route. |
| Global OIDC provider | Base router, root key. Creation rejects `IDENTITY_PROVIDER_MODE=org`; configure global identity-provider mode. |

[`verifyValidLicense`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/middlewares/verifyValidLicense.ts) checks license validity only for Enterprise. [`verifyValidSubscription`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/middlewares/verifySubscription.ts) checks active subscription and allowed tiers only for SaaS. The tier sets are from the pinned [`tierMatrix`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/lib/billing/tierMatrix.ts). These are source-level compatibility facts, not a statement of current commercial pricing or entitlements. `verifyLimits` can additionally reject SaaS writes when organization usage exceeds limits.

## Working API Key example

This creates an organization key suitable for ordinary site CRUD and site list data sources. The provider executing **key creation** must use an existing root key with `createApiKey`, `setApiKeyActions`, `listApiKeys`, `listApiKeyActions`, and `deleteApiKey`. The `action_ids` below are granted to the **new key**, not to its creator.

```hcl
resource "pangolin_organization_api_key" "sites" {
  org_id = "example"
  name   = "terraform-sites"

  action_ids = [
    "createSite",
    "getSite",
    "listSites",
    "updateSite",
    "deleteSite",
  ]
}

output "sites_api_credential" {
  value     = "${pangolin_organization_api_key.sites.id}.${pangolin_organization_api_key.sites.api_key}"
  sensitive = true
}
```

Set the provider's `api_key` or `PANGOLIN_API_KEY` to the complete **`id.secret`** credential. The resource's `api_key` attribute is the returned secret portion; it must be combined with its `id` as above. This is the exact format parsed by [Integration API authentication](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/middlewares/integration/verifyApiKey.ts#L25).

For read-only site inventory, use `action_ids = ["getSite", "listSites"]`. Expand the action set using the exact rows above when the same key manages more resource types. The API rejects an empty action set and unknown IDs; this is enforced by [`setApiKeyActions`](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/apiKeys/setApiKeyActions.ts).
