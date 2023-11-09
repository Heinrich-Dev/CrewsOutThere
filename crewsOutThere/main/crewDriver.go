package main

// #cgo CFLAGS:
// #define _DEFAULT_SOURCE
// #include <stdlib.h>
// #include <unistd.h>
// import "C"
import (
	"crewFinder/command"
	"crewFinder/crewTwilio"
	"fmt"

	"crewFinder/db"
	"crewFinder/encryption"
	"crewFinder/httpServ"
	"crewFinder/logging"
	"log"
	"net/http"
)

// Initialize http handler functions, connection to database, pull values from cot.conf, setup the parser, connect to twilio, and begin serving requests
func main() {
	// C.daemon(1, 0)
	logging.InitLogger()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	encryption.InitConf()
	db.DBAdminConnect()
	fmt.Println("Entering parserSetup")
	command.ParserSetup()
	crewTwilio.TwilioConnect()
	http.HandleFunc("/", httpServ.ReceiveText)
	http.HandleFunc("/status", httpServ.ReceiveTest)
	listenerAddr := encryption.GetListenerIp()
	err := http.ListenAndServe(listenerAddr, nil)
	// ListenAndServe is a blocking function so if we ever get out of it things have gone horribly wrong
	log.Fatal(err)

}
