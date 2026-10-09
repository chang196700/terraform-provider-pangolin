resource "pangolin_resource_target" "example" {
  ip          = "127.0.0.1"
  method      = "http"
  port        = 8080
  resource_id = "1"
  site_id     = "1"
}
