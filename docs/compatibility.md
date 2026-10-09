# Pangolin API compatibility

The behavior baseline is **Pangolin 1.24.0**, commit [`a235635d9e5a6431ea426431e3e862e965928cd3`](https://github.com/fosrl/pangolin/tree/a235635d9e5a6431ea426431e3e862e965928cd3). The supplied `openapi.yaml` is preserved unchanged. Its SHA-256 is `d9cbd22791ecd449182a3d8778a39bb0e36b5c96e62c983abc859368a700ab52`.

An OpenAPI registration is not proof that a route is available on the Integration API. The server collects handler registrations in a shared registry, while the [Integration server](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/integrationApiServer.ts) mounts only the routes in the [base integration router](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts) and, for commercial builds, its [private extension](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/routers/integration.ts).

## OpenAPI declarations unavailable on stock 1.24.0

These Terraform types retain the supplied schema for servers that implement those declarations. They are **not usable as complete API Key workflows on stock Pangolin 1.24.0**. API fixtures test their schema and wire behavior, not server availability. A root key, additional action IDs, or a commercial edition does not add a missing route.

| Terraform type | Declared operation | Stock 1.24.0 finding |
|---|---|---|
| Resource `pangolin_user_server_admin` | `POST /user/{userId}/server-admin` | No Integration route. The dashboard registers this operation only inside its non-SaaS server-administrator routes. Create, update, and destroy demotion are unavailable through API Key authentication. |
| Resource `pangolin_user_client` | `PUT /org/{orgId}/user/{userId}/client` | The Integration registration is explicitly commented out. No active route invokes `createUserClient`, so create is unavailable. The shared delete handler explicitly rejects clients with a `userId`, so destroy is also unavailable. Existing user clients can be imported, read, updated, archived, or blocked. |
| Data source `pangolin_users` | `GET /users` | No Integration route. Only the non-SaaS dashboard administrator route registers this global list. `GET /org/{orgId}/users` is a different, available organization-scoped operation. |
| Data source `pangolin_certificate` | `GET /org/{orgId}/certificate/{domainId}/{domain}` | No Integration route. The dashboard route requires a user session, certificate access, and the user's `getCertificate` action. |

Sources: the [commented user-client registration](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts#L1434), [dashboard global users and administrator operations](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/external.ts#L1373), and [dashboard certificate route](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/external.ts#L1891).

The [client delete handler](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/client/deleteClient.ts#L93) returns HTTP 400 for user-associated clients and requires archiving instead, even though the generic delete route exists. This does not affect machine clients managed with `pangolin_client`.

The single-user data source `pangolin_user` **is available**: `GET /user/{userId}` is registered in the Integration router with both root-key and `getUser` checks. This differs from the unavailable global user list. See the [actual registration](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/integration.ts#L1242).

Changing the endpoint from `/v1` to `/api/v1` does not solve these gaps. The [dashboard server](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/apiServer.ts#L31) uses `/api/v1`; its [authenticated router](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/external.ts#L82) requires `verifySessionUserMiddleware`. The [session verifier](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/auth/sessions/verifySession.ts#L7) reads a session cookie, not the Integration API bearer key. This provider implements API Key authentication only.

## Corrected event destination list route

The supplied OpenAPI declares `GET /org/{orgId}/event-streaming-destination` with a singular final segment. The actual [1.24.0 Integration router](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/private/routers/integration.ts#L353) mounts **`GET /org/{orgId}/event-streaming-destinations`**. The provider uses the working plural route for refresh and list data sources; create remains `PUT /org/{orgId}/event-streaming-destination`. The acceptance fixture verifies this distinction.

## Registered routes with session-context defects

| Operation | Baseline defect | Provider behavior |
|---|---|---|
| `DELETE /org/{orgId}` | The registered handler dereferences `req.user.userId`, but bearer authentication supplies `req.apiKey`, not a user session. It can return HTTP 500 for API Key deletion. | `pangolin_organization` retains its declared destroy operation and reports the failure. Organization destroy is not supported end-to-end on the stock baseline; there is no API Key admin-delete alternative. Creation, read, and updates have separate handlers. |
| `GET /org/{orgId}/user/{userId}` | The handler requires `req.userOrg`; the API Key middleware provides `req.apiKeyOrg`. | Membership refresh uses the organization list to identify the immutable user ID, then `/org/{orgId}/user-by-username` to read complete membership data. This needs both `listUsers` and `getOrgUser`. |
| `GET /public-resource/{resourceId}/access-tokens` | The list handler needs organization context, but the root-key resource-access middleware returns before setting it. | The token resource explicitly requires `org_id`, reads `/org/{orgId}/access-tokens`, and filters by `resource_id`. Import uses `org_id/resource_id/id`. |

Sources: [organization delete handler](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/org/deleteOrg.ts#L69), [API Key authentication context](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/middlewares/integration/verifyApiKey.ts#L48), [organization membership handler](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/user/getOrgUser.ts#L135), [access-token list handler](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/routers/accessToken/listAccessTokens.ts), and [resource-access middleware](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/server/middlewares/integration/verifyApiKeyResourceAccess.ts).

## Edition and validation boundary

The build selects the base router for OSS and the private extension for Enterprise/SaaS; this is explicit in the [build import resolver](https://github.com/fosrl/pangolin/blob/a235635d9e5a6431ea426431e3e862e965928cd3/esbuild.mjs#L185). A shared policy lifecycle, organization OIDC provider, site provisioning key, event streaming destination, alert rule, and standalone health check use private routes. Organization membership role removal also uses a private route. See [exact actions and edition checks](permissions.md).

Local Terraform acceptance tests run the real Terraform executable against independent HTTP fixtures. They verify configuration planning, response mapping, updates, import, disappearance, and destroy. They do not establish connectivity to a deployed Pangolin server, licensed feature availability, or behavior beyond the pinned contract. Optional live tests require a separately configured disposable server.

## Lost write responses

Writes are not retried automatically. If the server commits a mutation but its response is lost, inspect the remote object before retrying. In particular, a virtual key expiry update accepts a relative duration without an idempotency key; retrying that same expiry change can calculate a later absolute expiry. Ordinary updates with unchanged duration do not send an expiry change.
