package main

import (
	"net/http"
	"bytes"
	"encoding/json"
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
 type Payload struct{
 	webhookURL string
 	baseCurrency string
 	targetCurrency string
 	minTriggerValue json.Number
 	maxTriggerValue json.Number
 }

func sendDiscord(w http.ResponseWriter, r *http.Request){
	var jsonStr = []byte(`{"content":"Fuck you."}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}