package main

import (
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

/*
{
    "webhookURL": "http://remoteUrl:8080/randomWebhookPath",
    "baseCurrency": "EUR",
 	"targetCurrency": "NOK",
    "minTriggerValue": 1.50,
    "maxTriggerValue": 2.55
}
*/
type Payload struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency" bson:"targetCurrency"`
	MinTriggerValue float64       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

func sendDiscord(data []byte) {
	//var jsonStr= []byte(`{"content":"Fuck you."}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
