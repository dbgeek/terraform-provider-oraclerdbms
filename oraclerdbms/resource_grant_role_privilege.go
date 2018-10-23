package oraclerdbms

import (
	"bytes"
	"fmt"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGrantRolePrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantRolePrivilege,
		Delete: resourceOracleRdbmsDeleteGrantRolePrivilege,
		Read:   resourceOracleRdbmsReadGrantRolePrivilege,
		Update: nil,
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
	client := meta.(*providerConfiguration).Client

	resourceGrantRolePrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantRolePrivilege.Role = d.Get("role").(string)

	err := client.GrantService.GrantRolePriv(resourceGrantRolePrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantRolePrivIDHash(d.Get("grantee").(string), d.Get("role").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantRolePrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantRolePrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteGrantRolePrivilege")
	var resourceGrantRolePrivilege oraclehelper.ResourceGrantRolePrivilege
	client := meta.(*providerConfiguration).Client

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
	client := meta.(*providerConfiguration).Client

	resourceGrantRolePrivilege.Grantee = d.Get("grantee").(string)

	rolePrivs, err := client.GrantService.ReadGrantRolePrivs(resourceGrantRolePrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	if _, ok := rolePrivs[d.Get("role").(string)]; !ok {
		d.SetId("")
		return nil
	}

	return nil
}

func grantRolePrivIDHash(grantee string, role string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", grantee))
	buf.WriteString(fmt.Sprintf("%s-", role))
	return fmt.Sprintf("grantrolpriv-%d", hashcode.String(buf.String()))
}
