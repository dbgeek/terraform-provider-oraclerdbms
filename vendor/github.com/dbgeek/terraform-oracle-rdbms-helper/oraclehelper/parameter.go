package oraclehelper

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	queryParameter = `
SELECT
	num as num,
	name,
	type as type,
	value,
	default_value,
	display_value,
	isdefault,
	isses_modifiable,
	issys_modifiable,
	isinstance_modifiable,
	description,
	update_comment,
	hash as col_hash
FROM
	v$system_parameter
WHERE name = :1
`
)

type (
	//ResourceParameter ..
	ResourceParameter struct {
		Name    string
		Value   string
		Comment string
	}
	parameter struct {
		Num                  sql.NullString
		Name                 string
		ParameterType        sql.NullString
		Value                string
		DefaultValue         string
		DisplayValue         sql.NullString
		IsDefault            string
		IsSesModifiable      sql.NullString
		IsSysModifiable      sql.NullString
		IsInstanceModifiable sql.NullString
		Description          sql.NullString
		UpdateComment        sql.NullString
		ParameterHash        sql.NullString
	}
	parameterService struct {
		client *Client
	}
)

func (p *parameterService) Read(tf ResourceParameter) (*parameter, error) {
	log.Printf("[DEBUG] Read name: %s\n", tf.Name)
	param := &parameter{}

	err := p.client.DBClient.QueryRow(queryParameter, tf.Name).Scan(&param.Num,
		&param.Name,
		&param.ParameterType,
		&param.Value,
		&param.DefaultValue,
		&param.DisplayValue,
		&param.IsDefault,
		&param.IsSesModifiable,
		&param.IsSysModifiable,
		&param.IsInstanceModifiable,
		&param.Description,
		&param.UpdateComment,
		&param.ParameterHash,
	)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (p *parameterService) SetParameter(tf ResourceParameter) error {
	log.Printf("[DEBUG] SetParameter value: %s \n", tf.Value)
	sqlCommand := fmt.Sprintf("alter system set %s = %s", tf.Name, tf.Value)

	if tf.Comment != "" {
		sqlCommand += fmt.Sprintf(" comment='%s'", tf.Comment)
	}

	sqlCommand += " scope=both"

	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)

	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		log.Printf("[DEBUG] SetParameter return error: %s \n", err)
		return err
	}
	return nil
}

func (p *parameterService) ResetParameter(tf ResourceParameter) error {
	log.Println("[DEBUG] ResetParameter")
	sqlCommand := fmt.Sprintf("alter system reset %s scope=both", tf.Name)

	_, err := p.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}

	return nil
}
