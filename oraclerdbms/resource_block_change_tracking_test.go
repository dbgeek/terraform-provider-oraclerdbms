package oraclerdbms

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBlockChangeTracking(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCbt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_block_change_tracking.cbt", "status", "ENABLED"),
				),
			},
		},
	})
}

const testAccCbt = `
resource "oraclerdbms_block_change_tracking" "cbt" {
	status = "ENABLED"
	file_name = "/opt/oracle/product/12.2.0.1/dbhome_1/dbs/bb"
}
`
