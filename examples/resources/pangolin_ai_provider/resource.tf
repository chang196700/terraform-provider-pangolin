resource "pangolin_ai_provider" "example" {
  api_key      = "replace-with-your-upstream-key"
  name         = "example"
  org_id       = "example"
  type         = "openai"
  upstream_url = "https://api.openai.com/v1"
}
