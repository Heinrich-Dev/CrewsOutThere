package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

func AddTestTable() {
	query := "CREATE TABLE IF NOT EXISTS CREW(firstname VARCHAR(255) PRIMARY KEY,lastname VARCHAR(255))"
	// In case of any issues with code or network time out after 5 seconds
	ctx, cancelfunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelfunc()
	_, err := DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating crew table", err)
	} else {
		log.Println("Added Table!")
	}
	// Uncomment these lines when above section of code works
	insertRow("***", "***")
	insertRow("***", "***")
	insertRow("***", "***")
	insertRow("***", "***")
	crew, err := getCrew()
	if err != nil {
		log.Fatal(err)
	}
	for _, member := range crew {
		fmt.Printf("Firstname: %s Lastname: %s\n", member.firstName, member.lastName)
	}
	// Remove Test table to clean things up
	RemoveTestTable()
}

func RemoveTestTable() {
	if _, err := DB.Exec("DROP TABLE CREW"); err != nil {
		log.Print(err)
	}
}

func insertRow(firstname string, lastname string) {
	query := "INSERT INTO CREW (firstname, lastname) VALUES (?, ?)"
	insertResult, err := DB.ExecContext(context.Background(), query, firstname, lastname)
	if err != nil {
		log.Fatalf("impossible insert into CREW: %s", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	fmt.Printf("inserted id: %d\n", id)
}

func getCrew() ([]Crew, error) {
	rows, err := DB.Query("Select * from CREW")
	if err != nil {
		return nil, err
	}
	// Releases any resources held by the rows no matter how the function returns.
	// Looping all the way through rows also closes it implicitly
	defer rows.Close()

	var crew []Crew
	// Loop through the rows
	for rows.Next() {
		var cr Crew
		if err := rows.Scan(&cr.firstName, &cr.lastName); err != nil {
			return crew, err
		}
		crew = append(crew, cr)
	}
	if err = rows.Err(); err != nil {
		return crew, err
	}
	return crew, nil
}
