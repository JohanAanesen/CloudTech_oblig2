package funcs

import (
	"gopkg.in/mgo.v2/bson"
)

//Data struct
type Data struct {
	Base  string             `json:"base" bson:"base"`
	Date  string             `json:"date" bson:"date"`
	Rates map[string]float64 `json:"rates" bson:"rates"`
}

//Payload struct
type Payload struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency" bson:"targetCurrency"`
	CurrentRate     float64       `json:"currentRate" bson:"currentRate"`
	MinTriggerValue float64       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

//InvokedPayload struct
type InvokedPayload struct {
	BaseCurrency    string  `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency" bson:"targetCurrency"`
	CurrentRate     float64 `json:"currentRate" bson:"currentRate"`
	MinTriggerValue float64 `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

//LatestPayload struct
type LatestPayload struct {
	BaseCurrency   string `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency" bson:"targetCurrency"`
}

//DiscordWrap struct
type DiscordWrap struct {
	Content InvokedPayload `json:"content"`
}