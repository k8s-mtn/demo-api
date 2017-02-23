package main

import (
	"fmt"

	"github.com/sfreiberg/gotwilio"
)

var (
	twilio *gotwilio.Twilio

	twilioApplicationSid string
)

func setupTwilio(id string, token string, applicationID string) {
	twilioApplicationSid = applicationID
	twilio = gotwilio.NewTwilioClient(id, token)
}

func sendTwilio(number string, imgURL string) error {
	if twilio == nil {
		return fmt.Errorf("twilio is not configured")
	}

	_, _, err := twilio.SendMMS("", number, "", imgURL, "", twilioApplicationSid)
	if err != nil {
		return err
	}

	return nil
}
