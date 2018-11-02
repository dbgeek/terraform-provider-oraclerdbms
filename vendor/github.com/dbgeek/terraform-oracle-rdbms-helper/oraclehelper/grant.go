package oraclehelper

import (
	"fmt"
	"log"
	"strings"
)

const (
	queryTabPrivs = `
SELECT
	*
FROM
	dba_tab_privs tp
WHERE tp.grantee = UPPER(:1)
AND tp.owner = UPPER(:2)
AND tp.table_name = UPPER(:3)
`
	querySysPrivs = `
SELECT
	*
FROM
	dba_sys_privs sp
WHERE sp.grantee = UPPER(:1)
`
	queryRolePrivs = `
SELECT
	*
FROM
	dba_role_privs rp
WHERE rp.grantee = UPPER(:1)
`
)

type (
	//ResourceGrantObjectPrivilege ...
	ResourceGrantObjectPrivilege struct {
		Grantee    string
		Privilege  []string
		Owner      string
		ObjectName string
	}
	//ResourceGrantSystemPrivilege ..
	ResourceGrantSystemPrivilege struct {
		Grantee   string
		Privilege string
	}
	//ResourceGrantRolePrivilege ..
	ResourceGrantRolePrivilege struct {
		Grantee string
		Role    string
	}
	//GrantTable ...
	GrantTable struct {
		Grantee   string
		Owner     string
		TableName string
		Grantor   string
		Privilege string
		Grantable string
		Hierarchy string
		Common    string
		Type      string
		Inherited string
	}
	//GrantObjectPrivs ...
	GrantObjectPrivs struct {
		Grantee    string
		Owner      string
		ObjectName string
		Privileges []string
	}
	//GrantSysPrivs ..
	GrantSysPrivs struct {
		Grantee     string
		Privilege   string
		AdminOption string
		Common      string
		Inherited   string
	}
	//GrantRolePrivs ...
	GrantRolePrivs struct {
		Grantee        string
		GrantedRole    string
		AdminOption    string
		DelegateOption string
		DefaultRole    string
		Common         string
		Inherited      string
	}
	grantService struct {
		client *Client
	}
)

func (tp *grantService) ReadGrantObjectPrivilege(tf ResourceGrantObjectPrivilege) (GrantObjectPrivs, error) {
	log.Printf("[DEBUG] ReadGrantTab grantee: %s\n", tf.Grantee)
	var privileges []string
	rows, err := tp.client.DBClient.Query(queryTabPrivs, tf.Grantee, tf.Owner, tf.ObjectName)
	if err != nil {
		return GrantObjectPrivs{}, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return GrantObjectPrivs{}, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		var grantTable GrantTable
		log.Printf("[DEBUG] getting privs: ")
		if val, ok := m["GRANTEE"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["OWNER"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["TABLE_NAME"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["GRANTOR"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["PRIVILEGE"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["GRANTABLE"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["HIERARCHY"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["COMMON"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["TYPE"].(string); ok {
			grantTable.Grantee = val
		}
		if val, ok := m["INHERITED"].(string); ok {
			grantTable.Grantee = val
		}
		log.Printf("[DEBUG] getting privs: grantTable.Privilege")
		privileges = append(privileges, grantTable.Privilege)
	}

	return GrantObjectPrivs{Grantee: tf.Grantee,
		Owner:      tf.Owner,
		ObjectName: tf.ObjectName,
		Privileges: privileges}, nil
}

func (tp *grantService) ReadGrantSysPrivs(tf ResourceGrantSystemPrivilege) (map[string]GrantSysPrivs, error) {
	log.Printf("[DEBUG] ReadGrantSysPrivs grantee: %s\n", tf.Grantee)
	sysPrivs := make(map[string]GrantSysPrivs)
	rows, err := tp.client.DBClient.Query(querySysPrivs, tf.Grantee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		var sysPriv GrantSysPrivs
		if val, ok := m["GRANTEE"].(string); ok {
			sysPriv.Grantee = val
		}
		if val, ok := m["PRIVILEGE"].(string); ok {
			sysPriv.Privilege = val
		}
		if val, ok := m["ADMIN_OPTION"].(string); ok {
			sysPriv.AdminOption = val
		}
		if val, ok := m["COMMON"].(string); ok {
			sysPriv.Common = val
		}
		if val, ok := m["INHERITED"].(string); ok {
			sysPriv.Inherited = val
		}

		sysPrivs[sysPriv.Privilege] = sysPriv
	}

	return sysPrivs, nil
}

func (tp *grantService) ReadGrantRolePrivs(tf ResourceGrantRolePrivilege) (map[string]GrantRolePrivs, error) {
	log.Printf("[DEBUG] ReadGrantRolePrivs grantee: %s\n", tf.Grantee)
	rolePrivs := make(map[string]GrantRolePrivs)
	rows, err := tp.client.DBClient.Query(queryRolePrivs, tf.Grantee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		var rolePriv GrantRolePrivs
		if val, ok := m["GRANTEE"].(string); ok {
			rolePriv.Grantee = val
		}
		if val, ok := m["GRANTED_ROLE"].(string); ok {
			rolePriv.GrantedRole = val
		}
		if val, ok := m["ADMIN_OPTION"].(string); ok {
			rolePriv.AdminOption = val
		}
		if val, ok := m["DELEGATE_OPTION"].(string); ok {
			rolePriv.DelegateOption = val
		}
		if val, ok := m["DEFAULT_ROLE"].(string); ok {
			rolePriv.DefaultRole = val
		}
		if val, ok := m["COMMON"].(string); ok {
			rolePriv.Common = val
		}
		if val, ok := m["INHERITED"].(string); ok {
			rolePriv.Inherited = val
		}

		rolePrivs[rolePriv.GrantedRole] = rolePriv
	}

	return rolePrivs, nil
}

func (tp *grantService) RevokeObjectPrivilege(tf ResourceGrantObjectPrivilege) error {
	log.Printf("[DEBUG] RevokeGranTab grantee: %s privs: %s\n", tf.Grantee, strings.Join(tf.Privilege, ","))
	sqlCommand := fmt.Sprintf("revoke %s on %s.%s from %s", strings.Join(tf.Privilege, ","), tf.Owner, tf.ObjectName, tf.Grantee)
	log.Printf("[DEBUG] RevokeGranTab revoke sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (tp *grantService) GrantObjectPrivilege(tf ResourceGrantObjectPrivilege) error {
	log.Printf("[DEBUG] GrantGranTab grantee: %s\n", tf.Grantee)
	sqlCommand := fmt.Sprintf("grant %s on %s.%s to %s", strings.Join(tf.Privilege, ","), tf.Owner, tf.ObjectName, tf.Grantee)
	log.Printf("[DEBUG] sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}
func (tp *grantService) GrantRolePriv(tf ResourceGrantRolePrivilege) error {
	log.Printf("[DEBUG] GrantRolePriv grantee: %s\n", tf.Grantee)
	sqlCommand := fmt.Sprintf("grant %s to %s", tf.Role, tf.Grantee)
	log.Printf("[DEBUG] sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}
func (tp *grantService) GrantSysPriv(tf ResourceGrantSystemPrivilege) error {
	log.Printf("[DEBUG] GrantSysPriv grantee: %s\n", tf.Grantee)
	sqlCommand := fmt.Sprintf("grant %s to %s", tf.Privilege, tf.Grantee)
	log.Printf("[DEBUG] sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (tp *grantService) RevokeSysPriv(tf ResourceGrantSystemPrivilege) error {
	log.Printf("[DEBUG] RevokeSysPriv grantee: %s\n", tf.Grantee)
	sqlCommand := fmt.Sprintf("revoke %s from %s", tf.Privilege, tf.Grantee)
	log.Printf("[DEBUG] sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (tp *grantService) RevokeRolePriv(tf ResourceGrantRolePrivilege) error {
	log.Printf("[DEBUG] RevokeSysPriv grantee: %s\n", tf.Grantee)
	sqlCommand := fmt.Sprintf("revoke %s from %s", tf.Role, tf.Grantee)
	log.Printf("[DEBUG] sqlcommand: %s\n", sqlCommand)
	_, err := tp.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}
