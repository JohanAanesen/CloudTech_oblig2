package main

import (
	"gopkg.in/mgo.v2/bson"
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
//Data struct
type Data struct {
	Base  string
	Date  string
	Rates map[string]float64
}

type Payload struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency" bson:"targetCurrency"`
	MinTriggerValue float64       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

type InvokedPayload struct {
	BaseCurrency    string  `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency" bson:"targetCurrency"`
	CurrentRate     float64 `json:"currentRate" bson:"currentRate"`
	MinTriggerValue float64 `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

type LatestPayload struct {
	BaseCurrency   string `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency" bson:"targetCurrency"`
}