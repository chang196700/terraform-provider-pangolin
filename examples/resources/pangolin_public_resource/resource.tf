resource "pangolin_public_resource" "example" {
  domain_id = "your-domain-id"
  mode      = "http"
  name      = "example"
  org_id    = "example"
  subdomain = "app"
}
