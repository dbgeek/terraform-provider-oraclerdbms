package oraclehelper

import (
	"fmt"
	"log"
)

const (
	queryUser = `
SELECT
	u.username,
	u.default_tablespace,
	u.temporary_tablespace,
	u.profile
FROM
	dba_users u
WHERE u.username = :1
`
)

type (
	//ResourceUser ..
	ResourceUser struct {
		Username            string
		Password            string
		DefaultTablespace   string
		TemporaryTablespace string
		Profile             string
	}
	//User ..
	User struct {
		Username            string
		Password            string
		DefaultTablespace   string
		TemporaryTablespace string
		Profile             string
	}
	userService struct {
		client *Client
	}
)

func (u *userService) ReadUser(tf ResourceUser) (*User, error) {
	log.Printf("[DEBUG] ReadUser username: %s\n", tf.Username)
	param := &User{}

	err := u.client.DBClient.QueryRow(queryUser, tf.Username).Scan(&param.Username,
		&param.DefaultTablespace,
		&param.TemporaryTablespace,
		&param.Profile,
	)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (u *userService) CreateUser(tf ResourceUser) error {
	log.Println("[DEBUG] CreateUser")
	sqlCommand := fmt.Sprintf("create user %s identified by change_on_install", tf.Username)

	if tf.DefaultTablespace != "" {
		sqlCommand += fmt.Sprintf(" default tablespace %s", tf.DefaultTablespace)
	}
	if tf.TemporaryTablespace != "" {
		sqlCommand += fmt.Sprintf(" temporary tablespace %s", tf.TemporaryTablespace)
	}

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := u.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) ModifyUser(tf ResourceUser) error {
	log.Println("[DEBUG] ModifyUser")
	sqlCommand := fmt.Sprintf("alter user %s", tf.Username)

	if tf.DefaultTablespace != "" {
		sqlCommand += fmt.Sprintf(" default tablespace %s", tf.DefaultTablespace)
	}
	if tf.TemporaryTablespace != "" {
		sqlCommand += fmt.Sprintf(" temporary tablespace %s", tf.TemporaryTablespace)
	}

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := u.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) DropUser(tf ResourceUser) error {
	log.Println("[DEBUG] DeleteUser")
	sqlCommand := fmt.Sprintf("drop user %s", tf.Username)
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := u.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[DEBUG] drop user err: %s\n", err)
		return err
	}

	return nil
}
