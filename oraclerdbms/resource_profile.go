package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateProfile,
		Delete: resourceOracleRdbmsDeleteProfile,
		Read:   resourceOracleRdbmsReadProfile,
		Update: resourceOracleRdbmsUpdateProfile,
		Schema: map[string]*schema.Schema{
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
		},
	}
}

func resourceOracleRdbmsCreateProfile(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateProfile")
	var resourceProfile oraclehelper.ResourceProfile
	client := meta.(*providerConfiguration).Client
	resourceProfile.Profile = d.Get("profile").(string)
	client.ProfileService.CreateProfile(resourceProfile)

	d.SetId(d.Get("profile").(string))

	return resourceOracleRdbmsUpdateProfile(d, meta)
}

func resourceOracleRdbmsDeleteProfile(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteProfile")
	var resourceProfile oraclehelper.ResourceProfile
	client := meta.(*providerConfiguration).Client

	resourceProfile.Profile = d.Id()
	client.ProfileService.DeleteProfile(resourceProfile)
	return nil
}

func resourceOracleRdbmsReadProfile(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadProfile")
	return nil
}

func resourceOracleRdbmsUpdateProfile(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateProfile")

	return resourceOracleRdbmsReadProfile(d, meta)
}
