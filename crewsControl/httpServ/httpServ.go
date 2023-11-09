package httpServ

import (
	"crewFinder/command"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Receive an HTTP Post from Twilio, pull out relevant fields
func ReceiveText(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	myNumber, myMessage, timeStamp := parsePost(string(body))
	response := command.ValidateAndParse(myMessage, myNumber, timeStamp)
	// Default help request is handled by twilio so we do not want to write back a text
	if response != "DEFAULT_HELP_REQUEST" {
		phone := myNumber[3:len(myNumber)]
		command.MessageUser(phone, response)
	}
}

// Note in here we are dealing with an HTTP Post from Twilio, which should always have the same format, however, want to avoid magic numbers
func parsePost(postBody string) (string, string, int64) {
	var number string
	var message string
	splitPost := strings.Split(postBody, "&")
	// Iterate over all of the indices in the slice looking for the fields we want being body and from
	for _, field := range splitPost {
		// Split each string on the = and compare
		// Every field takes form of fieldName=fieldValue
		arg := strings.Split(field, "=")
		if strings.Compare(arg[0], "From") == 0 {
			number = arg[1]
		} else if strings.Compare(arg[0], "Body") == 0 {
			message = arg[1]
		}
	}
	// Make sure to capture the current time for bookkeeping purposes for the request queue in the database
	return number, message, time.Now().UnixMilli()
}

// This is used so we know the server is still running, responds to crontab pings
func ReceiveTest(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "running")
}
