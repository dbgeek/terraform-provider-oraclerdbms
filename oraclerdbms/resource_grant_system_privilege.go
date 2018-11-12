package oraclerdbms

import (
	"fmt"
	"log"
	"strings"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGrantSystemPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantSystemPrivilege,
		Delete: resourceOracleRdbmsDeleteGrantSystemPrivilege,
		Read:   resourceOracleRdbmsReadGrantSystemPrivilege,
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
	client := meta.(*oracleHelperType).Client

	resourceGrantSystemPrivilege.Grantee = d.Get("grantee").(string)
	resourceGrantSystemPrivilege.Privilege = d.Get("privilege").(string)
	err := client.GrantService.GrantSysPriv(resourceGrantSystemPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantSysPrivID(d.Get("grantee").(string), d.Get("privilege").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantSystemPrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantSystemPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsDeleteGrantSystemPrivilege")
	var resourceGrantSystemPrivilege oraclehelper.ResourceGrantSystemPrivilege
	client := meta.(*oracleHelperType).Client

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
	client := meta.(*oracleHelperType).Client

	//ToDo Break out as a function ??
	splitSystemPrivilege := strings.Split(d.Id(), "-")
	grantee := splitSystemPrivilege[0]
	privilege := splitSystemPrivilege[1]

	resourceGrantSystemPrivilege.Grantee = grantee

	sysPrivs, err := client.GrantService.ReadGrantSysPrivs(resourceGrantSystemPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	if _, ok := sysPrivs[privilege]; !ok {
		d.SetId("")
		return nil
	}

	d.Set("grantee", grantee)
	d.Set("privilege", privilege)

	return nil
}

func grantSysPrivID(grantee string, privilege string) string {
	return fmt.Sprintf("%s-%s", grantee, privilege)
}
