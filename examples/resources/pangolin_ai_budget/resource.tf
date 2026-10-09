resource "pangolin_ai_budget" "example" {
  amount = 1000000
  org_id = "example"
  period = "monthly"
  unit   = "tokens"
}
