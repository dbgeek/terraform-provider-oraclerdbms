package oraclerdbms

import (
	"fmt"
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"

	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceGrantObjectPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateGrantObjectPrivilege,
		Delete: resourceOracleRdbmsDeleteGrantObjectPrivilege,
		Read:   resourceOracleRdbmsReadGrantObjectPrivilege,
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
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"privilege": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
		},
	}
}

func resourceOracleRdbmsCreateGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateGrantObjectPrivilege")
	var privilegesList []string
	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}

	resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
		Grantee:    d.Get("grantee").(string),
		Owner:      d.Get("owner").(string),
		ObjectName: d.Get("object").(string),
		Privilege:  privilegesList,
	}

	err := client.GrantService.GrantObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	id := grantObjectPrivID(d.Get("grantee").(string), d.Get("owner").(string), d.Get("object").(string))
	d.SetId(id)
	return resourceOracleRdbmsReadGrantObjectPrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteGrantObjectPrivilege")
	var privilegesList []string

	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}

	resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
		Grantee:    d.Get("grantee").(string),
		Owner:      d.Get("owner").(string),
		ObjectName: d.Get("object").(string),
		Privilege:  privilegesList,
	}

	err := client.GrantService.RevokeObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		return err
	}
	return nil
}

func resourceOracleRdbmsReadGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOracleRdbmsReadGrantObjectPrivilege grantee:%s\n", d.Get("grantee"))
	client := meta.(*providerConfiguration).Client

	splitGrantObjectPrivilege := strings.Split(d.Id(), "-")
	grantee := splitGrantObjectPrivilege[0]
	owner := splitGrantObjectPrivilege[1]
	object := splitGrantObjectPrivilege[2]

	resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
		Grantee:    grantee,
		Owner:      owner,
		ObjectName: object,
	}
	grantedObject, err := client.GrantService.ReadGrantObjectPrivilege(resourceGrantObjectPrivilege)
	if err != nil {
		d.SetId("")
		return err
	}
	if !d.IsNewResource() {
		d.Set("grantee", grantee)
		d.Set("owner", owner)
		d.Set("object", object)
		d.Set("privilege", grantedObject.Privileges)
	}
	return nil
}

func grantObjectPrivID(grantee string, owner string, object string) string {
	return fmt.Sprintf("%s-%s-%s", grantee, owner, object)
}
