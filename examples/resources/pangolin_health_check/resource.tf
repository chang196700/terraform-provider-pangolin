resource "pangolin_health_check" "example" {
  hc_enabled  = true
  hc_hostname = "127.0.0.1"
  hc_mode     = "http"
  hc_port     = 8080
  name        = "example"
  org_id      = "example"
  site_id     = "1"
}
