 terraform {
   required_providers {
     devops-bootcamp = {
       source = "liatr.io/terraform/devops-bootcamp"
     }
   }
 }

 provider "devops-bootcamp" {
   host = "http://localhost:8080"
 }

 resource "devops-bootcamp_engineer_resource" "example" {
  name = "Conor"
  email = "conor@liatr.io"
}

 data "devops-bootcamp_engineer_data_source" "example" {
   id = "G6R3N"
 }

 output "first_engineer" {
   value = data.devops-bootcamp_engineer_data_source.example
 }
   