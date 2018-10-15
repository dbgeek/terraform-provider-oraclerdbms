package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"

	"log"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateRole,
		Delete: resourceOracleRdbmsDeleteRole,
		Read:   resourceOracleRdbmsReadRole,
		Update: nil, //resourceOracleRdbmsUpdateRole,
		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateRole(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateRole")
	var resourceRole oraclehelper.ResourceRole
	client := meta.(*providerConfiguration).Client

	if d.Get("role").(string) != "" {
		resourceRole.Role = d.Get("role").(string)
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
	var resourceRole oraclehelper.ResourceRole
	client := meta.(*providerConfiguration).Client

	resourceRole.Role = d.Id()
	err := client.RoleService.DropRole(resourceRole)
	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}
func resourceOracleRdbmsReadRole(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadRole")
	var resourceRole oraclehelper.ResourceRole
	client := meta.(*providerConfiguration).Client

	resourceRole.Role = d.Id()
	_, err := client.RoleService.ReadRole(resourceRole)
	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}
func resourceOracleRdbmsUpdateRole(d *schema.ResourceData, meta interface{}) error {

	return nil
}
