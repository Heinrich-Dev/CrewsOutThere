package db

import (
	"crewCLI/encryption"
	"database/sql"
	"log"

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

// Connects to the database with 'user' permissions
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
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}
