package oraclehelper

import (
	"fmt"
	"log"
)

const (
	queryProfile = `
SELECT 
  * 
FROM dba_profiles
WHERE profile = UPPER(:1)
`
)

type (
	//ResourceProfile ...
	ResourceProfile struct {
		Profile      string
		ResourceName string
		Limit        string
	}
	profile struct {
		Profile      string
		ResourceName string
		ResourceType string
		Limit        string
		Common       string
		Inherited    string
		Implicit     string
	}
	profileService struct {
		client *Client
	}
)

func (p *profileService) CreateProfile(tf ResourceProfile) error {
	log.Printf("[DEBUG] CreateProfile profile: %s\n", tf.Profile)
	sqlCommand := fmt.Sprintf("create profile %s limit PASSWORD_GRACE_TIME default", tf.Profile)
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}

func (p *profileService) ReadProfile(tf ResourceProfile) (map[string]string, error) {
	log.Printf("[DEBUG] Read name: %s\n", tf.Profile)
	profileparms := make(map[string]string)

	rows, err := p.client.DBClient.Query(queryProfile, tf.Profile)
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

		var profileparm profile
		if val, ok := m["PROFILE"].(string); ok {
			profileparm.Profile = val
		}
		if val, ok := m["RESOURCE_NAME"].(string); ok {
			profileparm.ResourceName = val
		}
		if val, ok := m["RESOURCE_TYPE"].(string); ok {
			profileparm.ResourceType = val
		}
		if val, ok := m["LIMIT"].(string); ok {
			profileparm.Limit = val
		}
		if val, ok := m["COMMON"].(string); ok {
			profileparm.Common = val
		}
		if val, ok := m["INHERITED"].(string); ok {
			profileparm.Inherited = val
		}
		if val, ok := m["IMPLICIT"].(string); ok {
			profileparm.Implicit = val
		}
		profileparms[profileparm.ResourceName] = profileparm.Limit
	}
	if len(profileparms) > 0 {
		profileparms["PROFILE"] = tf.Profile
	} else {
		return profileparms, fmt.Errorf("Can not find profile: %s in the dba_profiles", tf.Profile)
	}
	return profileparms, nil
}

func (p *profileService) UpdateProfile(tf ResourceProfile) error {
	log.Printf("[DEBUG] UpdateProfile resource: %s, value: %s \n", tf.ResourceName, tf.Limit)
	sqlCommand := fmt.Sprintf("alter profile %s limit %s %s", tf.Profile, tf.ResourceName, tf.Limit)

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[DEBUG] SetParameter return error: %s \n", err)
		return err
	}
	return nil
}

func (p *profileService) ResetProfileResourceLimite(tf ResourceProfile) error {
	log.Printf("[DEBUG] ResetProfileResourceLimite profile: %s, resource: %s", tf.Profile, tf.ResourceName)

	sqlCommand := fmt.Sprintf("alter profile %s limit %s default", tf.Profile, tf.ResourceName)
	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[DEBUG] SetParameter return error: %s \n", err)
		return err
	}
	return nil
}

func (p *profileService) DeleteProfile(tf ResourceProfile) error {
	log.Println("[DEBUG] DeleteProfile")
	sqlCommand := fmt.Sprintf("drop profile %s", tf.Profile)
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}
