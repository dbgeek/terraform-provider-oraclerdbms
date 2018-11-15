package oraclerdbms

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccStatsGlobal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccStatsConfigBasiGlobal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstest", "preference_value", "PARTITION"),
				),
			},
			resource.TestStep{
				Config: testAccStatsConfigBasicGlobalUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstest", "preference_value", "AUTO"),
				),
			},
		},
	})
}

func TestAccStatsTable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccStatsConfigBasicTable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstable", "preference_value", "PARTITION"),
				),
			},
			resource.TestStep{
				Config: testAccStatsConfigBasicTableUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oraclerdbms_stats.statstable", "preference_value", "GLOBAL"),
				),
			},
		},
	})
}

const (
	testAccStatsConfigBasiGlobal = `
resource "oraclerdbms_stats" "statstest" {
	preference_name = "GRANULARITY"
	preference_value = "PARTITION"
}
`
	testAccStatsConfigBasicGlobalUpdate = `
resource "oraclerdbms_stats" "statstest" {
	preference_name = "GRANULARITY"
	preference_value = "AUTO"
}
`
	testAccStatsConfigBasicTable = `
# TEST table in schema need to exists...
resource "oraclerdbms_stats" "statstable" {
	owner_name		= "SYSTEM"
	table_name		= "TEST"
	preference_name = "GRANULARITY"
	preference_value = "PARTITION"
}
`
	testAccStatsConfigBasicTableUpdate = `
# TEST table in schema need to exists...
resource "oraclerdbms_stats" "statstable" {
	owner_name		= "SYSTEM"
	table_name		= "TEST"
	preference_name = "GRANULARITY"
	preference_value = "GLOBAL"
}
`
)
