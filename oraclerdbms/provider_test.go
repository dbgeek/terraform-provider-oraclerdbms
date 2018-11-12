package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var client oraclehelper.Client

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProvider.ConfigureFunc = providerConfigure
	testAccProviders = map[string]terraform.ResourceProvider{
		"oraclerdbms": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"ORACLE_USERNAME", "ORACLE_PASSWORD", "ORACLE_DBHOST", "ORACLE_DBPORT", "ORACLE_SERVICE"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("ORACLE_USERNAME, ORACLE_PASSWORD,ORACLE_PASSWORD, ORACLE_DBHOST,ORACLE_DBPORT, ORACLE_SERVICE must be set for acceptance tests")
		}
	}
}
