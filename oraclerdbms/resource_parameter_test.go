package oraclerdbms

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccParameter(t *testing.T) {
	fmt.Println("[INFO] TestAccParameter")
	var parameterRsID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccParameterCheckDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccParameterConfigBasic,
				Check: resource.ComposeTestCheckFunc(testAccParameterCheck(
					"oraclerdbms_parameter.test", &parameterRsID),
				),
			},
			resource.TestStep{
				Config: testAccParameterConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					testAccParameterCheck("oraclerdbms_parameter.test", &parameterRsID),
					resource.TestCheckResourceAttr("oraclerdbms_parameter.test", "value", "444"),
				),
			},
		},
	})
}

func testAccParameterCheck(rn string, name *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Println("[INFO] testAccParameterCheck")
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("parameter id not set")
		}

		fmt.Printf("[INFO] rn: %v \n", rn)

		*name = rs.Primary.ID

		return nil
	}
}

func testAccParameterCheckDestroy(s *terraform.State) error {

	fmt.Println("[INFO] testAccParameterCheckDestroy")
	for _, rs := range s.RootModule().Resources {
		fmt.Printf("[INFO] resource type: %v \n", rs.Primary.Attributes)
		if rs.Type != "oraclerdbms_parameter" {
			fmt.Println("[INFO] rs type != oraclerdbms_parameter")
			continue
		}
	}
	return nil
}

func TestAccParameter_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_parameter.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccParameterConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccParameterConfigBasic = `
resource "oraclerdbms_parameter" "test" {
    name = "undo_retention"
	value = "666"
}
`

const testAccParameterConfigBasic2 = `
resource "oraclerdbms_parameter" "test" {
    name = "undo_retention"
	value = "444"
	update_comment = "acc test of comment"
}
`
