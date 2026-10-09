resource "pangolin_organization" "example" {
  name           = "example"
  org_id         = "example"
  subnet         = "100.90.0.0/16"
  utility_subnet = "100.91.0.0/16"
}
