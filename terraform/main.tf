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