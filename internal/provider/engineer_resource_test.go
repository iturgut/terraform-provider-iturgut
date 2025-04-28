// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccEngineerResource tests the basic CRUD operations for the engineer resource
func TestAccEngineerResource(t *testing.T) {
	// Skip this test if we're not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	// Use a unique email for this test to avoid conflicts
	testEmail := fmt.Sprintf("test.%d@example.com", os.Getpid())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngineerResourceConfig("Test Engineer", testEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "name", "Test Engineer"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer_resource.test", "email", testEmail),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer_resource.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_engineer_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Helper function to generate test configuration for a single engineer resource
func testAccEngineerResourceConfig(name, email string) string {
	return fmt.Sprintf(`
resource "devops-bootcamp_engineer_resource" "test" {
  name  = %[1]q
  email = %[2]q
}
`, name, email)
}
