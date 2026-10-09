resource "pangolin_organization_oidc_identity_provider" "example" {
  auth_url        = "https://idp.example.com/"
  client_id       = "example"
  client_secret   = "replace-with-a-secure-value"
  identifier_path = "example"
  name            = "example"
  org_id          = "example"
  scopes          = "example"
  token_url       = "https://idp.example.com/"
}
