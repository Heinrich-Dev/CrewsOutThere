package crewTwilio

import (
	"log"

	"iscotDown/encryption"

	"github.com/twilio/twilio-go"
)

// This client needs to be callable from other functions in package
var TwilioClient *twilio.RestClient
var CrewNumber string = "+1**********"

func TwilioConnect() {
	sid, auth_token := encryption.GetTwilioInfo()
	var myParams twilio.ClientParams
	myParams.Username = sid
	myParams.Password = auth_token

	TwilioClient = twilio.NewRestClientWithParams(myParams)
	if TwilioClient == nil {
		log.Fatal("Bad client")
	}
}
