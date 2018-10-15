package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_user.test", "username", "USER666"),
				),
			},
		},
	})
}

const testAccUserConfigBasic = `
resource "oraclerdbms_user" "test" {
	username = "USER666"
	password = "change_on_install"
	default_tablespace = "USERS"
}
`
