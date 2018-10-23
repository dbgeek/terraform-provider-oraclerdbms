package oraclerdbms

import (
	"bytes"
	"fmt"
	"log"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGrantSystemPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantSystemPrivilege,
		Delete: resourceOracleRdbmsDeleteGrantSystemPrivilege,
		Read:   resourceOracleRdbmsReadGrantSystemPrivilege,
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
			"privilege": &schema.Schema{
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

func resourceOracleRdbmsCreateGrantSystemPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsCreateGrantSystemPrivilege")
	var resourceGrantSystemPrivilege oraclehelper.ResourceGrantSystemPrivilege
	client := meta.(*providerConfiguration).Client

	resourceGrantSystemPrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantSystemPrivilege.Privilege = d.Get("privilege").(string)
	err := client.GrantService.GrantSysPriv(resourceGrantSystemPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantSysPrivIDHash(d.Get("grantee").(string), d.Get("privilege").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantSystemPrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantSystemPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsDeleteGrantSystemPrivilege")
	var resourceGrantSystemPrivilege oraclehelper.ResourceGrantSystemPrivilege
	client := meta.(*providerConfiguration).Client

	resourceGrantSystemPrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantSystemPrivilege.Privilege = d.Get("privilege").(string)

	err := client.GrantService.RevokeSysPriv(resourceGrantSystemPrivilege)
	if err != nil {
		return err
	}
	return nil
}

func resourceOracleRdbmsReadGrantSystemPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] resourceOracleRdbmsReadGrantSystemPrivilege grantee:%s\n", d.Get("grantee"))
	var resourceGrantSystemPrivilege oraclehelper.ResourceGrantSystemPrivilege
	client := meta.(*providerConfiguration).Client

	resourceGrantSystemPrivilege.Grantee = d.Get("grantee").(string)

	sysPrivs, err := client.GrantService.ReadGrantSysPrivs(resourceGrantSystemPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	if _, ok := sysPrivs[d.Get("privilege").(string)]; !ok {
		d.SetId("")
		return nil
	}

	return nil
}

func grantSysPrivIDHash(grantee string, privilege string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", grantee))
	buf.WriteString(fmt.Sprintf("%s-", privilege))
	return fmt.Sprintf("grantsyspriv-%d", hashcode.String(buf.String()))
}
