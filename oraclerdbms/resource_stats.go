package oraclerdbms

import (
	"fmt"
	"log"
	"strings"

	"github.com/dbgeek/terraform-oracle-rdbms-helper/oraclehelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStats() *schema.Resource {
	return &schema.Resource{
		Create: resourceOracleRdbmsCreateStats,
		Delete: resourceOracleRdbmsDeleteStats,
		Read:   resourceOracleRdbmsReadStats,
		Update: resourceOracleRdbmsUpdateStats,

		Schema: map[string]*schema.Schema{
			"preference_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"owner_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"table_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"preference_value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
		},
	}
}

func resourceOracleRdbmsCreateStats(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsCreateStats")
	client := meta.(*oracleHelperType).Client
	switch {
	case d.Get("owner_name").(string) == "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] global")

		resourceStats := oraclehelper.ResourceStats{
			Pname: d.Get("preference_name").(string),
			Pvalu: d.Get("preference_value").(string),
		}

		err := client.StatsService.SetGlobalPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
		id := fmt.Sprintf("STATS-GLOBAL-%s", d.Get("preference_name").(string))
		d.SetId(id)
	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] schema")
		resourceStats := oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			Pvalu:   d.Get("preference_value").(string),
		}
		err := client.StatsService.SetSchemaPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
		id := fmt.Sprintf("STATS-%s-%s", d.Get("owner_name").(string), d.Get("preference_name").(string))
		d.SetId(id)
	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) != "":
		log.Println("[DEBUG] table")
		resourceStats := oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			TaBName: d.Get("table_name").(string),
			Pvalu:   d.Get("preference_value").(string),
		}
		err := client.StatsService.SetTabPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
		id := fmt.Sprintf("STATS-TABLE-%s-%s", d.Get("table_name").(string), d.Get("preference_name").(string))
		d.SetId(id)
	}
	return nil
}

func resourceOracleRdbmsDeleteStats(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsDeleteStats")
	client := meta.(*oracleHelperType).Client

	switch {
	case d.Get("owner_name").(string) == "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] global")
		/*
			No solution how to reset a parameter global to default :/
			DBMS_STATS.RESET_GLOBAL_PREF_DEFAULTS; resets all parameters to default :/
		*/

	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] schema")
		resourceStats := oraclehelper.ResourceStats{
			Pname: d.Get("preference_name").(string),
		}

		result, err := client.StatsService.ReadGlobalPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
		err = client.StatsService.SetSchemaPre(oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			Pvalu:   result.Pvalu,
		})
		if err != nil {
			d.SetId("")
			return err
		}

	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) != "":
		log.Println("[DEBUG] table")

		resourceStats := oraclehelper.ResourceStats{
			Pname: d.Get("preference_name").(string),
		}

		result, err := client.StatsService.ReadGlobalPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
		err = client.StatsService.SetTabPre(oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			TaBName: d.Get("table_name").(string),
			Pvalu:   result.Pvalu,
		})
		if err != nil {
			d.SetId("")
			return err
		}
	}
	return nil
}

func resourceOracleRdbmsReadStats(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsReadStats")
	client := meta.(*oracleHelperType).Client
	switch {
	case d.Get("owner_name").(string) == "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] global")
		resourceStats := oraclehelper.ResourceStats{
			Pname: d.Get("preference_name").(string),
		}

		result, err := client.StatsService.ReadGlobalPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}

		d.Set("preference_value", result.Pvalu)

	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] schema")
		/*
			Oracle have no support for this
		*/
	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) != "":
		log.Println("[DEBUG] table")
		resourceStats := oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			TaBName: d.Get("table_name").(string),
		}
		result, err := client.StatsService.ReadTabPref(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}

		d.Set("preference_value", result.Pvalu)
	}
	return nil
}

func resourceOracleRdbmsUpdateStats(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceOracleRdbmsUpdateStats")
	client := meta.(*oracleHelperType).Client
	switch {
	case d.Get("owner_name").(string) == "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] global")
		resourceStats := oraclehelper.ResourceStats{
			Pname: d.Get("preference_name").(string),
			Pvalu: d.Get("preference_value").(string),
		}

		err := client.StatsService.SetGlobalPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) == "":
		log.Println("[DEBUG] schema")
		resourceStats := oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			Pvalu:   d.Get("preference_value").(string),
		}
		err := client.StatsService.SetSchemaPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
	case d.Get("owner_name").(string) != "" && d.Get("table_name").(string) != "":
		log.Println("[DEBUG] table")
		resourceStats := oraclehelper.ResourceStats{
			Pname:   d.Get("preference_name").(string),
			OwnName: d.Get("owner_name").(string),
			TaBName: d.Get("table_name").(string),
			Pvalu:   d.Get("preference_value").(string),
		}
		err := client.StatsService.SetTabPre(resourceStats)
		if err != nil {
			d.SetId("")
			return err
		}
	}

	return nil
}
