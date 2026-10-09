resource "pangolin_global_oidc_identity_provider" "example" {
  auth_url        = "https://idp.example.com/"
  client_id       = "example"
  client_secret   = "replace-with-a-secure-value"
  identifier_path = "example"
  name            = "example"
  scopes          = "example"
  token_url       = "https://idp.example.com/"
}
