package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccStats(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccStatsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstest", "preference_value", "PARTITION"),
				),
			},
			resource.TestStep{
				Config: testAccStatsConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstest", "preference_value", "AUTO"),
				),
			},
		},
	})
}

const (
	testAccStatsConfigBasic = `
resource "oraclerdbms_stats" "statstest" {
	preference_name = "GRANULARITY"
	preference_value = "PARTITION"
}
`
	testAccStatsConfigBasicUpdate = `
resource "oraclerdbms_stats" "statstest" {
	preference_name = "GRANULARITY"
	preference_value = "AUTO"
}
`
)
