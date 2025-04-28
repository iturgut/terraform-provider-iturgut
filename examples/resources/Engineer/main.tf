 terraform {
  required_providers {
    devops-bootcamp = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops-bootcamp" {
  # example configuration here
  host = "http://localhost:8080"
}

resource "devops-bootcamp_engineer_resource" "example" {
  # No configuration needed - this data source returns all engineers
  name = "Conor"
  email = "conor@liatr.io"
}

# Access the engineers data
output "first_engineer" {
  value = resource.devops-bootcamp_engineer_resource.example.name
}