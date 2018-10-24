package oraclerdbms

import (
	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strings"
)

func resourceParameter() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateParameter,
		Delete: resourceOracleRdbmsDeleteParameter,
		Read:   resourceOracleRdbmsReadParameter,
		Update: resourceOracleRdbmsUpdateParameter,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"update_comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: validation.StringInSlice([]string{
					"memory",
					"spfile",
					"both",
				}, true),
				Default: "both",
			},
		},
	}
}

// CreateParameter ..
func resourceOracleRdbmsCreateParameter(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] CreateParameter")
	var resourceParameter oraclehelper.ResourceParameter
	if d.Get("update_comment").(string) != "" {
		resourceParameter.Comment = d.Get("update_comment").(string)
	}
	resourceParameter.Name = d.Get("name").(string)
	resourceParameter.Value = d.Get("value").(string)
	resourceParameter.Scope = d.Get("scope").(string)

	client := meta.(*providerConfiguration).Client

	client.ParameterService.SetParameter(resourceParameter)

	d.SetId(d.Get("name").(string))

	return resourceOracleRdbmsUpdateParameter(d, meta)
}

// DeleteParameter ..
func resourceOracleRdbmsDeleteParameter(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] DeleteParameter")
	var resourceParameter oraclehelper.ResourceParameter
	client := meta.(*providerConfiguration).Client
	resourceParameter.Name = d.Id()

	client.ParameterService.ResetParameter(resourceParameter)
	d.SetId("")

	return nil
}

// ReadParameter ..
func resourceOracleRdbmsReadParameter(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] ReadParameter")
	var resourceParameter oraclehelper.ResourceParameter
	client := meta.(*providerConfiguration).Client
	resourceParameter.Name = d.Id()

	parm, err := client.ParameterService.Read(resourceParameter)
	if err != nil {
		log.Fatal("READERROR")
		d.SetId("")
	}
	if parm.Value == parm.DefaultValue {
		log.Printf("[DEBUG] ReadParameter Is default, value: %s state value: %s\n", parm.Value, d.Get("value").(string))
		d.SetId("")
		return nil
	}
	d.Set("value", parm.Value)
	d.Set("name", parm.Name)

	if d.Get("scope").(string) == "" {
		d.Set("scope", "BOTH")
	}

	log.Printf("[DEBUG] name: %s, value: %s, defaultvalue: %s \n", parm.Name, parm.Value, parm.DefaultValue)

	return nil
}

// UpdateParameter ..@
func resourceOracleRdbmsUpdateParameter(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] UpdateParameter")
	var resourceParameter oraclehelper.ResourceParameter
	if d.Get("update_comment").(string) != "" {
		resourceParameter.Comment = d.Get("update_comment").(string)
	}
	resourceParameter.Name = d.Id()
	resourceParameter.Value = d.Get("value").(string)

	if !d.IsNewResource() {
		client := meta.(*providerConfiguration).Client
		client.ParameterService.SetParameter(resourceParameter)
	}
	return resourceOracleRdbmsReadParameter(d, meta)
}
