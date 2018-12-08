package oraclerdbms

import (
	"fmt"
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	//	"strings"
)

func resourceAutotask() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateAutotask,
		Delete: resourceOracleRdbmsDeleteAutotask,
		Read:   resourceOracleRdbmsReadAutotask,
		Update: resourceOracleRdbmsUpdateAutotask,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"client_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateAutotask(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateAutotask")
	client := meta.(*oracleHelperType).Client
	var err error
	if d.Get("status").(bool) {
		err = client.AutoTaskService.EnableAutoTask(oraclehelper.ResourceAutoTask{
			ClientName: d.Get("client_name").(string),
			Status:     "YES"},
		)
	} else {
		err = client.AutoTaskService.EnableAutoTask(oraclehelper.ResourceAutoTask{
			ClientName: d.Get("client_name").(string),
			Status:     "NO"},
		)
	}

	d.SetId(d.Get("client_name").(string))

	return err
}
func resourceOracleRdbmsDeleteAutotask(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteAutotask")
	return nil
}

func resourceOracleRdbmsReadAutotask(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadAutotask")
	client := meta.(*oracleHelperType).Client

	resourceAutoTask, err := client.AutoTaskService.ReadAutoTask(oraclehelper.ResourceAutoTask{
		ClientName: d.Get("client_name").(string),
	})
	fmt.Print("INFO")
	if err != nil {
		return err
	}

	d.Set("status", resourceAutoTask.Status)

	return nil
}

func resourceOracleRdbmsUpdateAutotask(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateAutotask")
	client := meta.(*oracleHelperType).Client
	var err error
	if d.Get("status").(bool) {
		err = client.AutoTaskService.EnableAutoTask(oraclehelper.ResourceAutoTask{
			ClientName: d.Get("client_name").(string),
		})
	} else {
		err = client.AutoTaskService.DisableAutoTask(oraclehelper.ResourceAutoTask{
			ClientName: d.Get("client_name").(string),
		})
	}
	return err
}
