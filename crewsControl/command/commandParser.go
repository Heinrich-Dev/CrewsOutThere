package command

import (
	"fmt"
	"strings"
)

// Tells the parser what to do based on the status of the user
func ValidateAndParse(input string, phone string, timeStamp int64) string {
	// Convert the text from "+" delimited to " " and change and "%27" to "'"
	input = convertTextMessage(input)
	// Default help request is handled by twilio so we do not want to write back a text
	if strings.ToLower(input) == "help" {
		return "DEFAULT_HELP_REQUEST"
	}
	// Remove the %2B from the beginning of the phone number
	phone = phone[3:len(phone)]

	phoneWasDeferred := cleanupDB(phone, timeStamp)

	userStatus := validateUser(phone)

	response := "Error. Invalid user status."

	if phoneWasDeferred {
		userStatus = MEMBER
	}

	switch userStatus {
	case MEMBER:
		if strings.ToLower(input) == "yes" {
			response = "The request you may be responding to is no longer available"
		} else {
			response = parse(input, phone, timeStamp)
		}
	case CONTACTED:
		response = handleContacted(input, phone)
	case INVALID:
		setNameOfMember(phone, input)
		response = "Please confirm: is \"" + input + "\" your name? Type Yes or No."
	case CONFIRMING:
		if strings.ToLower(input) == "yes" {
			setMemberValidity(1, phone)
			response = "You have been successfully added to CrewsOutThere. Type \"help\" to get started."
		} else {
			setNameOfMember(phone, "")
			response = "Please respond with your name to be added to CrewsOutThere."
		}
	case NONMEMBER:
		response = "You are not a member of CrewsOutThere."
	}

	return response
}

// Parses the message and attempts to determine the user's intent.
// Returns a bool that represents if we find a command, along with information needed to fulfil the command.
// Reurns a string that represents the response to send back to the user.
func parse(input string, phone string, timeStamp int64) string {
	// This will be sent back to the user
	response := "Sorry I couldn't understand you. Type \"help\" for help."

	if wantsRegex.MatchString(input) {
		groups := wantsRegex.FindStringSubmatch(input)

		switch groups[2] {
		case "fly":
			if flyNotifRegex.MatchString(input) {
				response = handleFlyNotif(input, phone)
			} else {
				// Shouldn't be possible
				response = "Flight status not changed\n" + FlyUsage()
			}
		case "view":
			if showRegex.MatchString(input) {
				response = handleShow(input, phone)
			} else {
				response = ShowUsage()
			}
		case "invite":
			if inviteRegex.MatchString(input) {
				response = handleInvite(input, phone, timeStamp)
			} else {
				response = "Couldn't invite.\n" + InviteUsage()
			}
		case "be":
			if roleRegex.MatchString(input) {
				response = handleRole(input, phone)
			} else {
				response = "Couldnt set your role.\n" + RoleUsage()
			}
		}

		return response
	}

	if helpRegex.MatchString(input) {
		return handleHelp(input)
	}

	if needsRegex.MatchString(input) {
		if requestRegex.MatchString(input) {
			return handleRequest(input, phone, timeStamp)
		}

		// Did not type command correctly
		return "Request not sent.\n" + RequestUsage()
	}

	return response
}

/* Handler functions for the input parser
   Precondition: regex used in the function matches successfully */

// This function determines whether or not the user wants to enable flying notifications
// Returns a string that represents what the system will respond to the user
func handleFlyNotif(input string, phone string) string {
	groups := flyNotifRegex.FindStringSubmatch(input)
	var response string

	enableNotifs := groups[1] == ""
	specifiedAirport := groups[2] != ""
	iata := strings.ToUpper(groups[3])

	// Catch for if users type: "I want to fly kbli"
	if (specifiedAirport == false) && (iata != "") {
		response = "Flight status not changed\n" + FlyUsage()
	} else if enableNotifs {

		if specifiedAirport {
			err := addToFlies(iata, phone)
			if err != nil {
				return err.Error()
			}
			response = "You will now receive notifications for " + iata

		} else {
			err := updateNotify(phone, 1)
			if err != nil {
				return err.Error()
			}
			response = "You will now recieve flight notifications."
		}
	} else {

		if specifiedAirport {
			err := removeUserAtIATAFromFlies(phone, iata)
			if err != nil {
				return err.Error()
			}
			response = "You will no longer receive notifications for " + iata

		} else {
			err := updateNotify(phone, 0)
			if err != nil {
				return err.Error()
			}
			response = "You will no longer recieve flight notifications."
		}
	}

	return response
}

// This function provides a response to send to the user based on some key word
// Returns a string that represents what the system will respond to the user
func handleHelp(input string) string {
	groups := helpRegex.FindStringSubmatch(input)

	// This word is the topic they need help with
	switch strings.ToLower(groups[1]) {
	case "fly":
		return "Enables or disables all flight notifications\n" + FlyUsage()
	case "set role":
		return "Allows you to declare your role\n" + RoleUsage()
	case "view roles":
		return "Displays your roles\n" + ShowRoleUsage()
	case "set airport":
		return "Allows you to set your airports\n" + AirportUsage()
	case "view airports":
		return "Displays your airports\n" + ShowAirportUsage()
	case "invite":
		return "Allows you to invite a number\n" + InviteUsage()
	case "request":
		return "Allows you to request a role for a flight. You will be notified if your request is accepted.\n" + RequestUsage()
	}

	// If word is not in the list (or they just typed "help"), send general help message
	return GeneralUsage()
}

// This function shows a user their roles or airports, depending on which one is requested.
func handleShow(input string, phone string) string {
	groups := showRegex.FindStringSubmatch(input)

	filterUser := groups[1] == "my"

	response := ""

	switch groups[2] {
	case "roles":
		if filterUser {
			response = getEntriesFromWants(phone)
		} else {
			response = getAllRoles()
		}
	case "airports":
		if filterUser {
			response = getEntriesFromFlies(phone)
		} else {
			response = getAllAirports()
		}
	}

	return response
}

// This function invites a number to the service
func handleInvite(input string, requestPhone string, timeStamp int64) string {
	groups := inviteRegex.FindStringSubmatch(input)

	phoneNumber := groups[1]
	outcome := inviteUser(requestPhone, phoneNumber, timeStamp)
	if outcome == "" {
		return "Invited User!"
	}
	return outcome
}

// This function allows a user to add or remove a role
func handleRole(input string, phone string) string {
	groups := roleRegex.FindStringSubmatch(input)

	removingRole := groups[1] != ""
	role := strings.ToUpper(groups[3])

	if role == "" {
		return "Please specify a role."
	}

	response := ""

	if !removingRole {
		err := addToWants(role, phone)
		if err != nil {
			return err.Error()
		}
		response = "You are now a " + role
	} else {
		err := removeUserAtRoleFromWants(phone, role)
		if err != nil {
			return err.Error()
		}
		response = "You are no longer a " + role
	}

	return response
}

// Initiates a flight request and returns a response to the senders
func handleRequest(input string, phone string, timestamp int64) string {
	groups := requestRegex.FindStringSubmatch(input)

	role := strings.ToUpper(groups[2])
	airport := strings.ToUpper(groups[4])
	message := handleFlyRequest(phone, input, role, airport, timestamp)

	return message
}

/* End of handler functions */

// Used for testing
func printArray(arr []string) {
	for i, e := range arr {
		fmt.Println(i, "--", e)
	}
}
