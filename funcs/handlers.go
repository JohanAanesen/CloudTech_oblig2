package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

//HandleMain main functio for /
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

//HandlePost handles POST requests to main
func HandlePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload Payload
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if payload.BaseCurrency != "EUR" {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

	//payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	payload.CurrentRate = ReadLatest(payload.TargetCurrency)
	if err != nil {
		http.Error(w, "Currency not found", http.StatusBadRequest)
		return
	}

	db := DatabaseCon()

	payload.ID = bson.NewObjectId()

	test := db.DB("cloudtech2").C("webhooks").Insert(&payload)
	if test != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	fmt.Fprintf(w, "%s", payload.ID.Hex())
}

//HandleGet handles GET requests to main
func HandleGet(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := DatabaseCon()

	var payload Payload

	err := db.DB("cloudtech2").C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	//	payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	payload.CurrentRate = ReadLatest(payload.TargetCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

//HandleDelete handles DELETE request to main
func HandleDelete(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := DatabaseCon()

	err := db.DB("cloudtech2").C("webhooks").RemoveId(bson.ObjectIdHex(s))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	w.WriteHeader(http.StatusOK)
}

//HandleLatest handles POST requests to /latest
func HandleLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data LatestPayload

		//json decoder
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}
		if data.BaseCurrency != "EUR" {
			http.Error(w, "Not implemented: %s\n", http.StatusNotImplemented)
			return
		}

		value := ReadLatest(data.TargetCurrency)

		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	} else {
		http.Error(w, "Request method unsupported", http.StatusBadRequest)
	}
}

//HandleAverage handles POST request to /average
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

		if data.BaseCurrency != "EUR" {
			http.Error(w, "Not implemented: %s\n", http.StatusNotImplemented)
			return
		}
		//value := getFixerAverage(current_time, data.BaseCurrency, data.TargetCurrency)
		value := ReadAverage(data.TargetCurrency)
		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	}
}

//HandleEvaluation handles evaluation trigger for evaluation purposes
func HandleEvaluation(w http.ResponseWriter, r *http.Request) {
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

	for i := 0; i < count; i++ {
		//if payload[i].CurrentRate <= payload[i].MinTriggerValue {
			//Send webhook mintrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil {
				fmt.Printf("Json encoding went to shit: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)
	/*	} else if payload[i].CurrentRate >= payload[i].MaxTriggerValue {
			//Send webhook maxtrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil {
				fmt.Printf("Json encoding went to shit: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)

		}*/ /*else{
			//Don't send webhook? dunno
			var jsonStr= []byte(`{"content":"Within margins"}`)
			sendWebhook(url, jsonStr)
			w.WriteHeader(http.StatusOK)
		}*/

	}

	w.WriteHeader(http.StatusOK)

}

//SendWebhook sends the webhook to url with data provided
func SendWebhook(url string, data []byte) {
	//var jsonStr= []byte(`{"content":"shit"}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	//	req, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println(err)
		fmt.Println(ioutil.ReadAll(resp.Body))
	}

}

//UpdateCurrencies updates the currencies to all the registered webhooks and sends webhook if triggered
func UpdateCurrencies() {

	GetFixer("EUR")

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

	for i := 0; i < count; i++ {
		//newValue, err := GetFixer(payload[i].BaseCurrency, payload[i].TargetCurrency)
		newValue := ReadLatest(payload[i].TargetCurrency)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
		payload[i].CurrentRate = newValue

		err = c.UpdateId(payload[i].ID, payload[i])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
		if payload[i].CurrentRate <= payload[i].MinTriggerValue {
			//Send webhook mintrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil {
				fmt.Printf("Json encoding error: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)
		} else if payload[i].CurrentRate >= payload[i].MaxTriggerValue {
			//Send webhook maxtrigger
			var webhookPay InvokedPayload

			webhookPay.BaseCurrency = payload[i].BaseCurrency
			webhookPay.TargetCurrency = payload[i].TargetCurrency
			webhookPay.CurrentRate = payload[i].CurrentRate
			webhookPay.MinTriggerValue = payload[i].MinTriggerValue
			webhookPay.MaxTriggerValue = payload[i].MaxTriggerValue

			b, err := json.Marshal(webhookPay)
			if err != nil {
				fmt.Printf("Json encoding error: %s\n", err)
				return
			}
			SendWebhook(payload[i].WebhookURL, b)

		}
		//	fmt.Printf("Updated ID: %v\n %s\n", payload[i].ID.Hex(), http.StatusOK)
	}
}
