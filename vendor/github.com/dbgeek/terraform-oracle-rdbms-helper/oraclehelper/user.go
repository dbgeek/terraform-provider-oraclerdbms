package oraclehelper

import (
	"fmt"
	"github.com/mattrobenolt/size"
	"log"
)

const (
	queryUser = `
SELECT
	u.username,
	u.default_tablespace,
	u.temporary_tablespace,
	u.profile,
	u.account_status
FROM
	dba_users u
WHERE u.username = UPPER(:1)
`
	queryQuota = `
SELECT
	q.tablespace_name,
	TO_CHAR(q.max_bytes) AS max_bytes
FROM DBA_TS_QUOTAS q
WHERE q.username = UPPER(:1)
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
		AccountStatus       string
		Quota               map[string]string
	}
	//User ..
	User struct {
		Username            string
		Password            string
		DefaultTablespace   string
		TemporaryTablespace string
		Profile             string
		AccountStatus       string
		Quota               map[string]string
	}
	userService struct {
		client *Client
	}
)

func (u *userService) ReadUser(tf ResourceUser) (*User, error) {
	log.Printf("[DEBUG] ReadUser username: %s\n", tf.Username)
	quota := make(map[string]string)
	param := &User{}

	err := u.client.DBClient.QueryRow(queryUser, tf.Username).Scan(&param.Username,
		&param.DefaultTablespace,
		&param.TemporaryTablespace,
		&param.Profile,
		&param.AccountStatus,
	)
	if err != nil {
		return nil, err
	}

	rows, err := u.client.DBClient.Query(queryQuota, tf.Username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rowTablespace string
		var rowBytes string
		if err := rows.Scan(&rowTablespace, &rowBytes); err != nil {
			return nil, err
		}
		if rowBytes == "-1" {
			quota[rowTablespace] = "unlimited"
		} else {
			s, err := size.ParseCapacity(rowBytes)
			if err != nil {
				log.Printf("parse size")
			}
			quota[rowTablespace] = s.String()
		}

	}

	param.Quota = quota

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
	if tf.AccountStatus != "" {
		sqlCommand += fmt.Sprintf(" account %s", tf.AccountStatus)
	}
	if tf.Profile != "" {
		sqlCommand += fmt.Sprintf(" profile %s", tf.Profile)
	}
	if tf.Quota != nil {
		for k, v := range tf.Quota {
			sqlCommand += fmt.Sprintf(" quota %s on %s", v, k)
		}
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
	if tf.AccountStatus != "" {
		sqlCommand += fmt.Sprintf(" account %s", tf.AccountStatus)
	}
	if tf.Profile != "" {
		sqlCommand += fmt.Sprintf(" profile %s", tf.Profile)
	}
	if tf.Quota != nil {
		for k, v := range tf.Quota {
			sqlCommand += fmt.Sprintf(" quota %s on %s", v, k)
		}
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
