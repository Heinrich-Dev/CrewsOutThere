package main

import (
	"iscotDown/crewTwilio"
	"iscotDown/encryption"

	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// Currently contacts *********

// When invoked sends a message to the toContact phone number saying the server is down using twilio
func main() {
	encryption.InitConf()
	crewTwilio.TwilioConnect()
	reporterContact := encryption.GetReporterContact()
	ownerContact := encryption.GetOwnerContact()
	sendMessage(reporterContact)
	sendMessage(ownerContact)

}

func sendMessage(destNumber string) {
	params := &api.CreateMessageParams{}
	params.SetBody("cot is down")
	params.SetFrom(crewTwilio.CrewNumber)
	params.SetTo(destNumber)
	crewTwilio.TwilioClient.Api.CreateMessage(params)
}
