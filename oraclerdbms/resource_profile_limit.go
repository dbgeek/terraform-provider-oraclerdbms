package oraclerdbms

import (
	"fmt"
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"

	"log"
)

func resourceProfileLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateProfileLimit,
		Delete: resourceOracleRdbmsDeleteProfileLimit,
		Read:   resourceOracleRdbmsReadProfileLimit,
		Update: resourceOracleRdbmsUpdateProfileLimit,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	d.SetId(getProfileLimitID(d.Get("profile").(string), d.Get("resource_name").(string)))
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

	//ToDo Break out as a function ??
	splitProfileLimit := strings.Split(d.Id(), "-")
	profile := splitProfileLimit[0]
	resourceName := splitProfileLimit[1]
	resourceProfile := oraclehelper.ResourceProfile{
		Profile: profile,
	}

	client := meta.(*providerConfiguration).Client

	profileLimit, _ := client.ProfileService.ReadProfile(resourceProfile)

	d.Set("value", profileLimit[resourceName])
	d.Set("profile", profileLimit["PROFILE"])
	d.Set("resource_name", resourceName)

	return nil
}

func resourceOracleRdbmsUpdateProfileLimit(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateProfile")
	client := meta.(*providerConfiguration).Client

	resourceProfile := oraclehelper.ResourceProfile{
		Profile:      d.Get("profile").(string),
		ResourceName: d.Get("resource_name").(string),
		Limit:        d.Get("value").(string),
	}
	client.ProfileService.UpdateProfile(resourceProfile)

	return resourceOracleRdbmsReadProfileLimit(d, meta)
}

func getProfileLimitID(profile string, resourceName string) string {
	return fmt.Sprintf("%s-%s", profile, resourceName)
}
