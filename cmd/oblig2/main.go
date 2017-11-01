package main

import (
	"github.com/JohanAanesen/CloudTech_oblig2/funcs"
	"net/http"
	"os"
)

//Discord webhook
//var url = "https://discordapp.com/api/webhooks/371707670832349187/dPg6uA7eJL1K0wPxtfyde1ZQu_6LoC_O_SOqrQJ5b_VqcxpfsnGHE4TYKrNz95sAXW3o"

func main() {

	http.HandleFunc("/", funcs.HandleMain)
	http.HandleFunc("/latest", funcs.HandleLatest)
	http.HandleFunc("/average", funcs.HandleAverage)
	http.HandleFunc("/evaluationtrigger", funcs.HandleEvaluation)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
	//http.ListenAndServe(":8080", nil)
}
