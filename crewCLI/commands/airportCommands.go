package commands

import (
	"context"
	"crewCLI/db"
	"database/sql"
	"fmt"
	"log"
)

// Entry function to handle any airport related command
func HandleAirports(args []string) {
	if len(args) < 2 {
		log.Fatal("Too few args! For help, invoke this with the flags: '-a -h'")
	}
	switch args[1] {
	case "-i":
		if len(args) < 3 {
			log.Fatalf("Too few args!")
		}
		// Check to see if we have an airport comment, otherwise pass nil as airport comment
		var comment string
		if len(args) >= 4 {
			comment = args[3]
		}
		insertAirport(args[2], comment)
	case "-d":
		if len(args) < 3 {
			log.Fatalf("Too few args!")
		}
		deleteAirport(args[2])
	case "-l":
		listAirports()
	case "-h":
		fmt.Println("Usage: ./crewCLI -a -[<option>] [arg] [comment-optional]")
		fmt.Println("\t -i [IATA_Code] \"[comment]\" \t insert airport")
		fmt.Println("\t -d [IATA_code] \t\t delete airport")
		fmt.Println("\t -l \t\t\t\t list airports")
		fmt.Println("\t -h \t\t\t\t airport help")
	default:
		log.Fatalf("Invalid input. For help with airport actions run with -a -h")
	}
}

// Adds airport to Airports table in database
func insertAirport(airport_name string, comment string) {
	query := "INSERT INTO Airports (IATA_Code, comment) VALUES (?, ?)"
	insertResult, err := db.DB.ExecContext(context.Background(), query, airport_name, comment)
	if err != nil {
		log.Fatalf("Impossible insert into Airports: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Impossible to retrieve last inserted id: %s", err)
	}
	fmt.Println("Insert successful")
}

// Deletes airport from Airports table in database
func deleteAirport(airport_name string) {
	query := "DELETE FROM Airports WHERE IATA_Code = ?"
	deleteResult, err := db.DB.ExecContext(context.Background(), query, airport_name)
	if err != nil {
		log.Fatalf("Impossible delete into Airports: %s", err)
	}
	ra, err := deleteResult.RowsAffected()
	if err != nil {
		log.Fatalf("Impossible to know number of rows affected: %s", err)
	}

	fmt.Printf("Row deleted: %v\n", ra)
}

// Lists all airports in the Airports table in the database
func listAirports() {
	selectResult, err := db.DB.Query("SELECT * FROM Airports")

	if err != nil {
		log.Fatalf("Impossible to select airports: %s", err)
	}
	defer selectResult.Close()
	fmt.Println("\n--- Airports ---")
	for selectResult.Next() {
		var airport string
		var comment sql.NullString
		err = selectResult.Scan(&airport, &comment)

		if err != nil {
			log.Fatalf("Impossible to get row from selected results: %s", err)
		}

		fmt.Printf("%s", airport)
		if comment.Valid {
			fmt.Printf(": %s\n", comment.String)
		}
		fmt.Println()
	}
}
