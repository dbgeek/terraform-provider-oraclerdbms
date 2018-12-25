package oraclehelper

import (
	"log"
)

const (
	queryTablePref = `
SELECT 
	DBMS_STATS.GET_PREFS (:1,:2,:3) AS pvalue 
FROM dual
`
	queryGlobalPref = `
SELECT 
	DBMS_STATS.GET_PREFS (:1) AS pvalue 
FROM dual
`
	setGlobalPref = `
BEGIN
	DBMS_STATS.SET_GLOBAL_PREFS (
		pname     => :1,
		pvalue    => :2
	);
END;
`
	setSchemaPref = `
BEGIN
	DBMS_STATS.SET_SCHEMA_PREFS(
		pname	=> :1,
		ownname	=> :2,
		pvalue 	=> :3
	);
END;
`
	setTablePref = `
BEGIN
	DBMS_STATS.SET_TABLE_PREFS(
		pname	=> :1,
		ownname => :2,
		tabname => :3,
		pvalue 	=> :4
	);
END;
`
)

//Stats ..
type (
	//ResourceStats ....
	ResourceStats struct {
		Pname   string
		OwnName string
		TaBName string
		Pvalu   string
	}
	Stats struct {
		Pname   string
		OwnName string
		TaBName string
		Pvalu   string
	}
	statsService struct {
		client *Client
	}
)

func (r *statsService) ReadGlobalPre(tf ResourceStats) (*Stats, error) {
	log.Printf("[DEBUG] ReadGlobalPre pname: %s\n", tf.Pname)
	statsType := &Stats{}

	err := r.client.DBClient.QueryRow(queryGlobalPref, tf.Pname).Scan(&statsType.Pvalu)
	if err != nil {
		return nil, err
	}
	return statsType, nil
}

func (r *statsService) SetGlobalPre(tf ResourceStats) error {
	log.Printf("[DEBUG] SetGlobalPre pname: %s, pvalu: %s\n", tf.Pname, tf.Pvalu)

	_, err := r.client.DBClient.Exec(setGlobalPref, tf.Pname, tf.Pvalu)
	if err != nil {
		return err
	}
	return nil
}

func (r *statsService) SetSchemaPre(tf ResourceStats) error {
	log.Printf("[DEBUG] SetSchemaPre pname: %sowner: %s pvalue: %s\n", tf.Pname, tf.OwnName, tf.Pvalu)

	_, err := r.client.DBClient.Exec(setSchemaPref, tf.Pname, tf.OwnName, tf.Pvalu)
	if err != nil {
		return err
	}
	return nil
}

func (r *statsService) ReadTabPref(tf ResourceStats) (*Stats, error) {
	log.Printf("[DEBUG] ReadTabPref pname: %s owner: %s table: %s\n", tf.Pname, tf.OwnName, tf.TaBName)
	statsType := &Stats{}

	err := r.client.DBClient.QueryRow(queryTablePref, tf.Pname, tf.OwnName, tf.TaBName).Scan(&statsType.Pvalu)
	if err != nil {
		return nil, err
	}
	return statsType, nil
}
func (r *statsService) SetTabPre(tf ResourceStats) error {
	log.Printf("[DEBUG] SetTabPre pname: %s owner: %s table: %s\n", tf.Pname, tf.OwnName, tf.TaBName)

	_, err := r.client.DBClient.Exec(setTablePref, tf.Pname, tf.OwnName, tf.TaBName, tf.Pvalu)
	if err != nil {
		return err
	}
	return nil
}
