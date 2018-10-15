package oraclehelper

import (
	"fmt"
	"log"
)

const (
	queryRole = `
SELECT
	r.role,
	r.password_required,
	r.authentication_type,
	r.common,
	r.oracle_maintained
FROM
	dba_roles r
WHERE r.role = :1
`
)

//Role ..
type (
	//ResourceRole ....
	ResourceRole struct {
		Role string
	}
	Role struct {
		Role               string
		PasswordRequired   string
		AuthenticationType string
		Common             string
		OracleMaintained   string
	}
	roleService struct {
		client *Client
	}
)

func (r *roleService) ReadRole(tf ResourceRole) (*Role, error) {
	log.Printf("[DEBUG] ReadUser username: %s\n", tf.Role)
	roleType := &Role{}

	err := r.client.DBClient.QueryRow(queryRole, tf.Role).Scan(&roleType.Role,
		&roleType.PasswordRequired,
		&roleType.AuthenticationType,
		&roleType.Common,
		&roleType.OracleMaintained,
	)
	if err != nil {
		return nil, err
	}

	return roleType, nil
}

func (r *roleService) CreateRole(tf ResourceRole) error {
	log.Println("[DEBUG] CreateRole")
	sqlCommand := fmt.Sprintf("create role %s", tf.Role)

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := r.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (r *roleService) ModifyRole(tf ResourceRole) error {
	log.Println("[DEBUG] ModifyRole")
	sqlCommand := fmt.Sprintf("alter user %s", tf.Role)

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := r.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (r *roleService) DropRole(tf ResourceRole) error {
	log.Println("[DEBUG] DropRole")
	sqlCommand := fmt.Sprintf("drop role %s", tf.Role)
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := r.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[DEBUG] drop role err: %s\n", err)
		return err
	}

	return nil
}
