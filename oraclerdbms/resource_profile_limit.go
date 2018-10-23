package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"

	"log"
)

func resourceProfileLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateProfileLimit,
		Delete: resourceOracleRdbmsDeleteProfileLimit,
		Read:   resourceOracleRdbmsReadProfileLimit,
		Update: resourceOracleRdbmsUpdateProfileLimit,
		Schema: map[string]*schema.Schema{
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOracleRdbmsCreateProfileLimit(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateProfile")
	var resourceProfile oraclehelper.ResourceProfile
	client := meta.(*providerConfiguration).Client
	resourceProfile.Profile = d.Get("profile").(string)
	resourceProfile.ResourceName = d.Get("resource_name").(string)
	resourceProfile.Limit = d.Get("value").(string)

	client.ProfileService.UpdateProfile(resourceProfile)
	d.SetId(d.Get("resource_name").(string))
	return resourceOracleRdbmsUpdateProfileLimit(d, meta)
}

func resourceOracleRdbmsDeleteProfileLimit(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteProfile")
	var resourceProfile oraclehelper.ResourceProfile
	resourceProfile.Profile = d.Get("profile").(string)
	resourceProfile.ResourceName = d.Get("resource_name").(string)

	client := meta.(*providerConfiguration).Client
	client.ProfileService.ResetProfileResourceLimite(resourceProfile)

	return nil
}

func resourceOracleRdbmsReadProfileLimit(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadProfile")
	var resourceProfile oraclehelper.ResourceProfile
	resourceProfile.Profile = d.Get("profile").(string)
	client := meta.(*providerConfiguration).Client

	profileLimit, _ := client.ProfileService.ReadProfile(resourceProfile)

	d.Set("value", profileLimit[d.Get("resource_name").(string)])

	return nil
}

func resourceOracleRdbmsUpdateProfileLimit(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateProfile")
	var resourceProfile oraclehelper.ResourceProfile
	client := meta.(*providerConfiguration).Client
	resourceProfile.Profile = d.Get("profile").(string)
	resourceProfile.ResourceName = d.Get("resource_name").(string)
	resourceProfile.Limit = d.Get("value").(string)

	client.ProfileService.UpdateProfile(resourceProfile)
	return resourceOracleRdbmsReadProfileLimit(d, meta)
}
