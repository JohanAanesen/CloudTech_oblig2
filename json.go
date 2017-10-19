package main

import (
	"net/http"
	"bytes"
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
func sendDiscord(w http.ResponseWriter, r *http.Request){
	var jsonStr = []byte(`{"content":"Trondheim suger, hilsen golang."}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}