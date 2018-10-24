package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"

	"log"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateRole,
		Delete: resourceOracleRdbmsDeleteRole,
		Read:   resourceOracleRdbmsReadRole,
		Update: nil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
		},
	}
}

func resourceOracleRdbmsCreateRole(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateRole")

	client := meta.(*providerConfiguration).Client

	resourceRole := oraclehelper.ResourceRole{
		Role: d.Get("role").(string),
	}

	err := client.RoleService.CreateRole(resourceRole)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(d.Get("role").(string))
	return resourceOracleRdbmsReadRole(d, meta)
}

func resourceOracleRdbmsDeleteRole(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteRole")
	client := meta.(*providerConfiguration).Client
	resourceRole := oraclehelper.ResourceRole{
		Role: d.Id(),
	}

	err := client.RoleService.DropRole(resourceRole)
	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}
func resourceOracleRdbmsReadRole(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadRole")
	client := meta.(*providerConfiguration).Client
	resourceRole := oraclehelper.ResourceRole{
		Role: d.Id(),
	}

	role, err := client.RoleService.ReadRole(resourceRole)
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("role", role.Role)
	return nil
}
func resourceOracleRdbmsUpdateRole(d *schema.ResourceData, meta interface{}) error {

	return nil
}
