package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccProfileConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_profile.test", "profile", "TEST01"),
				),
			},
		},
	})
}

func TestAccProfile_importBasic(t *testing.T) {
	resourceName := "oraclerdbms_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccProfileConfigBasic = `
resource "oraclerdbms_profile" "test" {
    profile = "TEST01"
}
`
