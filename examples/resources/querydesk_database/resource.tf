resource "querydesk_database" "example" {
  name     = "terraform_test"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"
}
