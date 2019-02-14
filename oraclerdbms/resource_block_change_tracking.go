package oraclerdbms

import (
	"log"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
)

/*
	This version of oraclerdbms_block_change_tracking resource does not support OMF and ASM
*/
func resourceBlockChangeTracking() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateBlockChangeTracking,
		Delete: resourceOracleRdbmsDeleteBlockChangeTracking,
		Read:   resourceOracleRdbmsReadBlockChangeTracking,
		Update: resourceOracleRdbmsUpdateBlockChangeTracking,

		Schema: map[string]*schema.Schema{
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			/*
				How should we handle if we use omf to not generate diff
			*/
			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateBlockChangeTracking(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateBlockChangeTracking")
	client := meta.(*oracleHelperType).Client

	resourceBct := oraclehelper.ResourceBlockChangeTracking{
		FileName: d.Get("file_name").(string),
	}

	err := client.BlockChangeTrackingService.EnableBlockChangeTracking(resourceBct)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId("cbt")

	cbt, err := client.BlockChangeTrackingService.ReadBlockChangeTracking()
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("file_name", cbt.FileName)

	return nil
}

func resourceOracleRdbmsDeleteBlockChangeTracking(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteBlockChangeTracking")
	client := meta.(*oracleHelperType).Client

	err := client.BlockChangeTrackingService.DisableBlockChangeTracking()
	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}
func resourceOracleRdbmsReadBlockChangeTracking(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadBlockChangeTracking")
	client := meta.(*oracleHelperType).Client

	cbt, err := client.BlockChangeTrackingService.ReadBlockChangeTracking()
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("status", cbt.Status)
	d.Set("file_name", cbt.FileName)
	return nil
}
func resourceOracleRdbmsUpdateBlockChangeTracking(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateBlockChangeTracking")
	if d.HasChange("status") {

	}

	return nil
}
