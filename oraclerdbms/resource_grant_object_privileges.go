package oraclerdbms

import (
	"bytes"
	"fmt"
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGrantObjectPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantObjectPrivilege,
		Delete: resourceOracleRdbmsDeleteGrantObjectPrivilege,
		Read:   resourceOracleRdbmsReadGrantObjectPrivilege,
		Update: nil,
		Schema: map[string]*schema.Schema{
			"grantee": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privilege": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateGrantObjectPrivilege")
	var privilegesList []string
	var resourceGrantObjectPrivilege oraclehelper.ResourceGrantObjectPrivilege
	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}
	resourceGrantObjectPrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantObjectPrivilege.Privilege = privilegesList
	resourceGrantObjectPrivilege.Owner = d.Get("owner").(string)
	resourceGrantObjectPrivilege.ObjectName = d.Get("object").(string)

	err := client.GrantService.GrantObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantObjectPrivIDHash(d.Get("grantee").(string), d.Get("owner").(string), d.Get("object").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantObjectPrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteGrantObjectPrivilege")
	var privilegesList []string
	var resourceGrantObjectPrivilege oraclehelper.ResourceGrantObjectPrivilege
	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}

	resourceGrantObjectPrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantObjectPrivilege.Privilege = privilegesList
	resourceGrantObjectPrivilege.Owner = d.Get("owner").(string)
	resourceGrantObjectPrivilege.ObjectName = d.Get("object").(string)

	err := client.GrantService.RevokeObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		return err
	}
	return nil
}

func resourceOracleRdbmsReadGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOracleRdbmsReadGrantObjectPrivilege grantee:%s\n", d.Get("grantee"))
	var resourceGrantObjectPrivilege oraclehelper.ResourceGrantObjectPrivilege
	client := meta.(*providerConfiguration).Client

	resourceGrantObjectPrivilege.Grantee = d.Get("grantee").(string)

	_, err := client.GrantService.ReadGrantObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func grantObjectPrivIDHash(grantee string, owner string, object string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", grantee))
	buf.WriteString(fmt.Sprintf("%s-", owner))
	buf.WriteString(fmt.Sprintf("%s-", object))
	return fmt.Sprintf("grantobjpriv-%d", hashcode.String(buf.String()))
}
