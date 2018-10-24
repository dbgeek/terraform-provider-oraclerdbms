package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccProfileLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccProfileLimitConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_profile_limit.test1", "value", "33"),
				),
			},
		},
	})
}

func TestAccProfileLimit_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_profile_limit.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileLimitConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccProfileLimitConfigBasic = `
resource "oraclerdbms_profile" "test1" {
    profile = "TEST666"
}

resource "oraclerdbms_profile_limit" "test1" {
	resource_name = "IDLE_TIME"
	value = "33"
	profile = "${oraclerdbms_profile.test1.id}"
}
`
