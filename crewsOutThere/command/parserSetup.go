package command

import (
	"crewFinder/db"
	"log"
	"regexp"
)

type Status int64

const (
	MEMBER     Status = 0
	CONTACTED  Status = 1
	INVALID    Status = 2
	CONFIRMING Status = 3
	NONMEMBER  Status = 4
)

type Airport struct {
	IATA_Code string
	comment   *string
}

var wantsRegex *regexp.Regexp
var needsRegex *regexp.Regexp

var flyNotifRegex *regexp.Regexp
var helpRegex *regexp.Regexp
var showRegex *regexp.Regexp
var inviteRegex *regexp.Regexp
var roleRegex *regexp.Regexp
var requestRegex *regexp.Regexp

// Regex Group Arrays ~Regey

// These are the first words/phrases said after "I want to ..."
var wantStarters = []string{"be", "fly", "view", "invite"}

// Different ways to say don't. Note: spaces are baked into the regex
var dont = []string{" dont", " don't", " do not"}

// Different ways to say "from"
var from = []string{"from", "out of", "at", "in"}

func ParserSetup() {
	wantsRegexSetup()
	needsRegexSetup()
	flyNotifRegexSetup()
	helpRegexSetup()
	showCommandsSetup()
	inviteSetup()
	roleRegexSetup()
	requestRegexSetup()
}

// Creates regex for the help command
func helpRegexSetup() {
	// Checks if there is a command word after a help word.
	// *Important: Does NOT match if there is no command word after the help word.
	regexString := "(?i)^help\\s*([a-zA-Z\\s]*?)\\s*$"

	// Compile regex (Will crash if something went wrong)
	helpRegex = regexp.MustCompile(regexString)
}

// Creates Regex for enabling and disabling fly notifications
func flyNotifRegexSetup() {

	// This captures whether or not the user wants to fly or not.
	regexString := "(?i)^I" + regexGroup(dont) + "? want to fly[\\s]*" + regexGroup(from) + "?[\\s]*([a-zA-Z]*)?\\s*$"

	// Compile regex (Will crash if something went wrong)
	flyNotifRegex = regexp.MustCompile(regexString)

}

// Creates regex for showing lists to the user
func showCommandsSetup() {
	// Captures if someone wants to show only their or all roles/airports
	regexString := "(?i)^I want to view ([a-zA-Z]*)\\s*(roles|airports|detailed)?\\s*(roles|airports)?"

	// Compile regex (Will crash if something went wrong)
	showRegex = regexp.MustCompile(regexString)
}

// Creates regex for inviting users
func inviteSetup() {
	regexString := "(?i)^I want to invite (\\d{10,11})[\\s]*$"

	// Compile regex (Will crash if something went wrong)
	inviteRegex = regexp.MustCompile(regexString)

}

// Creates regex that allows the user to add or remove themselves to/from a role
func roleRegexSetup() {
	// Captures the role the user specifies, and whether they want to remove or add a role
	regexString := "(?i)^I" + regexGroup(dont) + "? want to be (a|an) ([a-zA-Z]*)"

	// Compile regex (Will crash if something went wrong)
	roleRegex = regexp.MustCompile(regexString)
}

// Creates regex that figures out the command the user wants
func wantsRegexSetup() {
	regexString := "(?i)^I" + regexGroup(dont) + "? want to " + regexGroup(wantStarters)

	// Compile regex (Will crash if something went wrong)
	wantsRegex = regexp.MustCompile(regexString)
}

// Creates regex that figures out the command the user wants
func needsRegexSetup() {
	regexString := "(?i)^I need"

	// Compile regex (Will crash if something went wrong)
	needsRegex = regexp.MustCompile(regexString)
}

func requestRegexSetup() {
	regexString := "(?i)I need (a|an) ([a-zA-Z]*) (at|out of|from|in) ([a-zA-Z]*)[\\s]?(.*)"

	// Compile regex (Will crash if something went wrong)
	requestRegex = regexp.MustCompile(regexString)
}

// Finds all airports from the database, returns an Airport struct array with IATA Codes and comments.
func getAirports() ([]Airport, error) {
	rows, err := db.DB.Query("SELECT * FROM Airports")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var l_airports []Airport
	for rows.Next() {
		var row Airport
		if err := rows.Scan(&row.IATA_Code, &row.comment); err != nil {
			log.Fatal(err)
		}

		l_airports = append(l_airports, row)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return l_airports, nil
}

// Generates a group that selects a word from the words string.
func regexGroup(words []string) string {
	// Start Group
	regex := "("

	// Add all command words to regex
	for _, word := range words {
		regex += "\\b" + word + "\\b|"
	}

	// Replace final pipe with a ")" (End group)
	regex = regex[:len(regex)-1] + ")"

	return regex
}

// Generates a group that selects and word from the words string.
func regexWord(word string) string {
	return "(\\b" + word + "\\b)"
}
