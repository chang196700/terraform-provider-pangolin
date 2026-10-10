# Integration API coverage

Source: unchanged `openapi.yaml` (SHA-256 `d9cbd22791ecd449182a3d8778a39bb0e36b5c96e62c983abc859368a700ab52`). Behavior baseline: Pangolin 1.24.0, `a235635d9e5a6431ea426431e3e862e965928cd3`.

Every operation is classified below. `Mapped` means the behavior is owned by the named Terraform entity; it does not mean each alternate API route is called. `Alternative` routes expose behavior through another route or lookup. `Excluded` operations are not persistent managed state. This inventory is not evidence of live server validation.

| Method | Path | Disposition | Terraform owner / rationale |
|---|---|---|---|
| GET | `/` | Excluded | Operational telemetry or API discovery. |
| DELETE | `/access-token/{accessTokenId}` | Mapped | public_resource_access_token |
| DELETE | `/admin/org/{orgId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/admin/orgs` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| DELETE | `/ai-budget/{budgetId}` | Mapped | ai_budget |
| GET | `/ai-budget/{budgetId}` | Mapped | ai_budget |
| POST | `/ai-budget/{budgetId}` | Mapped | ai_budget |
| DELETE | `/ai-model/{modelId}` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| GET | `/ai-model/{modelId}` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| POST | `/ai-model/{modelId}` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| GET | `/ai-model/{modelId}/ai-budgets` | Mapped | ai_budget |
| DELETE | `/ai-provider/{providerId}` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| GET | `/ai-provider/{providerId}` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| POST | `/ai-provider/{providerId}` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| GET | `/ai-provider/{providerId}/ai-budgets` | Mapped | ai_budget |
| GET | `/ai-provider/{providerId}/catalog-models` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| PUT | `/ai-provider/{providerId}/model` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| GET | `/ai-provider/{providerId}/models` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| PUT | `/ai-provider/{providerId}/target` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| GET | `/ai-provider/{providerId}/targets` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| POST | `/certificate/{certId}` | Excluded | Certificate issuance/restart action; certificate metadata is read-only. |
| DELETE | `/client/{clientId}` | Mapped | client / user_client / user_devices |
| GET | `/client/{clientId}` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}/archive` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}/block` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}/private-resources` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}/rebuild-associations-cache` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| POST | `/client/{clientId}/site-resources` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/client/{clientId}/unarchive` | Mapped | client / user_client / user_devices |
| POST | `/client/{clientId}/unblock` | Mapped | client / user_client / user_devices |
| GET | `/client/{clientId}/verify-associations-cache` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| GET | `/domain/check-namespace-availability` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| GET | `/domains/namepaces` | Unavailable in 1.24.0 | Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md. |
| GET | `/idp` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| PUT | `/idp/oidc` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| DELETE | `/idp/{idpId}` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| GET | `/idp/{idpId}` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| POST | `/idp/{idpId}/oidc` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| GET | `/idp/{idpId}/org` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| DELETE | `/idp/{idpId}/org/{orgId}` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| POST | `/idp/{idpId}/org/{orgId}` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| PUT | `/idp/{idpId}/org/{orgId}` | Mapped | global_oidc_identity_provider / identity_provider_organization |
| GET | `/maintenance/info` | Excluded | Operational telemetry or API discovery. |
| GET | `/openapi.json` | Excluded | Operational telemetry or API discovery. |
| GET | `/openapi.yaml` | Excluded | Operational telemetry or API discovery. |
| PUT | `/org` | Mapped | organization |
| DELETE | `/org/{orgId}` | Upstream defect | organization: stock 1.24.0 handler dereferences session user after Bearer authentication; see compatibility.md. |
| GET | `/org/{orgId}` | Mapped | organization |
| POST | `/org/{orgId}` | Mapped | organization |
| GET | `/org/{orgId}/access-tokens` | Mapped | public_resource_access_token |
| PUT | `/org/{orgId}/ai-budget` | Mapped | ai_budget |
| GET | `/org/{orgId}/ai-budgets` | Mapped | ai_budget |
| GET | `/org/{orgId}/ai-catalog-models` | Mapped | ai_model / ai_catalog_models / resource AI associations |
| PUT | `/org/{orgId}/ai-provider` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| GET | `/org/{orgId}/ai-provider/{niceId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/ai-providers` | Mapped | ai_provider / public_resource.ai_providers / private_resource.ai_providers |
| PUT | `/org/{orgId}/alert-rule` | Mapped | alert_rule |
| DELETE | `/org/{orgId}/alert-rule/{alertRuleId}` | Mapped | alert_rule |
| GET | `/org/{orgId}/alert-rule/{alertRuleId}` | Mapped | alert_rule |
| POST | `/org/{orgId}/alert-rule/{alertRuleId}` | Mapped | alert_rule |
| GET | `/org/{orgId}/alert-rules` | Mapped | alert_rule |
| PUT | `/org/{orgId}/api-key` | Mapped | organization_api_key |
| DELETE | `/org/{orgId}/api-key/{apiKeyId}` | Mapped | organization_api_key |
| GET | `/org/{orgId}/api-key/{apiKeyId}/actions` | Mapped | organization_api_key |
| POST | `/org/{orgId}/api-key/{apiKeyId}/actions` | Mapped | organization_api_key |
| GET | `/org/{orgId}/api-keys` | Mapped | organization_api_key |
| PUT | `/org/{orgId}/blueprint` | Excluded | Blueprint apply has no corresponding destroy lifecycle. |
| GET | `/org/{orgId}/blueprint/{blueprintId}` | Mapped | blueprint / blueprints data sources |
| GET | `/org/{orgId}/blueprints` | Mapped | blueprint / blueprints data sources |
| GET | `/org/{orgId}/certificate/{domainId}/{domain}` | Unavailable in 1.24.0 | Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md. |
| PUT | `/org/{orgId}/client` | Mapped | client / user_client / user_devices |
| GET | `/org/{orgId}/client/{niceId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/clients` | Mapped | client / user_client / user_devices |
| POST | `/org/{orgId}/create-invite` | Excluded | Invitation consumption and delivery are workflows; invitations are read-only. |
| GET | `/org/{orgId}/domain/{domainId}` | Mapped | domain / domains / domain_dns_records data sources |
| PATCH | `/org/{orgId}/domain/{domainId}` | Mapped | domain |
| POST | `/org/{orgId}/domain/{domainId}` | Mapped | domain |
| GET | `/org/{orgId}/domain/{domainId}/dns-records` | Mapped | domain / domains / domain_dns_records data sources |
| GET | `/org/{orgId}/domains` | Mapped | domain / domains / domain_dns_records data sources |
| GET | `/org/{orgId}/event-streaming-destination` | Corrected | event_streaming_destination: runtime route is /org/{orgId}/event-streaming-destinations (plural). |
| PUT | `/org/{orgId}/event-streaming-destination` | Mapped | event_streaming_destination |
| DELETE | `/org/{orgId}/event-streaming-destination/{destinationId}` | Mapped | event_streaming_destination |
| POST | `/org/{orgId}/event-streaming-destination/{destinationId}` | Mapped | event_streaming_destination |
| PUT | `/org/{orgId}/health-check` | Mapped | health_check |
| DELETE | `/org/{orgId}/health-check/{healthCheckId}` | Mapped | health_check |
| POST | `/org/{orgId}/health-check/{healthCheckId}` | Mapped | health_check |
| GET | `/org/{orgId}/health-checks` | Mapped | health_check |
| GET | `/org/{orgId}/idp` | Mapped | organization_oidc_identity_provider |
| PUT | `/org/{orgId}/idp/oidc` | Mapped | organization_oidc_identity_provider |
| DELETE | `/org/{orgId}/idp/{idpId}` | Mapped | organization_oidc_identity_provider |
| GET | `/org/{orgId}/idp/{idpId}` | Mapped | organization_oidc_identity_provider |
| POST | `/org/{orgId}/idp/{idpId}/oidc` | Mapped | organization_oidc_identity_provider |
| GET | `/org/{orgId}/invitations` | Mapped | invitations data source |
| DELETE | `/org/{orgId}/invitations/{inviteId}` | Excluded | Invitation consumption and delivery are workflows; invitations are read-only. |
| GET | `/org/{orgId}/logs/access` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/access/export` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/action` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/action/export` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/export` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/filters` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/overview` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/providers` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/resources` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/users-roles` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/ai/usage/virtual-api-keys` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/analytics` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/connection` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/connection/export` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/logs/request` | Excluded | Operational telemetry or API discovery. |
| GET | `/org/{orgId}/my-virtual-api-keys` | Excluded | Dashboard personal projection requires a user session; virtual_api_key uses organization-scoped Integration routes. |
| GET | `/org/{orgId}/my-virtual-api-keys/{virtualApiKeyId}` | Excluded | Dashboard personal projection requires a user session; virtual_api_key uses organization-scoped Integration routes. |
| GET | `/org/{orgId}/pick-client-defaults` | Helper | client: called only during creation when defaults are needed. |
| GET | `/org/{orgId}/pick-site-defaults` | Excluded | Preview allocation is not persistent; Newt creation allocates credentials. |
| PUT | `/org/{orgId}/private-resource` | Mapped | private_resource |
| GET | `/org/{orgId}/private-resources` | Mapped | private_resource |
| PUT | `/org/{orgId}/public-resource` | Mapped | public_resource |
| GET | `/org/{orgId}/public-resource-names` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/public-resources` | Mapped | public_resource |
| POST | `/org/{orgId}/reset-bandwidth` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| PUT | `/org/{orgId}/resource` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/resource-policies` | Mapped | public_resource_policy / public_resource.access |
| POST | `/org/{orgId}/resource-policy` | Mapped | public_resource_policy / public_resource.access |
| GET | `/org/{orgId}/resource-policy/{niceId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/resource/{niceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/resources` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/resources-names` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/org/{orgId}/role` | Mapped | role / organization_user |
| GET | `/org/{orgId}/roles` | Mapped | role / organization_user |
| PUT | `/org/{orgId}/site` | Mapped | site |
| PUT | `/org/{orgId}/site-provisioning-key` | Mapped | site_provisioning_key |
| DELETE | `/org/{orgId}/site-provisioning-key/{siteProvisioningKeyId}` | Mapped | site_provisioning_key |
| PATCH | `/org/{orgId}/site-provisioning-key/{siteProvisioningKeyId}` | Mapped | site_provisioning_key |
| GET | `/org/{orgId}/site-provisioning-keys` | Mapped | site_provisioning_key |
| PUT | `/org/{orgId}/site-resource` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/site-resources` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/site/{niceId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/site/{siteId}/private-resources` | Mapped | private_resource |
| GET | `/org/{orgId}/site/{siteId}/resource/nice/{niceId}` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/site/{siteId}/resources` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/org/{orgId}/sites` | Mapped | site |
| PUT | `/org/{orgId}/user` | Mapped | organization_user / user / users / user_server_admin |
| GET | `/org/{orgId}/user-by-username` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| GET | `/org/{orgId}/user-devices` | Mapped | client / user_client / user_devices |
| DELETE | `/org/{orgId}/user/{userId}` | Mapped | organization_user / user / users / user_server_admin |
| GET | `/org/{orgId}/user/{userId}` | Alternative | organization_user: handler requires session context; use org users and username lookup with Bearer authentication. |
| POST | `/org/{orgId}/user/{userId}` | Mapped | organization_user |
| GET | `/org/{orgId}/user/{userId}/check` | Alternative | Alternate lookup or admin projection; Terraform identifies entities by immutable ID. |
| PUT | `/org/{orgId}/user/{userId}/client` | Unavailable in 1.24.0 | Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md. |
| GET | `/org/{orgId}/users` | Mapped | organization_user / user / users / user_server_admin |
| PUT | `/org/{orgId}/virtual-api-key` | Mapped | virtual_api_key |
| GET | `/org/{orgId}/virtual-api-keys` | Mapped | virtual_api_key |
| POST | `/org/{orgId}/virtual-api-keys/email-identity-keys` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| GET | `/orgs` | Mapped | organization |
| GET | `/private-resource/{resourceId}/clients` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/clients` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/clients/add` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/clients/remove` | Mapped | private_resource |
| GET | `/private-resource/{resourceId}/roles` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/roles` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/roles/add` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/roles/remove` | Mapped | private_resource |
| GET | `/private-resource/{resourceId}/users` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/users` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/users/add` | Mapped | private_resource |
| POST | `/private-resource/{resourceId}/users/remove` | Mapped | private_resource |
| DELETE | `/private-resource/{siteResourceId}` | Mapped | private_resource |
| GET | `/private-resource/{siteResourceId}` | Mapped | private_resource |
| POST | `/private-resource/{siteResourceId}` | Mapped | private_resource |
| POST | `/public-resource-policy/{resourceId}/access-control` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourceId}/access-control` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/public-resource-policy/{resourcePolicyId}` | Mapped | public_resource_policy / public_resource.access |
| POST | `/public-resource-policy/{resourcePolicyId}` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/public-resource-policy/{resourcePolicyId}/header-auth` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}/header-auth` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/public-resource-policy/{resourcePolicyId}/password` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}/password` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/public-resource-policy/{resourcePolicyId}/pincode` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}/pincode` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/public-resource-policy/{resourcePolicyId}/rules` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}/rules` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/public-resource-policy/{resourcePolicyId}/whitelist` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/public-resource-policy/{resourcePolicyId}/whitelist` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/public-resource/{resourceId}` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/access-token` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/access-tokens` | Alternative | public_resource_access_token: use the organization token list and filter resourceId; the direct route fails for root keys. |
| POST | `/public-resource/{resourceId}/header-auth` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/password` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/pincode` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/policies` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/roles` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/roles` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/roles/add` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/roles/remove` | Mapped | public_resource |
| PUT | `/public-resource/{resourceId}/rule` | Mapped | public_resource |
| DELETE | `/public-resource/{resourceId}/rule/{ruleId}` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/rule/{ruleId}` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/rules` | Mapped | public_resource |
| PUT | `/public-resource/{resourceId}/target` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/targets` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/users` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/users` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/users/add` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/users/remove` | Mapped | public_resource |
| GET | `/public-resource/{resourceId}/whitelist` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/whitelist` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/whitelist/add` | Mapped | public_resource |
| POST | `/public-resource/{resourceId}/whitelist/remove` | Mapped | public_resource |
| PUT | `/resource-policy/{resourceId}/access-control` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/resource-policy/{resourcePolicyId}` | Mapped | public_resource_policy / public_resource.access |
| GET | `/resource-policy/{resourcePolicyId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}` | Mapped | public_resource_policy / public_resource.access |
| PUT | `/resource-policy/{resourcePolicyId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/access-control` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}/header-auth` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/header-auth` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}/password` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/password` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}/pincode` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/pincode` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}/rules` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/rules` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource-policy/{resourcePolicyId}/whitelist` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource-policy/{resourcePolicyId}/whitelist` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/resource/{resourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/access-token` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/access-tokens` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/ai-budgets` | Mapped | ai_budget |
| GET | `/resource/{resourceId}/ai-models` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-models` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-models/add` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-models/remove` | Mapped | public_resource |
| GET | `/resource/{resourceId}/ai-providers` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-providers` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-providers/add` | Mapped | public_resource |
| POST | `/resource/{resourceId}/ai-providers/remove` | Mapped | public_resource |
| POST | `/resource/{resourceId}/header-auth` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/password` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/pincode` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/policies` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/roles` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/roles` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/roles/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/roles/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource/{resourceId}/rule` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/resource/{resourceId}/rule/{ruleId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/rule/{ruleId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/rules` | Alternative | Legacy alias; use the canonical Terraform resource. |
| PUT | `/resource/{resourceId}/target` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/targets` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/users` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/users` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/users/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/users/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/resource/{resourceId}/whitelist` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/whitelist` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/whitelist/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/resource/{resourceId}/whitelist/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/role/{roleId}` | Mapped | role / organization_user |
| GET | `/role/{roleId}` | Mapped | role / organization_user |
| POST | `/role/{roleId}` | Mapped | role / organization_user |
| POST | `/role/{roleId}/add/{userId}` | Mapped | role / organization_user |
| GET | `/role/{roleId}/ai-budgets` | Mapped | ai_budget |
| DELETE | `/site-resource/{siteResourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/site-resource/{siteResourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/site-resource/{siteResourceId}/ai-budgets` | Mapped | ai_budget |
| GET | `/site-resource/{siteResourceId}/ai-models` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-models` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-models/add` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-models/remove` | Mapped | private_resource |
| GET | `/site-resource/{siteResourceId}/ai-providers` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-providers` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-providers/add` | Mapped | private_resource |
| POST | `/site-resource/{siteResourceId}/ai-providers/remove` | Mapped | private_resource |
| GET | `/site-resource/{siteResourceId}/clients` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/clients` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/clients/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/clients/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/site-resource/{siteResourceId}/roles` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/roles` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/roles/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/roles/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| GET | `/site-resource/{siteResourceId}/users` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/users` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/users/add` | Alternative | Legacy alias; use the canonical Terraform resource. |
| POST | `/site-resource/{siteResourceId}/users/remove` | Alternative | Legacy alias; use the canonical Terraform resource. |
| DELETE | `/site/{siteId}` | Mapped | site |
| GET | `/site/{siteId}` | Mapped | site |
| POST | `/site/{siteId}` | Mapped | site |
| POST | `/site/{siteId}/approve` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| POST | `/site/{siteId}/reject` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| POST | `/site/{siteId}/restart` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| DELETE | `/target/{targetId}` | Mapped | resource_target |
| GET | `/target/{targetId}` | Mapped | resource_target |
| POST | `/target/{targetId}` | Mapped | resource_target |
| GET | `/user/{userId}` | Mapped | organization_user / user / users / user_server_admin |
| POST | `/user/{userId}/2fa` | Excluded | Imperative workflow; not a symmetric Terraform lifecycle. |
| POST | `/user/{userId}/add-role/{roleId}` | Mapped | role / organization_user |
| DELETE | `/user/{userId}/remove-role/{roleId}` | Mapped | role / organization_user |
| POST | `/user/{userId}/server-admin` | Unavailable in 1.24.0 | Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md. |
| GET | `/users` | Unavailable in 1.24.0 | Declared by OpenAPI but not mounted on the Bearer Integration API; see compatibility.md. |
| DELETE | `/virtual-api-key/{virtualApiKeyId}` | Mapped | virtual_api_key |
| GET | `/virtual-api-key/{virtualApiKeyId}` | Mapped | virtual_api_key |
| POST | `/virtual-api-key/{virtualApiKeyId}` | Mapped | virtual_api_key |
| GET | `/virtual-api-key/{virtualApiKeyId}/ai-budgets` | Mapped | ai_budget |

Total: **312 operations** in the supplied OpenAPI document.

## Fixed-source routes omitted from `openapi.yaml`

The pinned Pangolin source registers the following Bearer Integration routes even though they are absent from the supplied OpenAPI document. The domain lifecycle routes are implemented by `pangolin_domain`; the resource-name route is exposed as `pangolin_resource_names`. The restart endpoint remains an imperative operation and is intentionally not modeled as Terraform state.

| Method | Path | Provider disposition |
|---|---|---|
| PUT | `/org/{orgId}/domain` | Mapped: `pangolin_domain` create |
| DELETE | `/org/{orgId}/domain/{domainId}` | Mapped: `pangolin_domain` destroy |
| POST | `/org/{orgId}/domain/{domainId}/restart` | Excluded: imperative restart |
| GET | `/org/{orgId}/resource-names` | Mapped: `pangolin_resource_names` |
| GET | `/org/{orgId}/public-resource-names` | Alias of `resource-names` |

The domain update `POST` and deprecated `PATCH` routes are present in the OpenAPI document and are mapped to `pangolin_domain`; their runtime request body is supplied by the fixed server source.
