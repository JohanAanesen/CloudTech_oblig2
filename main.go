package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)



//Discord webhook
//var url = "https://discordapp.com/api/webhooks/370393359900082200/_RASdjfNlTsFm9QMprDIFukfV05u7_vfN8nBjgoJ7y0_D_JmLXYdoWVbY8guoCkbOAVx"
//Discord webhook
var url = "https://discordapp.com/api/webhooks/371707670832349187/dPg6uA7eJL1K0wPxtfyde1ZQu_6LoC_O_SOqrQJ5b_VqcxpfsnGHE4TYKrNz95sAXW3o"
//Slack webhook
//var url = "https://hooks.slack.com/services/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"


func HandleDiscord(w http.ResponseWriter, r *http.Request) {
	type data1 struct {
		Content string `json:"content"`
	}
	var data data1
	data.Content = "jørg1 er k00l"

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "sheeet %s\n", err)
	}

	sendWebhook(url, b)

}

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
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
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
		value, err := getFixer(data.BaseCurrency, data.TargetCurrency)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}

		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	}else{
		http.Error(w, "Request method unsupported", http.StatusBadRequest)
	}
}
//WORKS
func HandleAverage(w http.ResponseWriter, r *http.Request) {
	current_time := time.Now().Local()
	if r.Method == "POST" {
		var data LatestPayload

		//json decoder
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Somethings wrong: %s\n", http.StatusBadRequest)
			return
		}
		value := getFixerAverage(current_time, data.BaseCurrency, data.TargetCurrency)

		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprintf(w, "%v", value)
	}
}

func HandleEvaluation(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w,"fuck off m8")
}

func main() {
	////////////NEEDS ITS OWN WORKER////////////
/*	for range time.NewTicker(24 * time.Second).C {
		updateCurrencies()
	}*/

	http.HandleFunc("/discord", HandleDiscord)
	http.HandleFunc("/", HandleMain)
	http.HandleFunc("/latest", HandleLatest)
	http.HandleFunc("/average", HandleAverage)
	http.HandleFunc("/evaluationtriggger", HandleEvaluation)



	//port := os.Getenv("PORT")
	//http.ListenAndServe(":"+port, nil)
	http.ListenAndServe(":8080", nil)
}
