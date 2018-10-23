package oraclehelper

import (
	"fmt"
	"log"
)

const (
	queryProfile = `
SELECT
	profile,
	resource_name,
	resource_type,
	limit,
	common,
	inherited,
	implicit
FROM
	dba_profiles
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
		Commoon      string
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

	profileparms["PROFILE"] = tf.Profile
	rows, err := p.client.DBClient.Query(queryProfile, tf.Profile)
	if err != nil {
		//return nil, err
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var profileparm profile
		rows.Scan(&profileparm.Profile,
			&profileparm.ResourceName,
			&profileparm.ResourceType,
			&profileparm.Limit,
			&profileparm.Commoon,
			&profileparm.Inherited,
			&profileparm.Implicit,
		)
		profileparms[profileparm.ResourceName] = profileparm.Limit
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
