package oblig2

import (
	"net/http"
	"strings"
)

//Discord webhook
//var url = "https://discordapp.com/api/webhooks/370393359900082200/_RASdjfNlTsFm9QMprDIFukfV05u7_vfN8nBjgoJ7y0_D_JmLXYdoWVbY8guoCkbOAVx"
//Discord webhook
//var url = "https://discordapp.com/api/webhooks/371707670832349187/dPg6uA7eJL1K0wPxtfyde1ZQu_6LoC_O_SOqrQJ5b_VqcxpfsnGHE4TYKrNz95sAXW3o"
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
