resource "pangolin_organization_api_key" "example" {
  action_ids = ["createSite", "getSite", "listSites", "updateSite", "deleteSite"]
  name       = "example"
  org_id     = "example"
}
