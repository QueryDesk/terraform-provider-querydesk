terraform {
  required_providers {
    querydesk = {
      source = "registry.terraform.io/querydesk/querydesk"
    }
  }
}

provider "querydesk" {
  api_key = "SFMyNTY.g2gDbQAAAB5rZXlfMDFIM0JFWjlENkJSMVc1NUcwSjk5TUswMktuBgDGpe_WiAFiAAFRgA.LiCcHky6wRmzciNtrP2vSQzz4QEvv9qL255BinhfF7I"
  host    = "http://localhost:4000"
}

resource "querydesk_database" "example" {
  name     = "terraform"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "terraform"
}

resource "querydesk_database_user" "example" {
  database_id      = querydesk_database.example.id
  username         = "postgres"
  password         = "postgres"
  reviews_required = 0
}
