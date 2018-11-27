package oraclerdbms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

func TestAccGrantObjPrivsOnSchema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testGrantOnSchemaAddCleanUp(),
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGrantObjPrivsOnSchemaConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_grant_object_privilege.grantobjschematest", "owner", "SYSTEM"),
				),
			},
		},
	})
}
func TestAccGrantObjPrivs_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_grant_object_privilege.grantobjtest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
func testDBVersion() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		oracleHelper := testAccProvider.Meta().(*oracleHelperType)
		fmt.Printf("[DEBUG] This Dbversion %s\n", oracleHelper.Client.DBVersion)
		return nil
	}
}
func testGrantOnSchemaAddSetup() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*oracleHelperType).Client
		client.DBClient.Exec("CREATE TABLE SYSTEM.TST_TBL1(col number)")
		return nil
	}
}
func testGrantOnSchemaAddTable() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*oracleHelperType).Client
		client.DBClient.Exec("CREATE TABLE SYSTEM.TST_TBL2(col number)")
		return nil
	}
}

func testGrantOnSchemaAddCleanUp() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*oracleHelperType).Client
		client.DBClient.Exec("DROP TABLESYSTEM.TST_TBL1")
		client.DBClient.Exec("DROP TABLE SYSTEM.TST_TBL2")
		return nil
	}
}

func testAccGrantObjPrivsOnSchemaConfigBasicDiff() string {
	return `
	resource "oraclerdbms_grant_object_privilege" "grantobjschematest" {
		grantee = "${oraclerdbms_user.userobjpriv2.id}"
		privilege = ["SELECT","UPDATE"]
		owner = "SYSTEM"
		object_type = "TABLE"
	}
	resource "oraclerdbms_user" "userobjpriv2" {
		username = "USER99"
		default_tablespace = "USERS"
	}
	`
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
	default_tablespace = "USERS"
}
`
	testAccGrantObjPrivsOnSchemaConfigBasic = `
resource "oraclerdbms_grant_object_privilege" "grantobjschematest" {
	grantee = "${oraclerdbms_user.userobjpriv2.id}"
	privilege = ["SELECT","UPDATE"]
	owner = "SYSTEM"
	object_type = "TABLE"
}
resource "oraclerdbms_user" "userobjpriv2" {
	username = "USER99"
	default_tablespace = "USERS"
}
`

	testAccGrantObjPrivsOnSchemaConfigBasic2 = `
resource "oraclerdbms_grant_object_privilege" "grantobjschematest" {
	grantee = "${oraclerdbms_user.userobjpriv2.id}"
	privilege = ["SELECT"]
	owner = "SYSTEM"
	object_type = "TABLE"
}
resource "oraclerdbms_user" "userobjpriv2" {
	username = "USER99"
	default_tablespace = "USERS"
}
`
)
