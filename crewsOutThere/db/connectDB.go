package db

import (
	// "context"
	"crewFinder/encryption"
	"database/sql"
	"time"

	// "time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var userPassword string
var adminPassword string
var db_password_err error

type Crew struct {
	firstName string
	lastName  string
}

// Connect to the cot database as a user
func DBConnect() {
	userPassword, adminPassword = encryption.GetDBPasswords()
	db_ip := encryption.GetDBIp()

	// Using the dbuser login info
	cfg := mysql.Config{
		User:   "cot.user",
		Passwd: userPassword,
		Net:    "tcp",
		Addr:   db_ip,
		DBName: "cotdb",
	}
	// Get a database handle
	var err error
	// If db hasn't started yet or is down, infinitely loop until we connect
	for {
		DB, err = sql.Open("mysql", cfg.FormatDSN())
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
	}
}

// Connect to the cot database as an admin
func DBAdminConnect() {
	userPassword, adminPassword = encryption.GetDBPasswords()
	db_ip := encryption.GetDBIp()
	cfg := mysql.Config{
		User:   "cot.admin",
		Passwd: adminPassword,
		Net:    "tcp",
		Addr:   db_ip,
		DBName: "cotdb",
	}
	// Get a database handle
	var err error
	for {
		DB, err = sql.Open("mysql", cfg.FormatDSN())
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
	}
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(10)
	DB.SetConnMaxLifetime(time.Second * 20)
}
