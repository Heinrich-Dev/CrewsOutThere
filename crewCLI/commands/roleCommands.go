package commands

import (
	"context"
	"crewCLI/db"
	"database/sql"
	"fmt"
	"log"
)

// Entry function to handle any role related command
func HandleRoles(args []string) {
	if len(args) < 2 {
		log.Fatal("Too few args! For help, invoke this with the flags: '-r -h'")
	}
	switch args[1] {
	case "-i":
		if len(args) < 3 {
			log.Fatalf("Too few args!")
		}
		var desc string
		if len(args) == 4 {
			desc = args[3]
		}
		insertRole(args[2], desc)
	case "-d":
		if len(args) < 3 {
			log.Fatalf("Too few args!")
		}
		deleteRole(args[2])
	case "-l":
		listRoles()
	case "-h":
		fmt.Println("Usage: ./crewCLI -r -[<option>] [arg] \"[description-optional]\"")
		fmt.Println("\t -i [role_name] \t insert role")
		fmt.Println("\t -d [role_name] \t delete role")
		fmt.Println("\t -l \t\t\t list roles")
		fmt.Println("\t -h \t\t\t role help")
	default:
		log.Fatalf("Invalid input. For help with role actions run with -r -h")
	}
}

// Adds a role entry to the Roles table
// role_desc is optional and can be nil
func insertRole(role_name string, role_desc string) {
	query := "INSERT INTO Roles (Role_Name, description) VALUES (?, ?)"
	insertResult, err := db.DB.ExecContext(context.Background(), query, role_name, role_desc)
	if err != nil {
		log.Fatalf("Impossible insert into ROLES: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Impossible to retrieve last inserted id: %s", err)
	}
	fmt.Printf("Insert successful")
}

// Removes the specified role_name from the Roles table
func deleteRole(role_name string) {
	query := "DELETE FROM Roles WHERE Role_Name = ?"
	deleteResult, err := db.DB.ExecContext(context.Background(), query, role_name)
	if err != nil {
		log.Fatalf("Impossible delete into ROLES: %s", err)
	}
	ra, err := deleteResult.RowsAffected()
	if err != nil {
		log.Fatalf("Impossible to know number of rows affected: %s", err)
	}

	fmt.Printf("Row deleted: %v\n", ra)
}

// Lists all role entries in the Roles table
func listRoles() {
	selectResult, err := db.DB.Query("SELECT * FROM Roles")

	if err != nil {
		log.Fatalf("impossible to select roles: %s", err)
	}
	defer selectResult.Close()
	fmt.Println("\n--- Roles ---")
	for selectResult.Next() {
		var role string
		// desc can be null, sql.NullString is a struct that can be null but has a string field if it exists
		var desc sql.NullString
		err = selectResult.Scan(&role, &desc)

		if err != nil {
			log.Fatalf("Impossible to get row from selected results: %s", err)
		}
		fmt.Printf("%s", role)

		if desc.Valid {
			fmt.Printf(": %s\n", desc.String)
		}

	}
	fmt.Println()
}
