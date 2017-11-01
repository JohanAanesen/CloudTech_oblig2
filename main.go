package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//Discord webhook
//var url = "https://discordapp.com/api/webhooks/370393359900082200/_RASdjfNlTsFm9QMprDIFukfV05u7_vfN8nBjgoJ7y0_D_JmLXYdoWVbY8guoCkbOAVx"
//Discord webhook
var url = "https://discordapp.com/api/webhooks/371707670832349187/dPg6uA7eJL1K0wPxtfyde1ZQu_6LoC_O_SOqrQJ5b_VqcxpfsnGHE4TYKrNz95sAXW3o"
//Slack webhook
//var url = "https://hooks.slack.com/services/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"

//WORKS
func HandleMain(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case "GET":
		HandleGet(URL[1], w, r)
	case "POST":
		HandlePost(w, r)
	case "DELETE":
		HandleDelete(URL[1], w, r)
	default:
		http.Error(w, "Request not supported.", http.StatusNotImplemented)
	}
}

//WORKS
func HandleLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data LatestPayload

		//json decoder
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}
		if data.BaseCurrency != "EUR"{
			http.Error(w, "Not implemented: %s\n", http.StatusNotImplemented)
			return
		}
		/*value, err := getFixer(data.BaseCurrency, data.TargetCurrency)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}*/
		value := ReadLatest(data.TargetCurrency)

		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	}else{
		http.Error(w, "Request method unsupported", http.StatusBadRequest)
	}
}
//WORKS
func HandleAverage(w http.ResponseWriter, r *http.Request) {
	//current_time := time.Now().Local()
	if r.Method == "POST" {
		var data LatestPayload

		//json decoder
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}

		if data.BaseCurrency != "EUR"{
			http.Error(w, "Not implemented: %s\n", http.StatusNotImplemented)
			return
		}
		//value := getFixerAverage(current_time, data.BaseCurrency, data.TargetCurrency)
		value := ReadAverage(data.TargetCurrency)
		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	}
}

func HandleEvaluation(w http.ResponseWriter, r *http.Request){
//	fmt.Fprint(w,"fuck off m8")
//	updateCurrencies(w)

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("webhooks")
	count, _ := c.Count()

	var payload []Payload

	err := c.Find(nil).All(&payload)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for i := 0; i < count; i++{
		if payload[i].CurrentRate <= payload[i].MinTriggerValue{
			//Send webhook mintrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil{
				fmt.Printf("Json encoding went to shit: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)
		}else if payload[i].CurrentRate >= payload[i].MaxTriggerValue{
			//Send webhook maxtrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil{
				fmt.Printf("Json encoding went to shit: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)

		}/*else{
			//Don't send webhook? dunno
			var jsonStr= []byte(`{"content":"Within margins"}`)
			sendWebhook(url, jsonStr)
			w.WriteHeader(http.StatusOK)
		}*/

	}

	w.WriteHeader(http.StatusOK)

}


func main() {
	////////////NEED TO FIGURE THIS ONE OUT////////////
/*	for range time.NewTicker(24 * time.Second).C {
		updateCurrencies()
	}*/

	http.HandleFunc("/", HandleMain)
	http.HandleFunc("/latest", HandleLatest)
	http.HandleFunc("/average", HandleAverage)
	http.HandleFunc("/evaluationtrigger", HandleEvaluation)


	//port := os.Getenv("PORT")
	//http.ListenAndServe(":"+port, nil)
	http.ListenAndServe(":8080", nil)
}
