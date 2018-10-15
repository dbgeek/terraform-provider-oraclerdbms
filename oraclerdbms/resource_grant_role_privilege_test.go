package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGrantRolePrivs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGrantRolePrivsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_grant_role_privilege.grantroleprivs", "role", "ROLEPRIVS"),
				),
			},
		},
	})
}

const testAccGrantRolePrivsConfigBasic = `
resource "oraclerdbms_grant_role_privilege" "grantroleprivs" {
	grantee = "${oraclerdbms_user.userrolepriv.id}"
	role = "${oraclerdbms_role.roleprivstest.id}"

}
resource "oraclerdbms_user" "userrolepriv" {
	username = "USERROLEPRIVS"
	password = "change_on_install"
	default_tablespace = "USERS"
}

resource "oraclerdbms_role" "roleprivstest" {
	role = "ROLEPRIVS"
}
`
