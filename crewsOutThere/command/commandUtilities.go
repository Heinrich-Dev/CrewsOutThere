package command

import (
	"context"
	"crewFinder/db"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Validates what status the user is to the system
// MEMBER is a fully verified member in the system
// CONTACTED is when they are a fully verified member and have also been contacted by an active request
// INVALID is when invited but haven't given their name
// CONFIRMING is when they have given their name and they need to confirm that it is correct
// NONMEMBER is phone number not in system
func validateUser(phone string) Status {
	var pNumber string
	var added_by string
	var name string
	var notify string
	var isValid int
	var timestampMem int64

	// stmt, err := db.DB.Prepare("SELECT * FROM Members WHERE phone_number = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Members WHERE phone_number = ?"
	err := db.DB.QueryRow(query, phone).Scan(&pNumber, &name, &added_by, &notify, &isValid, &timestampMem)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error querying database: %s", err)
		}
		return NONMEMBER
	} else if name == "" {
		return INVALID
	} else if isValid == 0 {
		return CONFIRMING
	}

	// Check if the user has been contacted
	var rNumber string
	var cNumber string
	var timestamp int64

	// stmt, err = db.DB.Prepare("SELECT * FROM Contacts WHERE contacted_phone = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query = "SELECT * FROM Contacts WHERE contacted_phone = ?"
	err = db.DB.QueryRow(query, phone).Scan(&rNumber, &cNumber, &timestamp)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error querying database: %s", err)
		}
		return MEMBER
	} else {
		return CONTACTED
	}
}

// Check to see if a given IATA Code is in the Airports table
func isIATAInAirports(iata_code string) bool {
	// stmt, err := db.DB.Prepare("SELECT * FROM Airports WHERE iata_code = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	var comment sql.NullString

	query := "SELECT * FROM Airports WHERE iata_code = ?"
	err := db.DB.QueryRow(query, iata_code).Scan(&iata_code, &comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatalf("Error querying database: %s", err)
	}
	return true
}

// Check to see if a give role name is in roles
func isRoleInRoles(roleName string) bool {
	// stmt, err := db.DB.Prepare("SELECT * FROM Roles WHERE role_name = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	var roleMessage sql.NullString

	query := "SELECT * FROM Roles WHERE role_name = ?"
	err := db.DB.QueryRow(query, roleName).Scan(&roleName, &roleMessage)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatalf("Error querying database: %s", err)
	}
	return true
}

// Check to see if given phone number is already flying at the given airport
func isMemberAlreadyFlyingAtAirport(phone string, iata string) bool {
	// stmt, err := db.DB.Prepare("SELECT * FROM Flies WHERE phone_number = ? AND iata_code = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Flies WHERE phone_number = ? AND iata_code = ?"
	err := db.DB.QueryRow(query, phone, iata).Scan(&phone, &iata)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatalf("Error querying database: %s", err)
	}
	return true
}

// Check to see if given phone number is already wanting the given role
func isMemberAlreadyWantingRole(phone string, role string) bool {
	// stmt, err := db.DB.Prepare("SELECT * FROM Wants WHERE phone_number = ? AND role_name = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Wants WHERE phone_number = ? AND role_name = ?"
	err := db.DB.QueryRow(query, phone, role).Scan(&phone, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatalf("Error querying database: %s", err)
	}
	return true
}

// Get the name of a member based on their phone number and return it
func getNameOfMember(phone string) string {
	var name string
	// stmt, err := db.DB.Prepare("Select name FROM Members WHERE phone_number = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "Select name FROM Members WHERE phone_number = ?"
	err := db.DB.QueryRow(query, phone).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows returned")
		}
		log.Fatalf("Error querying database: %s", err)
	}
	return name
}

// Adds the phone number and added_by number a user being invited to CrewsOutThere
// Note that the phone number validations occur in inviteUser which calls this function
func addUnverifiedMember(phone string, added_by string, timestamp int64) {
	name := ""
	notify := 1
	isValid := 0
	query := "INSERT INTO Members (phone_number, name, added_by, notify, is_valid, timestamp) Values (?, ?, ?, ?, ?, ?)"
	insertResult, err := db.DB.ExecContext(context.Background(), query, phone, name, added_by, notify, isValid, timestamp)
	if err != nil {
		log.Fatalf("impossible insert into Members: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
}

// Adds a basic test member to the Members table
func AddTestMember() {
	phone := "11234567890"
	name := "TEST MEMBER"
	added_by := "11234567890"
	notify := 1
	isValid := 1
	timestamp := 0
	query := "INSERT INTO Members (phone_number, name, added_by, notify, is_valid, timestamp) VALUES (?, ?, ?, ?, ?, ?)"
	insertResult, err := db.DB.ExecContext(context.Background(), query, phone, name, added_by, notify, isValid, timestamp)

	if err != nil {
		log.Fatalf("Impossible insert into Members: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Impossible to retrieve last inserted id: %s", err)
	}
}

// Display the contents of the members table. For debugging
func DisplayMembers() {
	selectResult, err := db.DB.Query("SELECT * FROM Members")

	if err != nil {
		log.Fatalf("impossible to select from Members: %s", err)
	}
	defer selectResult.Close()
	fmt.Printf("\n--- Members ---\n")
	for selectResult.Next() {
		var phone string
		var name string
		var added_by string
		var notify int
		var timestamp int64
		err = selectResult.Scan(&phone, &name, &added_by, &notify, &timestamp)

		if err != nil {
			log.Fatalf("impossible to get row from selected results: %s", err)
		}

		notify_text := "Do not notify"

		if notify == 1 {
			notify_text = "Do notify"
		}

		fmt.Printf("%s: %s. Added by %s. %s. Time added: %s.\n", name, phone, added_by, notify_text, timestamp)
	}

	err = selectResult.Err()
	if err != nil {
		log.Fatalf("Error with select query: %s", err)
	}
}

// Changes the input message from unicode to a usable string
func convertTextMessage(input string) string {
	// Convert "+" to " "
	ret := ""
	split := strings.Split(input, "+")
	for i := 0; i < len(split)-1; i++ {
		ret += split[i] + " "
	}
	ret += split[len(split)-1]

	// Convert %22 to """
	split = strings.Split(ret, "%22")
	ret = ""
	for i := 0; i < len(split)-1; i++ {
		ret += split[i] + "'"
	}
	ret += split[len(split)-1]

	// Convert "%27" to "'"
	split = strings.Split(ret, "%27")
	ret = ""
	for i := 0; i < len(split)-1; i++ {
		ret += split[i] + "'"
	}
	ret += split[len(split)-1]

	// Trim off leading and trailing spaces
	ret = strings.TrimSpace(ret)
	return ret
}

// This is to try to handle a caught out of sync mysql error
// Because the error only happens at notifRequestersOfTimedOutRequest
// Which happens to be the first query on a new pass through the system
// This function will be tailored to work with that
func retryQuery(query string, cutoff int64) *sql.Rows {
	retryResult, err := db.DB.Query(query, cutoff)
	if err != nil {
		log.Fatalf("Error returned when handling query error: %s\n", err)
	}
	log.Println("Handled sync error")
	return retryResult
}
