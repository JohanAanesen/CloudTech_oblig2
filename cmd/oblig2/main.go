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
//	http.HandleFunc("/addlatest", HandleNew)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
//	http.ListenAndServe(":8080", nil)
}

/*
func HandleNew(w http.ResponseWriter, r *http.Request) {

	json1, err := http.Get("http://api.fixer.io/latest") //+ "," + s2)
	if err != nil {
		fmt.Printf("fixer.io is not responding, %s\n", err)
		return
	}

	//data object
	var data funcs.Data

	//json decoder
	err = json.NewDecoder(json1.Body).Decode(&data)
	if err != nil { //err handler
		fmt.Printf("Error: %s\n", err)
		return
	}

	//Storing data in db
	funcs.SaveData(data)

	fmt.Println("ok.")
}*/
