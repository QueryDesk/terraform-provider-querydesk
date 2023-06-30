resource "querydesk_database" "example" {
  name     = "terraform_test"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"
}

resource "querydesk_database_user" "example" {
  database_id      = querydesk_database.example.id
  username         = "postgres"
  password         = "postgres"
  reviews_required = 0
}