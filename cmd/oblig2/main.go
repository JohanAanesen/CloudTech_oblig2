package oblig2

import (
	"net/http"
	"github.com/JohanAanesen/CloudTech_oblig2/funcs"
	"github.com/JohanAanesen/CloudTech_oblig2/handlers"
)

//Discord webhook
//var url = "https://discordapp.com/api/webhooks/370393359900082200/_RASdjfNlTsFm9QMprDIFukfV05u7_vfN8nBjgoJ7y0_D_JmLXYdoWVbY8guoCkbOAVx"
//Discord webhook
//var url = "https://discordapp.com/api/webhooks/371707670832349187/dPg6uA7eJL1K0wPxtfyde1ZQu_6LoC_O_SOqrQJ5b_VqcxpfsnGHE4TYKrNz95sAXW3o"
//Slack webhook
//var url = "https://hooks.slack.com/services/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"

func main() {
	////////////NEED TO FIGURE THIS ONE OUT////////////
/*	for range time.NewTicker(24 * time.Second).C {
		updateCurrencies()
	}*/

	http.HandleFunc("/", funcs.HandleMain)
	http.HandleFunc("/latest", funcs.HandleLatest)
	http.HandleFunc("/average", funcs.HandleAverage)
	http.HandleFunc("/evaluationtrigger", funcs.HandleEvaluation)


	//port := os.Getenv("PORT")
	//http.ListenAndServe(":"+port, nil)
	http.ListenAndServe(":8080", nil)
}
