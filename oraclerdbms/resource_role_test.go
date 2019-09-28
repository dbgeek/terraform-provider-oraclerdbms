package oraclerdbms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_role.roletest", "role", "ROLE666"),
				),
			},
		},
	})
}

func TestAccRole_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_role.roletest"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				//ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

const testAccRoleConfigBasic = `
resource "oraclerdbms_role" "roletest" {
	role = "ROLE666"

}
`
