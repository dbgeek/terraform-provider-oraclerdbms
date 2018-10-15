package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateUser,
		Delete: resourceOracleRdbmsDeleteUser,
		Read:   resourceOracleRdbmsReadUser,
		Update: resourceOracleRdbmsUpdateUser,
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"default_tablespace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"temporary_tablespace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateUser")
	var user oraclehelper.ResourceUser
	client := meta.(*providerConfiguration).Client

	if d.Get("username").(string) != "" {
		user.Username = d.Get("username").(string)
	}
	if d.Get("password").(string) != "" {
		user.Password = d.Get("password").(string)
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
	client.UserService.CreateUser(user)

	d.SetId(d.Get("username").(string))

	return resourceOracleRdbmsUpdateUser(d, meta)
}

func resourceOracleRdbmsDeleteUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteUser")
	var user oraclehelper.ResourceUser
	client := meta.(*providerConfiguration).Client

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
	client := meta.(*providerConfiguration).Client

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
	if user.DefaultTablespace != "" {
		d.Set("default_tablespace", user.DefaultTablespace)
	}
	if user.DefaultTablespace != "" {
		d.Set("temporary_tablespace", user.TemporaryTablespace)
	}
	if user.Profile != "" {
		d.Set("profile", user.Profile)
	}

	return nil
}

func resourceOracleRdbmsUpdateUser(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateUser")
	if !d.IsNewResource() {
		var resourceUser oraclehelper.ResourceUser
		client := meta.(*providerConfiguration).Client
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
		client.UserService.ModifyUser(resourceUser)
	}
	return resourceOracleRdbmsReadUser(d, meta)
}
