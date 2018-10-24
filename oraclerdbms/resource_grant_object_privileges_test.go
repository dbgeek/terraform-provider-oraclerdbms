package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGrantObjPrivs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGrantObjPrivsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_grant_object_privilege.grantobjtest", "owner", "SYSTEM"),
				),
			},
			resource.TestStep{
				Config: testAccGrantObjPrivsConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_grant_object_privilege.grantobjtest", "owner", "SYSTEM"),
				),
			},
		},
	})
}
func TestAccGrantObjPrivs_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_grant_object_privilege.grantobjtest"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGrantObjPrivsConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const (
	testAccGrantObjPrivsConfigBasic = `
resource "oraclerdbms_grant_object_privilege" "grantobjtest" {
	grantee = "${oraclerdbms_user.userobjpriv.id}"
	privilege = ["SELECT","UPDATE"]
	owner = "SYSTEM"
	object = "TEST"

}
resource "oraclerdbms_user" "userobjpriv" {
	username = "USER999"
	password = "change_on_install"
	default_tablespace = "USERS"
}
`

	testAccGrantObjPrivsConfigBasic2 = `
resource "oraclerdbms_grant_object_privilege" "grantobjtest" {
	grantee = "${oraclerdbms_user.userobjpriv.id}"
	privilege = ["SELECT"]
	owner = "SYSTEM"
	object = "TEST"

}
resource "oraclerdbms_user" "userobjpriv" {
	username = "USER999"
	password = "change_on_install"
	default_tablespace = "USERS"
}
`
)
