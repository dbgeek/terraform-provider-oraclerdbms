package oraclerdbms

import (
	"log"
	"strings"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateUser,
		Delete: resourceOracleRdbmsDeleteUser,
		Read:   resourceOracleRdbmsReadUser,
		Update: resourceOracleRdbmsUpdateUser,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceOracleRdbmsUserMigrate,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"account_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OPEN",
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"default_tablespace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"temporary_tablespace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"quota": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateUser")
	var user oraclehelper.ResourceUser
	client := meta.(*oracleHelperType).Client

	if d.Get("username").(string) != "" {
		user.Username = d.Get("username").(string)
	}
	if d.Get("account_status").(string) != "" {
		v := d.Get("account_status").(string)
		switch {
		case v == "OPEN":
			user.AccountStatus = "UNLOCK"
		case v == "LOCKED":
			user.AccountStatus = "LOCK"
		case strings.HasPrefix(v, "EXPIRED"):
			user.AccountStatus = "EXPIRED"
		}
	}
	if d.Get("profile").(string) != "" {
		user.Profile = d.Get("profile").(string)
	}
	if d.Get("default_tablespace").(string) != "" {
		user.DefaultTablespace = d.Get("default_tablespace").(string)
	}
	if d.Get("temporary_tablespace").(string) != "" {
		user.TemporaryTablespace = d.Get("temporary_tablespace").(string)
	}
	quotaMap := map[string]string{}
	if v, ok := d.GetOk("quota"); ok {
		for key, value := range v.(map[string]interface{}) {
			quotaMap[key] = value.(string)
		}
		user.Quota = quotaMap
	}
	client.UserService.CreateUser(user)

	d.SetId(d.Get("username").(string))

	return resourceOracleRdbmsUpdateUser(d, meta)
}

func resourceOracleRdbmsDeleteUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteUser")
	var user oraclehelper.ResourceUser
	client := meta.(*oracleHelperType).Client

	user.Username = d.Id()
	err := client.UserService.DropUser(user)
	if err != nil {
		log.Fatalf("Error droping user err msg: %v", err)
	}

	return nil
}

func resourceOracleRdbmsReadUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadUser")
	var resourceUser oraclehelper.ResourceUser
	client := meta.(*oracleHelperType).Client

	resourceUser.Username = d.Id()
	user, err := client.UserService.ReadUser(resourceUser)
	log.Printf("[DEBUG] Resource read user user: %v\n", user)
	if err != nil {
		log.Println("exit error not nil")
		d.SetId("")
		return nil
	}
	if user.Username == "" {
		log.Println("exit username nil")
		d.SetId("")
		return nil
	}

	if user == nil {
		log.Println("exit user nil")
		d.SetId("")
		return nil
	}
	if user.Username != "" {
		d.Set("username", user.Username)
	}
	if user.AccountStatus != "" {
		switch {
		case user.AccountStatus == "OPEN":
			d.Set("account_status", user.AccountStatus)
		case user.AccountStatus == "LOCKED":
			d.Set("account_status", user.AccountStatus)
		case strings.HasPrefix(user.AccountStatus, "EXPIRED"):
			d.Set("account_status", "EXPIRED")
		}
	}
	if user.DefaultTablespace != "" {
		d.Set("default_tablespace", user.DefaultTablespace)
	}
	if user.DefaultTablespace != "" {
		d.Set("temporary_tablespace", user.TemporaryTablespace)
	}
	if user.Profile != "" {
		d.Set("profile", user.Profile)
	}
	if len(user.Quota) > 0 {
		d.Set("quota", user.Quota)
	}

	return nil
}

func resourceOracleRdbmsUpdateUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateUser")
	if !d.IsNewResource() {
		var resourceUser oraclehelper.ResourceUser
		client := meta.(*oracleHelperType).Client
		if d.Get("username").(string) != "" {
			resourceUser.Username = d.Get("username").(string)
		}
		if d.Get("profile").(string) != "" {
			resourceUser.Profile = d.Get("profile").(string)
		}
		if d.Get("default_tablespace").(string) != "" {
			resourceUser.DefaultTablespace = d.Get("default_tablespace").(string)
		}
		if d.Get("temporary_tablespace").(string) != "" {
			resourceUser.TemporaryTablespace = d.Get("temporary_tablespace").(string)
		}
		if d.HasChange("account_status") {
			v := d.Get("account_status").(string)
			switch {
			case v == "OPEN":
				resourceUser.AccountStatus = "UNLOCK"
			case v == "LOCKED":
				resourceUser.AccountStatus = "LOCK"
			case strings.HasPrefix(v, "EXPIRED"):
				resourceUser.AccountStatus = "EXPIRED"
			}
		}
		quotaMap := map[string]string{}
		if v, ok := d.GetOk("quota"); ok {
			for key, value := range v.(map[string]interface{}) {
				quotaMap[key] = value.(string)
			}
			resourceUser.Quota = quotaMap
		}
		client.UserService.ModifyUser(resourceUser)
	}
	return resourceOracleRdbmsReadUser(d, meta)
}
