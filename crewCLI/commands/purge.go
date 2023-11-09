package commands

import (
	"context"
	"crewCLI/db"
	"fmt"
	"log"
)

// HandlePurge is a helper function of cli.go
// HandlePurge takes in arguments specified on the command line, -h shows how
// each argument is handled.
func HandlePurge(args []string) {
	if len(args) < 2 {
		log.Fatal("Too few arguments. Invoke with flag -h for specific options.")
	}
	switch args[1] {
	case "-h":
		fmt.Println("Usage: ./crewCLI [flag] [phone_number]")
		fmt.Println("phone_number will be removed from all tables within the database.")
		fmt.Println("optional flag -h for help")
		fmt.Println("optional flag -a to purge all users added by phone_number")
		fmt.Println("optional flag -r to purge all users added by phone_number and all numbers added by those numbers and so on")
		fmt.Println("optional flag -n to purge all entries from Requester, Contacts, and Deferred")
		fmt.Println("WARNING: purge assumes first arg is a phone number.")
	case "-a":
		fmt.Printf("Purging all users added by %s from all tables in database...\n", args[2])
		PurgeAddedBy(args[2], 0)
	case "-r":
		fmt.Printf("Purging all users associated with %s from all tables in database...\n", args[2])
		PurgeAddedBy(args[2], 1)
	case "-n":
		fmt.Printf("Purging %s from all tables in database...\n", args[2])
		Nuke(args[2])
	default:
		fmt.Printf("Purging %s from all tables in database...\n", args[1])
		Purge(args[1])
	}
}

// Purge queries the cotdb database, finds all tables that have the phone_number attribute, and
// removes the given phone_number from all tables
func Purge(phone_number string) {
	// get all tables in the database with phone_number attribute
	res, err := db.DB.Query("SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cotdb' AND COLUMN_NAME = 'phone_number'")
	if err != nil {
		log.Fatalf("Error while trying to find tables: %s", err)
	}
	var table string
	for res.Next() {
		err := res.Scan(&table)
		fmt.Printf("--Deleting %s from table %s--\n", phone_number, table)
		if err != nil {
			log.Fatalf("Error while scanning table names into slice: %s", err)
		}
		// determine if individual user is being deleted or all users added by phone number
		query := "DELETE FROM " + table + " WHERE phone_number = ?"
		deleteResult, err := db.DB.ExecContext(context.Background(), query, phone_number)
		if err != nil {
			log.Fatalf("Error while deleting phone numbers from tables: %s", err)
		}

		ra, err := deleteResult.RowsAffected()
		if err != nil {
			log.Fatalf("Error finding rows affected: %s", err)
		}
		fmt.Printf("Rows deleted: %v\n", ra)
	}
}

// Similar to purge but only affects tables related to flight requests
func Nuke(phone_number string) {
	// get all tables in the database with phone_number attribute
	selectResult, err := db.DB.Query("SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'cotdb' AND COLUMN_NAME = 'requester_phone'")
	if err != nil {
		log.Fatalf("Error while trying to find tables: %s", err)
	}
	var table string
	for selectResult.Next() {
		err := selectResult.Scan(&table)
		fmt.Printf("--Deleting %s from table %s--\n", phone_number, table)
		if err != nil {
			log.Fatalf("Error while scanning table names into slice: %s", err)
		}
		// determine if individual user is being deleted or all users added by phone number
		query := "DELETE FROM " + table + " WHERE requester_phone = ?"
		deleteResult, err := db.DB.ExecContext(context.Background(), query, phone_number)
		if err != nil {
			log.Fatalf("Error while deleting phone numbers from tables: %s", err)
		}

		ra, err := deleteResult.RowsAffected()
		if err != nil {
			log.Fatalf("Error finding rows affected: %s", err)
		}
		fmt.Printf("Rows deleted: %v\n", ra)
	}
	fmt.Printf("--Deleting %s from table Requester--\n", phone_number)
	query2 := "DELETE FROM Requester WHERE phone_number = ?"
	deleteResult2, err := db.DB.ExecContext(context.Background(), query2, phone_number)
	if err != nil {
		log.Fatalf("Error while deleting phone numbers from tables: %s", err)
	}

	ra, err := deleteResult2.RowsAffected()
	if err != nil {
		log.Fatalf("Error finding rows affected: %s", err)
	}
	fmt.Printf("Rows deleted: %v\n", ra)
}

// PurgeAddedBy takes a string that is a phone number that has added other users.
// PurgeAddedBy queries the database to find all users added by that phone number
// and removes them using Purge.
// flag specifies if the function will be run recursively
func PurgeAddedBy(added_by string, flag int) {
	fmt.Printf("=== DELETING ALL USERS ADDED BY %s ===\n", added_by)
	// get all users added by phone number
	query := "SELECT phone_number FROM Members WHERE added_by = ?"
	selectResult, err := db.DB.Query(query, added_by)
	if err != nil {
		log.Fatalf("Error while finding members: %s", err)
	}
	var phone_number string
	// loop through all numbers, removing them from the database
	Purge(added_by)
	for selectResult.Next() {
		err := selectResult.Scan(&phone_number)
		if err != nil {
			log.Fatalf("Error while saving phone_number: %s", err)
		}
		if flag == 1 {
			PurgeAddedBy(phone_number, flag)
		}
	}
}
