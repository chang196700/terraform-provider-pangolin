resource "pangolin_private_resource" "example" {
  destination           = "10.0.0.10"
  mode                  = "host"
  name                  = "example"
  org_id                = "example"
  site_ids              = ["1"]
  tcp_port_range_string = "443"
  udp_port_range_string = ""
}
