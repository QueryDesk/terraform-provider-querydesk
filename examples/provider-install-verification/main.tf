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
  username = "postgres"
  password = "postgres"
  adapter  = "postgres"
  hostname = "localhost"
  database = "coverbot_dev"
  # TODO: figure out computed values
  ssl              = false
  restrict_access  = false
  reviews_required = 0
}

output "database" {
  value = querydesk_database.example
}
