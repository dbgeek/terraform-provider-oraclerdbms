package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGrantSysPrivs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGrantSysPrivsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_grant_system_privilege.grantsysprivs", "privilege", "CREATE SESSION"),
				),
			},
		},
	})
}

const testAccGrantSysPrivsConfigBasic = `
resource "oraclerdbms_grant_system_privilege" "grantsysprivs" {
	grantee = "${oraclerdbms_user.usersyspriv.id}"
	privilege = "CREATE SESSION"

}
resource "oraclerdbms_user" "usersyspriv" {
	username = "USER999"
	password = "change_on_install"
	default_tablespace = "USERS"
}
`
