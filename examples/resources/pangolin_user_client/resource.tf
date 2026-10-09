resource "pangolin_user_client" "example" {
  name    = "example"
  olm_id  = "existing-user-olm-id"
  org_id  = "example"
  type    = "olm"
  user_id = "example-id"
}
