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

 data "devops-bootcamp_engineer_data_source" "example" {
   # No configuration needed - this data source returns all engineers
   id = "NFA2Z"
 }

 # Access the engineers data
 output "first_engineer" {
   value = data.devops-bootcamp_engineer_data_source.example.name
 }
   