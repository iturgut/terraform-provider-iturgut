// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEngineerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a resource first
			{
				Config: testAccEngineerDataSourceSetup,
			},
			// Then test the data source
			{
				Config: testAccEngineerDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.devops-bootcamp_engineer_data_source.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Test Engineer"),
					),
				},
			},
		},
	})
}

const testAccEngineerDataSourceSetup = `
resource "devops-bootcamp_engineer_resource" "test" {
  name  = "Test Engineer"
  email = "test.datasource@example.com"
}
`

const testAccEngineerDataSourceConfig = `
resource "devops-bootcamp_engineer_resource" "test" {
  name  = "Test Engineer"
  email = "test.datasource@example.com"
}

data "devops-bootcamp_engineer_data_source" "test" {
  id = devops-bootcamp_engineer_resource.test.id
}
`
