package oraclerdbms

import (
	"fmt"
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateDatabase,
		Delete: resourceOracleRdbmsDeleteDatabase,
		Read:   resourceOracleRdbmsReadDatabase,
		Update: resourceOracleRdbmsUpdateDatabase,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"forcelogging": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToUpper(val.(string))
					if !(v == "NO" || v == "YES") {
						errs = append(errs, fmt.Errorf("%q must be YES or NO, got: %s", key, v))
					}
					return
				},
			},
			"flashback": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToUpper(val.(string))
					if !(v == "YES" || v == "NO") {
						errs = append(errs, fmt.Errorf("%q must be YES or NO, got: %s", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceOracleRdbmsCreateDatabase(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateDatabase")
	client := meta.(*oracleHelperType).Client

	resourceDatabase := oraclehelper.ResourceDatabase{
		ForceLogging: d.Get("forcelogging").(string),
		FlashBackOn:  d.Get("flashback").(string),
	}

	err := client.DatabaseService.ModifyDatabase(resourceDatabase)
	if err != nil {
		return err
	}

	database, _ := client.DatabaseService.ReadDatabase()

	d.SetId(database.Name)

	return nil

}

func resourceOracleRdbmsDeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteDatabase")
	return nil

}

func resourceOracleRdbmsReadDatabase(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadDatabase")
	client := meta.(*oracleHelperType).Client

	resourceDatabase, err := client.DatabaseService.ReadDatabase()
	if err != nil {
		return err
	}

	d.Set("flashback", resourceDatabase.FlashBackOn)
	d.Set("forcelogging", resourceDatabase.ForceLogging)

	return nil
}

func resourceOracleRdbmsUpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateDatabase")
	client := meta.(*oracleHelperType).Client

	resourceDatabase := oraclehelper.ResourceDatabase{}

	if d.HasChange("flashback") {
		resourceDatabase.FlashBackOn = d.Get("flashback").(string)
	}

	if d.HasChange("forcelogging") {
		resourceDatabase.ForceLogging = d.Get("forcelogging").(string)
	}

	err := client.DatabaseService.ModifyDatabase(resourceDatabase)
	if err != nil {
		return err
	}
	return nil
}
