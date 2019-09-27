package oraclerdbms

import (
	"fmt"
	"log"
	"strings"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		CustomizeDiff: updateComputed,

		SchemaVersion: 1,
		MigrateState:  resourceOracleRdbmsGrantObjectPrivilegeMigrate,

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
			// This need to be changes to TypeSet. There is mutiple of privileges that need to be checked. Now only support SELECT
			"objects_sha256": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Computed: true,
			},
			// This need to be changes to TypeSet. There is mutiple of privileges that need to be checked. Now only support SELECT
			"privs_sha256": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func updateComputed(d *schema.ResourceDiff, meta interface{}) error {
	log.Println("[INFO] updateComputed")
	if d.Get("object").(string) == "" {
		log.Println("[INFO] updateComputed inside sha256 calc")
		objMap := map[string]string{}
		if v, ok := d.GetOk("objects_sha256"); ok {
			for key, value := range v.(map[string]interface{}) {
				objMap[key] = value.(string)
			}
		}
		privsMap := map[string]string{}
		if v, ok := d.GetOk("privs_sha256"); ok {
			for key, value := range v.(map[string]interface{}) {
				privsMap[key] = value.(string)
			}
		}
		for k, v := range objMap {
			if privsMap[k] != v {
				d.SetNewComputed("objects_sha256")
				d.SetNewComputed("privs_sha256")
				break
			}
		}
	}
	/*if d.Get("privs_sha256").(string) != d.Get("objects_sha256").(string) {
		d.SetNewComputed("objects_sha256")
		d.SetNewComputed("privs_sha256")
	}*/
	//log.Printf("[DEBUG] updateComputed,privs_sha256: %s, objects_sha256: %s\n", d.Get("privs_sha256").(string), d.Get("objects_sha256").(string))

	return nil
}
func resourceOracleRdbmsCreateGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsCreateGrantObjectPrivilege")
	var privilegesList []string
	client := meta.(*oracleHelperType).Client
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
		//var hash map[string]string
		hash := make(map[string]string)
		for _, v := range privilegesList {
			hash[v], err = client.GrantService.GetHashSchemaPrivsToUser(oraclehelper.ResourceGrantObjectPrivilege{
				Grantee:    d.Get("grantee").(string),
				Owner:      d.Get("owner").(string),
				ObjectName: d.Get("object").(string),
				Privilege:  []string{v},
			})
			if err != nil {
				return err
			}
		}
		d.Set("objects_sha256", hash)
		d.Set("privs_sha256", hash)
		id = grantObjectPrivID(d.Get("grantee").(string), d.Get("owner").(string), d.Get("owner").(string)) // hash)
		d.SetId(id)
		return resourceOracleRdbmsReadGrantObjectPrivilege(d, meta)
	}
	return nil
}

func resourceOracleRdbmsDeleteGrantObjectPrivilege(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] resourceOracleRdbmsDeleteGrantObjectPrivilege")
	var privilegesList []string

	client := meta.(*oracleHelperType).Client
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
	var privilegesList []string
	rawPrivileges := d.Get("privilege")
	rawPrivilegesList := rawPrivileges.(*schema.Set).List()
	for _, v := range rawPrivilegesList {
		str := v.(string)
		privilegesList = append(privilegesList, str)
	}
	client := meta.(*oracleHelperType).Client

	splitGrantObjectPrivilege := strings.Split(d.Id(), "-")
	grantee := splitGrantObjectPrivilege[0]
	owner := splitGrantObjectPrivilege[1]
	object := splitGrantObjectPrivilege[2]
	/*
		If object attribute is set the we do the privilege on object level or if owner & grantee not set
		we asume that we trying to import a resource
	*/
	if d.Get("object").(string) != "" || (d.Get("owner").(string) == "" && d.Get("grantee").(string) == "") {
		log.Printf("[INFO] resourceOracleRdbmsReadGrantObjectPrivilege inside object level")
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
		d.Set("grantee", grantee)
		d.Set("owner", owner)
		d.Set("object", object)
		d.Set("privilege", grantedObject.Privileges)
		objHash := make(map[string]string)
		privHash := make(map[string]string)
		for _, v := range grantedObject.Privileges {
			objHash[v] = owner
			privHash[v] = object
		}
		d.Set("objects_sha256", objHash)
		d.Set("privs_sha256", privHash)

		return nil
	}
	/*
		Doing the privilege on schema level. For now we only handle tables on schema level
	*/
	log.Println("[INFO] resourceOracleRdbmsReadGrantObjectPrivilege schema level")
	switch {
	case d.Get("object_type").(string) == "TABLE":
		log.Println("[INFO] resourceOracleRdbmsReadGrantObjectPrivilege schema level on table level")

		privsHash := make(map[string]string)
		tableHash := make(map[string]string)
		for _, v := range privilegesList {
			resourceGrantObjectPrivilege := oraclehelper.ResourceGrantObjectPrivilege{
				Grantee:   grantee,
				Owner:     owner,
				Privilege: []string{v},
			}
			pHash, err := client.GrantService.GetHashSchemaPrivsToUser(resourceGrantObjectPrivilege)
			privsHash[v] = pHash
			log.Printf("[INFO] resourceOracleRdbmsReadGrantObjectPrivilege hash: %s object: %s \n", privsHash, object)
			if err != nil {
				return err
			}
			tableHash[v], err = client.GrantService.GetHashSchemaAllTables(resourceGrantObjectPrivilege)
			if err != nil {
				return err
			}
			if !d.IsNewResource() {
				d.Set("objects_sha256", tableHash)
				d.Set("privs_sha256", privsHash)
			}
		}
		return nil
	}
	return nil
}

func grantObjectPrivID(grantee string, owner string, object string) string {

	return fmt.Sprintf("%s-%s-%s", grantee, owner, object)

}
