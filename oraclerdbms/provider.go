package oraclerdbms

import (
	"fmt"
	"log"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type oracleHelperType struct {
	Client *oraclehelper.Client
}

// Provider ....
func Provider() terraform.ResourceProvider {
	log.Println("[DEBUG] Initializing oraclerdbms ResourceProvider")
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"service": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLE_SERVICE", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("service must not be an empty string"))
					}

					return
				},
			},
			"dbhost": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLE_DBHOST", nil),
			},
			"dbport": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLE_DBPORT", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLE_USERNAME", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Username must not be an empty string"))
					}
					return
				},
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLE_PASSWORD", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Username must not be an empty string"))
					}
					return
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oraclerdbms_autotask":               resourceAutotask(),
			"oraclerdbms_block_change_tracking":  resourceBlockChangeTracking(),
			"oraclerdbms_database":               resourceDatabase(),
			"oraclerdbms_grant_object_privilege": resourceGrantObjectPrivilege(),
			"oraclerdbms_grant_role_privilege":   resourceGrantRolePrivilege(),
			"oraclerdbms_grant_system_privilege": resourceGrantSystemPrivilege(),
			"oraclerdbms_parameter":              resourceParameter(),
			"oraclerdbms_profile":                resourceProfile(),
			"oraclerdbms_profile_limit":          resourceProfileLimit(),
			"oraclerdbms_role":                   resourceRole(),
			"oraclerdbms_stats":                  resourceStats(),
			"oraclerdbms_user":                   resourceUser(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config oraclehelper.Cfg

	if d.Get("username").(string) != "" {
		config.Username = d.Get("username").(string)
	}
	if d.Get("password").(string) != "" {
		config.Password = d.Get("password").(string)
	}
	if d.Get("dbhost").(string) != "" {
		config.DbHost = d.Get("dbhost").(string)
	}
	if d.Get("dbport").(string) != "" {
		config.DbPort = d.Get("dbport").(string)
	}
	if d.Get("service").(string) != "" {
		config.DbService = d.Get("service").(string)
	}

	client := oraclehelper.NewClient(config)

	log.Println("[DEBUG] Initializing Oracle DB Helper client")
	return &oracleHelperType{
		Client: client,
	}, nil
}
