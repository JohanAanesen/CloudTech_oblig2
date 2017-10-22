package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)
//Data map
type Data struct {
	Base string
	Date string
	Rates map[string]float64
}
//Discord webhook
var url = "https://discordapp.com/api/webhooks/370393359900082200/_RASdjfNlTsFm9QMprDIFukfV05u7_vfN8nBjgoJ7y0_D_JmLXYdoWVbY8guoCkbOAVx"


func getFixer(s1 string, s2 string)(float64, error){

	json1, err := http.Get("http://api.fixer.io/latest?base=" + s1) //+ "," + s2)
	r := json1.Body
	if err != nil{
		fmt.Printf("fixer.io is not responding, %s\n", err)
		return 0, err
	}

	//data object
	var data Data

	//json decoder
	err = json.NewDecoder(r).Decode(&data)

	//err handler
	if err != nil {
		fmt.Printf("shit, %s\n", err)
		return 0, err
	}

	//number := data["rates"][s2].(float64)
	number := data.Rates[s2]
	return number, nil
}

func HandleExchange(w http.ResponseWriter, r *http.Request){
	//content-type set to JSON
	http.Header.Add(w.Header(), "content-type", "application/json")

//	json1, err := http.Get(GitHubURL + URL[4] + "/" + URL[5])
	result, err := getFixer("NOK", "JPY")
	if err != nil{
		http.Error(w, "Currency not found", http.StatusBadRequest)
	}

	fmt.Println(result)
	fmt.Fprint(w, result)

}

func main() {

	http.HandleFunc("/exchange", HandleExchange)
	http.HandleFunc("/discord", sendDiscord)

	//port := os.Getenv("PORT")
	//http.ListenAndServe(":"+port, nil)
	http.ListenAndServe(":8080", nil)
}
