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
				Optional: true,
				Computed: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) == strings.ToLower(new)
				},
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"object_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
	log.Println("[INFO] resourceOracleRdbmsCreateGrantObjectPrivilege")
	var privilegesList []string
	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}
	/*
		If object attribute is set the we do the privilege on object level
	*/
	if d.Get("object").(string) != "" {
		log.Println("[INFO] resourceOracleRdbmsCreateGrantObjectPrivilege Object level flow")
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
	/*
		initial we only handle object_type = TABLE for privileges on schema level.
	*/
	switch {
	case d.Get("object_type").(string) == "TABLE":
		resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
			Grantee:    d.Get("grantee").(string),
			Owner:      d.Get("owner").(string),
			ObjectName: d.Get("object").(string),
			Privilege:  privilegesList,
		}
		err := client.GrantService.GrantTableSchemaToUser(resourceGrantObjectPrivilege)
		if err != nil {
			d.SetId("")
			return err
		}
		id := grantObjectPrivID(d.Get("grantee").(string), d.Get("owner").(string), d.Get("object").(string))
		var hash string
		hash, err = client.GrantService.GetHashSchemaPrivsToUser(resourceGrantObjectPrivilege)
		id = grantObjectPrivID(d.Get("grantee").(string), d.Get("owner").(string), hash)
		d.SetId(id)
		return resourceOracleRdbmsReadGrantObjectPrivilege(d, meta)
	}
	return resourceOracleRdbmsReadGrantObjectPrivilege(d, meta)
}

func resourceOracleRdbmsDeleteGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsDeleteGrantObjectPrivilege")
	var privilegesList []string

	client := meta.(*providerConfiguration).Client
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}
	if d.Get("object").(string) != "" {
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
	resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
		Grantee:   d.Get("grantee").(string),
		Owner:     d.Get("owner").(string),
		Privilege: privilegesList,
	}
	err := client.GrantService.RevokeTableSchemaFromUser(resourceGrantObjectPrivilege)
	if err != nil {
		return err
	}
	return nil

}

func resourceOracleRdbmsReadGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] resourceOracleRdbmsReadGrantObjectPrivilege grantee:%s\n", d.Get("grantee"))
	client := meta.(*providerConfiguration).Client

	splitGrantObjectPrivilege := strings.Split(d.Id(), "-")
	grantee := splitGrantObjectPrivilege[0]
	owner := splitGrantObjectPrivilege[1]
	object := splitGrantObjectPrivilege[2]
	/*
		If object attribute is set the we do the privilege on object level or if owner & grantee not set
		we asume that we trying to import a resource
	*/
	if d.Get("object").(string) != "" || (d.Get("owner").(string) == "" && d.Get("grantee").(string) == "") {
		resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
			Grantee:    grantee,
			Owner:      owner,
			ObjectName: object,
		}
		grantedObject, err := client.GrantService.ReadGrantObjectPrivilege(resourceGrantObjectPrivilege)
		if err != nil {
			log.Printf("[INFO] ReadGrantObjectPrivilege failed with error: %v", err)
			d.SetId("")
			return err
		}
		log.Printf("[INFO] ReadGrantObjectPrivilege returned: %v\n", grantedObject)
		if !d.IsNewResource() {
			d.Set("grantee", grantee)
			d.Set("owner", owner)
			d.Set("object", object)
			d.Set("privilege", grantedObject.Privileges)
		}
		return nil
	}
	/*
		Doing the privilege on schema level. For now we only handle tables on schema level
	*/
	log.Println("[DEBUG] resourceOracleRdbmsReadGrantObjectPrivilege schema level")
	switch {
	case d.Get("object_type").(string) == "TABLE":
		log.Println("[DEBUG] resourceOracleRdbmsReadGrantObjectPrivilege schema level on table level")
		resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
			Grantee: grantee,
			Owner:   owner,
		}
		hash, err := client.GrantService.GetHashSchemaAllTables(resourceGrantObjectPrivilege)
		log.Printf("[DEBUG] resourceOracleRdbmsReadGrantObjectPrivilege hash: %s object: %s \n", hash, object)
		if err != nil {
			return err
		}
		if hash != object {
			log.Printf("[WARN] Hash diff between in state (%s) and the new calculated hash (%s), removing from state", object, hash)
			d.SetId("")
		}
		return nil
	}
	return nil
}

func grantObjectPrivID(grantee string, owner string, object string) string {

	return fmt.Sprintf("%s-%s-%s", grantee, owner, object)

}
