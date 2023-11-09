package main

import (
	"bufio"
	"crewFinder/command"
	"crewFinder/db"
	"crewFinder/encryption"
	"crewFinder/httpServ"
	"fmt"
	"os"
	"strings"
	"time"

	"net/http"
)

// Initialize http handler functions, connection to database, pull values from cot.conf, setup the parser, connect to twilio, and begin serving requests
func main() {
	http.HandleFunc("/", httpServ.ReceiveText)
	http.HandleFunc("/status", httpServ.ReceiveTest)
	// logging.InitLogger()
	encryption.InitConf()
	db.DBAdminConnect()
	command.ParserSetup()
	// crewTwilio.TwilioConnect()
	go http.ListenAndServe(":3000", nil)
	// ListenAndServe is a blocking function so if we ever get out of it things have gone horribly wrong
	// If you want to respond as a specific number, enter command as: [phone_number]:[rest of message]
	// Otherwise it will default to defaultPhone
	// for {
	// 	db.NestedQueryTest()
	// }
	// fmt.Println("Exiting on 0")
	// os.Exit(0)
	phonePrefix := "%2b"
	// mask := 0
	for {
		// db.NestedQueryTest()
		defaultPhone := "11234567890"
		var inText string
		// if mask%3 == 0 {
		// 	inText = "I need an MO out of kbli"
		// } else if mask%3 == 1 {
		// 	inText = "I need an AP out of kbli"
		// } else {
		// 	inText = "4564564567:Yes"
		// }
		// fmt.Println(mask)
		// mask++
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("\n\nEnter Text: ")
		if scanner.Scan() {
			inText = scanner.Text()
		}

		// If user input phone number, parse it for use
		inputPhone, request := splitStringByFirstColon(inText)
		if inputPhone != "" {
			// Need to prepend a 1 to non-default phones to pass verification
			defaultPhone = "1" + inputPhone
			inText = request
		}
		defaultPhone = phonePrefix + defaultPhone

		fmt.Printf("Phone number: %s\nMessage: %s\n", defaultPhone, inText)

		response := command.ValidateAndParse(inText, defaultPhone, time.Now().UnixMilli())
		fmt.Printf("\nRESPONSE (%d chars): %s\n", len(response), response)
	}

}

// Used to allow us to specify what phone number we want to message as
func splitStringByFirstColon(input string) (string, string) {
	index := strings.Index(input, ":")
	if index == -1 {
		// Return empty strings if colon is not found
		return "", ""
	}

	phone_number := input[:index]
	request := input[index+1:]

	return phone_number, request
}
