resource "pangolin_organization_user" "example" {
  idp_id   = "1"
  org_id   = "example"
  role_ids = ["1"]
  type     = "oidc"
  username = "user@example.com"
}
