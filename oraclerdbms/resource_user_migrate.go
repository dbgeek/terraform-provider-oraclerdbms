package oraclerdbms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceOracleRdbmsUserMigrate(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found oraclerdbms user v0; migrating to v1")
		return migrateUserV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateUserV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	delete(is.Attributes, "password")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)

	return is, nil
}
