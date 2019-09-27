package oraclerdbms

import (
	"fmt"
	"log"
	"strings"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGrantRolePrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantRolePrivilege,
		Delete: resourceOracleRdbmsDeleteGrantRolePrivilege,
		Read:   resourceOracleRdbmsReadGrantRolePrivilege,
		Update: nil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"grantee": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
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

func resourceOracleRdbmsCreateGrantRolePrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateGrantRolePrivilege")
	var resourceGrantRolePrivilege oraclehelper.ResourceGrantRolePrivilege
	client := meta.(*oracleHelperType).Client

	resourceGrantRolePrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantRolePrivilege.Role = d.Get("role").(string)

	err := client.GrantService.GrantRolePriv(resourceGrantRolePrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantRolePrivID(d.Get("grantee").(string), d.Get("role").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantRolePrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantRolePrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteGrantRolePrivilege")
	var resourceGrantRolePrivilege oraclehelper.ResourceGrantRolePrivilege
	client := meta.(*oracleHelperType).Client

	resourceGrantRolePrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantRolePrivilege.Role = d.Get("role").(string)

	err := client.GrantService.RevokeRolePriv(resourceGrantRolePrivilege)
	if err != nil {
		return err
	}
	return nil
}

func resourceOracleRdbmsReadGrantRolePrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOracleRdbmsReadGrantRolePrivilege grantee:%s\n", d.Get("grantee"))
	var resourceGrantRolePrivilege oraclehelper.ResourceGrantRolePrivilege
	client := meta.(*oracleHelperType).Client

	//ToDo Break out as a function ??
	splitRolerivilege := strings.Split(d.Id(), "-")
	grantee := splitRolerivilege[0]
	role := splitRolerivilege[1]

	resourceGrantRolePrivilege.Grantee = grantee

	rolePrivs, err := client.GrantService.ReadGrantRolePrivs(resourceGrantRolePrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	if _, ok := rolePrivs[role]; !ok {
		d.SetId("")
		return nil
	}
	d.Set("grantee", grantee)
	d.Set("role", role)

	return nil
}

func grantRolePrivID(grantee string, role string) string {
	return fmt.Sprintf("%s-%s", grantee, role)
}
