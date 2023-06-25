resource "querydesk_database" "example" {
  name     = "terraform_test"
  adapter  = "postgres"
  hostname = "localhost"
  database = "mydb"
}
