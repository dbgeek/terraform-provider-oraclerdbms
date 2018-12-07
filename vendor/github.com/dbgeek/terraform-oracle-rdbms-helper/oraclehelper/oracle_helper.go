package oraclehelper

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	//oci8 need for the Oracle OCI connections
	_ "github.com/mattn/go-oci8"
)

type (
	// Cfg bla bla
	Cfg struct {
		Username  string
		Password  string
		DbHost    string
		DbPort    string
		DbService string
	}
	// Client fkfkkf
	Client struct {
		DBClient               *sql.DB
		DBVersion              *version.Version
		ParameterService       *parameterService
		ProfileService         *profileService
		UserService            *userService
		RoleService            *roleService
		GrantService           *grantService
		StatsService           *statsService
		SchedulerWindowService *schedulerWindowService
		AutoTaskService        *autoTaskService
		DatabaseService        *databaseService
	}
)

const (
	queryDbVersion = `
SELECT 
	version 
FROM v$instance
`
)

// NewClient fkfkf
func NewClient(cfg Cfg) *Client {
	var err error
	var db *sql.DB
	var dBVersion string
	if cfg.DbHost == "" && cfg.DbPort == "" {
		db, err = sql.Open("oci8", fmt.Sprintf("%s/%s@%s", cfg.Username, cfg.Password, cfg.DbService))
		if err != nil {
			log.Fatal(err)

		}
	} else {
		//user/name@host:port/sid
		log.Printf("[DEBUG] dbhost connection string, username: %s, password: %s, dbhost: %s, dbport: %s, dbservice: %s \n", cfg.Username, cfg.Password, cfg.DbHost, cfg.DbPort, cfg.DbService)
		db, err = sql.Open("oci8", fmt.Sprintf("%s/%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.DbHost, cfg.DbPort, cfg.DbService))

		log.Printf("[DEBUG] connection str %s", fmt.Sprintf("%s/%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.DbHost, cfg.DbPort, cfg.DbService))

		if err != nil {
			log.Fatal(err)

		}
		db.Ping()
		if err != nil {
			log.Printf("[DEBUG] ping failed")
			log.Fatal(err)

		}
	}

	err = db.QueryRow(queryDbVersion).Scan(&dBVersion)
	if err != nil {
		log.Fatalf("Query db version failed and return error: %v\n", err)
	}

	c := &Client{DBClient: db}
	c.ParameterService = &parameterService{client: c}
	c.ProfileService = &profileService{client: c}
	c.UserService = &userService{client: c}
	c.RoleService = &roleService{client: c}
	c.GrantService = &grantService{client: c}
	c.StatsService = &statsService{client: c}
	c.SchedulerWindowService = &schedulerWindowService{client: c}
	c.AutoTaskService = &autoTaskService{client: c}
	c.DatabaseService = &databaseService{client: c}
	c.DBVersion, _ = version.NewVersion(dBVersion)
	log.Printf("[DEBUG] dbversion: %v", c.DBVersion)

	return c
}
