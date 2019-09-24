package oraclerdbms

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var (
	autoTasks = [3]string{"sql tuning advisor", "auto optimizer stats collection", "auto space advisor"}
)

func TestAccAutotask(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGAutotaskConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_autotask.autotask", "status", "true"),
				),
			},
		},
	})
}

const testAccGAutotaskConfigBasic = `
resource "oraclerdbms_autotask" "autotask" {
	client_name = "sql tuning advisor"
	status = true
}
`
