package oraclerdbms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
			resource.TestStep{
				Config: testAccUserConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_user.test", "default_tablespace", "SYSTEM"),
					resource.TestCheckResourceAttr("oraclerdbms_user.test", "account_status", "LOCKED"),
					resource.TestCheckResourceAttr("oraclerdbms_user.test", "quota.%", "1"),
					resource.TestCheckResourceAttr("oraclerdbms_user.test", "quota.SYSTEM", "10M"),
				),
			},
		},
	})
}
func TestAccUser_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic,
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
	testAccUserConfigBasic = `
resource "oraclerdbms_user" "test" {
	username = "USER666"
	default_tablespace = "USERS"
	account_status = "OPEN"
	quota = {
		SYSTEM = "10M"
	}
}
`
	testAccUserConfigBasicUpdate = `
resource "oraclerdbms_user" "test" {
	username = "USER666"
	default_tablespace = "SYSTEM"
	account_status = "LOCKED"
	quota = {
		SYSTEM = "10M"
	}
}
`
)
