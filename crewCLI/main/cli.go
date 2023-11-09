package main

import (
	"crewCLI/commands"
	"crewCLI/db"
	"crewCLI/encryption"
	"fmt"
	"log"
	"os"
)

var badargs = "Invalid input, for help please invoke program with the arg '-h'"

func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Fatalf("No args! %s", badargs)
	}
	// get variables used for authenticating and connecting
	encryption.InitConf()
	db.DBConnect()
	// Remove invocation from args
	args = append(args[:0], args[1:]...)
	switch args[0] {
	case "-a":
		commands.HandleAirports(args)
	case "-r":
		commands.HandleRoles(args)
	case "-p":
		commands.HandlePurge(args)
	case "-h":
		fmt.Println("Usage: ./crewCLI -[command_type] -[<option>] [args]")
		fmt.Println("\t -a \t airport commands")
		fmt.Println("\t -r \t role commands")
		fmt.Println("\t -p \t purge user")
		fmt.Println("\t -h \t for command help")
		fmt.Println("\t any command type followed by a -h will provide hints for formatting")
	default:
		log.Fatalf("Command not recognized! %s", badargs)
	}
}
