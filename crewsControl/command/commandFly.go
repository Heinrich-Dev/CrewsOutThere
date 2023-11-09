package command

import (
	"crewFinder/db"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

/*
If "Yes":
* - send the responders name and phone # back to the original sender (req phone # entry in contacts)
* - add everyone associated with that timeStamp's top entry in deferred to contacts (and remove those entries from deferred)
* - send everyone assoc. with that timeStamp the message in from their top entry of deferred
* - remove everyone associated with that timestamp from contacts

* If "No":
* - remove them from contacts
* - add their top entry of deferred to contacts (and remove that entry from deferred)
* - send them the message from the top entry of deferred
*/
func handleContacted(input string, contactedPhone string) string {
	requesterPhone, _, timeStamp := getItemFromCont(contactedPhone)

	response := ""

	if strings.ToLower(input) == "yes" {
		removeRequestFromRequesterAtTimeStamp(timeStamp)
		MessageUser(requesterPhone, buildReponseMessage(contactedPhone))
		handleOutgoingContacts(contactedPhone, timeStamp)
		response = getNameOfMember(requesterPhone) + " has been notified of your response."

	} else if strings.ToLower(input) == "no" {
		response = moveSingleEntryFromDefToCont(contactedPhone, timeStamp)
	} else {
		// If user did not say yes or no, prompt them again
		response = "Please confirm Yes or No to this request."
	}

	return response
}

// This queries for everyone requested at the given timestamp to tell them that that request has been filled, then
// send them their next request
func handleOutgoingContacts(contactedPhone string, timeStamp int64) {
	// First, get all of the people who have been contacted by the original requester as well as the person who
	// accepted the request (in order to send them their next deferred message)
	// stmt, err := db.DB.Prepare("SELECT * FROM Contacts WHERE contacted_phone = ? OR timestamp = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Contacts WHERE contacted_phone = ? OR timestamp = ?"
	toBeContacted, err := db.DB.Query(query, contactedPhone, timeStamp)

	if err != nil {
		log.Fatalf("Impossible select from Contacts: %s", err)
	}

	defer toBeContacted.Close()

	// Then, remove the original requester from Contacts
	removeRequestFromContactsAtContactedPhone(contactedPhone)

	// Next, for each of the results, move their oldest entry from Deferred into Contacts, then contact them
	for toBeContacted.Next() {
		var rPhone string
		var cPhone string
		err = toBeContacted.Scan(&rPhone, &cPhone, &timeStamp)

		if err != nil {
			log.Fatalf("Impossible to get row from selected results: %s", err)
		}

		message := moveSingleEntryFromDefToCont(cPhone, timeStamp)

		MessageUser(cPhone, message)
	}
}

// Returns the new message to send as found in top entry of def_contacts, or returns that they have no remaining requests
func moveSingleEntryFromDefToCont(contactedPhone string, timeStamp int64) string {
	// Get the message from the top Deferred table entry
	newRequesterPhone, _, message, newRequestTimeStamp := getTopItemFromDef(contactedPhone)

	// Remove the old entry from contacts (happens either way) (must happen before the potential new entry is added)
	removeRequestFromContactsAtContactedPhone(contactedPhone)

	// If user was deferred, add their item to contacts
	if message != "" {
		addItemToContacts(newRequesterPhone, contactedPhone, newRequestTimeStamp)
	}

	// Remove the new entry from def_cont
	removeRequestFromDeferredAtTimeStamp(newRequestTimeStamp)

	return message
}

// Sends out the given message to the given phone number
// Splits the message up into groups of 480 chars split at newlines. Required for viewing all roles and airports
// Also will split messages not containing newlines in halves recursively. This will take a while to get back to
// the user, but this should rarely happen
func MessageUser(phone string, message string) {
	if message == "" {
		return
	}
	if len(message) > 480 {
		messageList := strings.Split(message, "\n")
		newMessage := ""

		if len(messageList[0]) > 960 {
			MessageUser(phone, messageList[0][0:len(messageList[0])/2])
			MessageUser(phone, messageList[0][len(messageList[0])/2:len(messageList[0])])
		} else {
			for i := 0; i < len(messageList); i++ {
				if len(newMessage)+len(messageList[i])+1 > 480 {
					MessageUser(phone, newMessage)
					newMessage = ""
				}
				newMessage += messageList[i] + "\n"
			}
			MessageUser(phone, newMessage)
		}

	} else {
		fmt.Printf("\t %s, your number is: %s\n", message, phone)
	}
}

// Handles an incoming flight request by messaging relevant users and adding users already in contacted to deferred
func handleFlyRequest(requestPhone string, requestMessage string, role string, airport string, timestamp int64) string {
	// Used to determine if we have found a matching user in which case we add the request to the requester table which also requires a tracking variable as it can happen in two places
	foundMatchingUsers := false

	// Build a properly formatted request message
	requestMessageFull := buildRequestMessage(requestPhone, requestMessage)

	// Add every phone number with matching role and name with number in contacts to be placed in deferred
	// stmt, err := db.DB.Prepare("SELECT Phone_Number FROM Flies NATURAL JOIN Wants NATURAL JOIN Members WHERE Role_Name = ? AND IATA_Code = ? AND Phone_Number != ? AND Notify = 1 AND Phone_Number IN (SELECT Phone_Number FROM Contacts)")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT Phone_Number FROM Flies NATURAL JOIN Wants NATURAL JOIN Members WHERE Role_Name = ? AND IATA_Code = ? AND Phone_Number != ? AND Notify = 1 AND Phone_Number IN (SELECT Phone_Number FROM Contacts)"
	matchingCrew, err := db.DB.Query(query, role, airport, requestPhone)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error querying database in handleflyrequest: %s", err)
		}
	}

	defer matchingCrew.Close()

	// For each user that appears in above query, add them to the deferred table
	for matchingCrew.Next() {
		var contactPhone string
		// Make sure we know we have found matching users
		if !foundMatchingUsers {
			foundMatchingUsers = true
		}
		err = matchingCrew.Scan(&contactPhone)

		if err != nil {
			log.Fatalf("Impossible to get row from selected results: %s", err)
		}
		addItemToDeferred(requestPhone, contactPhone, requestMessageFull, timestamp)
	}

	// If found matching users add request to database
	if foundMatchingUsers {
		addItemToRequester(requestPhone, requestMessageFull, timestamp)
	}
	// Message every phone number with matching role and name with number not in contacts
	// stmt, err = db.DB.Prepare("SELECT Phone_Number FROM Flies NATURAL JOIN Wants NATURAL JOIN Members WHERE Role_Name = ? AND IATA_Code = ? AND Phone_Number != ? AND Notify = 1 AND Phone_Number NOT IN (SELECT Phone_Number FROM Contacts)")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query = "SELECT Phone_Number FROM Flies NATURAL JOIN Wants NATURAL JOIN Members WHERE Role_Name = ? AND IATA_Code = ? AND Phone_Number != ? AND Notify = 1 AND Phone_Number NOT IN (SELECT Phone_Number FROM Contacts)"
	matchingCrew, err = db.DB.Query(query, role, airport, requestPhone)
	if err != nil {
		// It's okay if there's no rows here as there could be valid crew members but they are in contacts
		if err != sql.ErrNoRows {
			log.Fatalf("Error querying database: %s", err)
		}
	}

	defer matchingCrew.Close()

	// Send message out to each user and add them to the contacts table
	for matchingCrew.Next() {
		var contactPhone string
		err = matchingCrew.Scan(&contactPhone)

		// If we haven't added request to database as we hadn't found matching users to add to deferred
		if !foundMatchingUsers {
			foundMatchingUsers = true
			addItemToRequester(requestPhone, requestMessageFull, timestamp)
		}

		if err != nil {
			log.Fatalf("Impossible to get row from selected results: %s", err)
		}

		// First, add to contacts table
		addItemToContacts(requestPhone, contactPhone, timestamp)

		// Then, send message out to user
		// params := &api.CreateMessageParams{}
		// params.SetBody(requestMessageFull)
		// params.SetFrom(crewTwilio.CrewNumber)
		// params.SetTo(contactPhone)
		// _, err := crewTwilio.TwilioClient.Api.CreateMessage(params)
		MessageUser(contactPhone, requestMessageFull)

		if err != nil {
			return err.Error()
		}
	}

	if !foundMatchingUsers {
		return "No users were found to be registered under both " + role + " and " + airport + " with notifications on. Your request could not be created."
	}

	return "Request created"
}

// Takes an input message from a user and formats it to be sent out as the request message to other users
func buildRequestMessage(requestPhone string, message string) string {
	var requesterName string
	// stmt, err := db.DB.Prepare("SELECT Name from Members Where Phone_Number = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()
	query := "SELECT Name from Members Where Phone_Number = ?"
	err := db.DB.QueryRow(query, requestPhone).Scan(&requesterName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("Requester phone number not in members")
		}
		log.Fatalf("Error querying database: %s", err)
	}
	// TODO Check if user message ends in punctuation so if not append a period to validate capital T in type
	requestMessage := "\t" + requesterName + " is building a crew: " + message + " type Yes or No or ignore."
	return requestMessage
}


// Builds a simple response message based on a contacted user
func buildReponseMessage(contactedPhone string) string {
	name := getNameOfMember(contactedPhone)
	message := "\t" + name + " (" + contactedPhone + ") has agreed to your request."

	return message
}

// Adds the given information as a new entry in the Contacts table
func addItemToContacts(requesterPhone string, contactedPhone string, timeStamp int64) {
	// stmt, err := db.DB.Prepare("INSERT INTO Contacts (requester_phone, contacted_phone, timestamp) VALUES (?, ?, ?)")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "INSERT INTO Contacts (requester_phone, contacted_phone, timestamp) VALUES (?, ?, ?)"
	_, err := db.DB.Exec(query, requesterPhone, contactedPhone, timeStamp)
	if err != nil {
		log.Fatalf("Error inserting into Contacts: %s", err)
	}
}

// Adds the given information as a new entry in the requester table
func addItemToRequester(requesterPhone string, request_message string, timestamp int64) {
	// stmt, err := db.DB.Prepare("INSERT INTO Requester (phone_number, request_message, timestamp) VALUES (?, ?, ?)")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "INSERT INTO Requester (phone_number, request_message, timestamp) VALUES (?, ?, ?)"
	_, err := db.DB.Exec(query, requesterPhone, request_message, timestamp)
	if err != nil {
		log.Fatalf("Error inserting into Requester: %s", err)
	}
}

// Adds the given information as a new entry in the deferred table
func addItemToDeferred(requesterPhone string, contactedPhone string, requestMessage string, timeStamp int64) {
	// stmt, err := db.DB.Prepare("INSERT INTO Deferred (requester_phone, contacted_phone, request_message, timestamp) VALUES (?, ?, ?, ?)")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "INSERT INTO Deferred (requester_phone, contacted_phone, request_message, timestamp) VALUES (?, ?, ?, ?)"
	_, err := db.DB.Exec(query, requesterPhone, contactedPhone, requestMessage, timeStamp)
	if err != nil {
		log.Fatalf("Error inserting into Deferred: %s", err)
	}
}

// Removes a request entry from the Contacts table at the given timestamp
func removeRequestFromContactsAtContactedPhone(contactedPhone string) {
	// stmt, err := db.DB.Prepare("DELETE FROM Contacts WHERE contacted_phone = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "DELETE FROM Contacts WHERE contacted_phone = ?"
	_, err := db.DB.Exec(query, contactedPhone)
	if err != nil {
		log.Fatalf("Impossible delete from Contacts: %s", err)
	}
}

// Removes a request entry from the Deferred table at the given timestamp
func removeRequestFromDeferredAtTimeStamp(timeStamp int64) {
	// stmt, err := db.DB.Prepare("DELETE FROM Deferred WHERE timestamp = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "DELETE FROM Deferred WHERE timestamp = ?"
	_, err := db.DB.Exec(query, timeStamp)
	if err != nil {
		log.Fatalf("Impossible delete from Deferred: %s", err)
	}
}

// Removes a request entry from the Requester table at the given timestamp
func removeRequestFromRequesterAtTimeStamp(timeStamp int64) {
	// stmt, err := db.DB.Prepare("DELETE FROM Requester WHERE timestamp = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "DELETE FROM Requester WHERE timestamp = ?"
	_, err := db.DB.Exec(query, timeStamp)
	if err != nil {
		log.Fatalf("Impossible delete from Deferred: %s", err)
	}
}

// Returns all of the information in Contacts at the given contacted number
// There will only ever be one entry in contacts for a given contacted number
func getItemFromCont(contactedPhone string) (string, string, int64) {
	var rPhone string
	var cPhone string
	var timeStamp int64
	// stmt, err := db.DB.Prepare("SELECT * FROM Contacts WHERE contacted_phone = ?")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Contacts WHERE contacted_phone = ?"
	err := db.DB.QueryRow(query, contactedPhone).Scan(&rPhone, &cPhone, &timeStamp)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error querying database: %s", err)
		}
	}

	return rPhone, cPhone, timeStamp
}

// Returns all of the information from the oldest request entry in Deferred for the given phone number
func getTopItemFromDef(contactedPhone string) (string, string, string, int64) {
	var rPhone string
	var cPhone string
	var message string
	var newRequestTimeStamp int64
	// stmt, err := db.DB.Prepare("SELECT * FROM Deferred WHERE contacted_phone = ? ORDER BY timestamp ASC")
	// if err != nil {
	// 	log.Fatalf("Error creating prepared statement: %s", err)
	// }
	// defer stmt.Close()

	query := "SELECT * FROM Deferred WHERE contacted_phone = ? ORDER BY timestamp ASC"
	err := db.DB.QueryRow(query, contactedPhone).Scan(&rPhone, &cPhone, &message, &newRequestTimeStamp)
	if err != nil {
		if err == sql.ErrNoRows {
			// This will be used to check if user is not deferred
			message = ""
		} else {
			log.Fatalf("Error querying database: %s", err)
		}
	}

	return rPhone, cPhone, message, newRequestTimeStamp
}
