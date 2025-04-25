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

 data "devops-bootcamp_engineer_data_source" "example" {
  // gets the engineer with id "NFA2Z"
   id = "NFA2Z"
 }

 output "first_engineer" {
  // returns the name of the engineer with id "NFA2Z"
   value = data.devops-bootcamp_engineer_data_source.example.name
 }
   