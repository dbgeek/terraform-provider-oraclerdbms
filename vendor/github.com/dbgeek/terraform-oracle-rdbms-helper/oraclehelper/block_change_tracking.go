package oraclehelper

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	queryBlockChangeTracking = `
SELECT 
    bct.status,
    bct.filename
FROM v$block_change_tracking bct
`
)

type (
	//ResourceBlockChangeTracking ...
	ResourceBlockChangeTracking struct {
		Status   string
		FileName string
	}
	blockChangeTrackingService struct {
		client *Client
	}
)

func (b *blockChangeTrackingService) DisableBlockChangeTracking() error {
	log.Printf("[DEBUG] DisableBlockChangeTracking\n")
	sqlCommand := fmt.Sprintf("ALTER DATABASE DISABLE BLOCK CHANGE TRACKING")
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)
	_, err := b.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}
	return nil
}

func (b *blockChangeTrackingService) EnableBlockChangeTracking(tf ResourceBlockChangeTracking) error {
	log.Printf("[DEBUG] EnableBlockChangeTracking with filename: %s\n", tf.FileName)
	sqlCommand := fmt.Sprintf("ALTER DATABASE ENABLE BLOCK CHANGE TRACKING USING FILE '%s'", tf.FileName)
	log.Printf("[DEBUG] sqlcommand: %s", sqlCommand)
	_, err := b.client.DBClient.Exec(sqlCommand)
	if err != nil {
		return err
	}
	return nil
}

func (b *blockChangeTrackingService) ReadBlockChangeTracking() (*ResourceBlockChangeTracking, error) {
	log.Printf("[DEBUG] ReadBlockChangeTracking\n")
	var fileName sql.NullString
	var status string
	resourceBlockChangeTracking := &ResourceBlockChangeTracking{}
	err := b.client.DBClient.QueryRow(
		queryBlockChangeTracking,
	).Scan(&status,
		&fileName,
	)
	if err != nil {
		return nil, err
	}
	resourceBlockChangeTracking.Status = status
	if fileName.Valid {
		resourceBlockChangeTracking.FileName = fileName.String
	}
	return resourceBlockChangeTracking, nil
}
