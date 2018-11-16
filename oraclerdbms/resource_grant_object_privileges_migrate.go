package oraclerdbms

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func resourceOracleRdbmsGrantObjectPrivilegeMigrate(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found oraclerdbms VPC GrantObjectPrivilege v0; migrating to v1")
		return migrateGrantObjectPrivilegeV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateGrantObjectPrivilegeV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty VPC State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)
	is.Attributes["privs_sha256.%"] = "1"
	is.Attributes["privs_sha256.select"] = is.Attributes["privs_sha256"]
	delete(is.Attributes, "privs_sha256")

	is.Attributes["objects_sha256.%"] = "1"
	is.Attributes["objects_sha256.select"] = is.Attributes["objects_sha256"]
	delete(is.Attributes, "objects_sha256")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)

	return is, nil
}
