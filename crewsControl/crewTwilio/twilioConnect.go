package crewTwilio

import (
	"fmt"
	"log"
	"os"

	"crewFinder/encryption"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// This client needs to be callable from other functions in package
var TwilioClient *twilio.RestClient

// This is the crewsOutThere number
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

// Test function to simulate outbound sending of invite message, use for debugging
func TestInvite(phone_number string) string {
	fmt.Println("My auth token is:", os.Getenv("TWILIO_AUTH_TOKEN"))
	params := &api.CreateMessageParams{}
	params.SetBody("This is a message from CrewsOutThere through the Twilio client")
	params.SetFrom(CrewNumber)
	params.SetTo(phone_number)
	resp, err := TwilioClient.Api.CreateMessage(params)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	return *resp.Sid
}
