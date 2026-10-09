resource "pangolin_alert_rule" "example" {
  all_sites  = true
  event_type = "site_offline"
  name       = "example"
  org_id     = "example"
}
