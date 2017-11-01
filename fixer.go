package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func GetFixer(s1 string){//, s2 string)(float64, error) {

	json1, err := http.Get("http://api.fixer.io/latest?base=" + s1) //+ "," + s2)
	if err != nil {
		fmt.Printf("fixer.io is not responding, %s\n", err)
		return
	}

	//data object
	var data Data

	//json decoder
	err = json.NewDecoder(json1.Body).Decode(&data)
	if err != nil { //err handler
		fmt.Printf("shit, %s\n", err)
		return
	}

	//Storing data in db
	SaveData(data)

//	return data.Rates[s2], nil
}
/*
func GetFixerAverage(t time.Time, s1 string, s2 string) float64 {
	var total float64
	//creates copy of time
	timeCopy := t
	//loops through 7 iterations
	for i := t.Day(); i > t.Day()-7; i-- {
		json1, err := http.Get("http://api.fixer.io/" + timeCopy.Format("2006-01-02") + "?base=" + s1)
		//err handler
		if err != nil {
			fmt.Printf("fixer.io is not responding, %s\n", err)
			return 0
		}
		//sets timecopy to yesterday
		timeCopy = timeCopy.AddDate(0, 0, -1)

		//data object
		var data Data

		//json decoder
		err = json.NewDecoder(json1.Body).Decode(&data)
		//err handler
		if err != nil {
			fmt.Printf("Something went wrong decoding json from fixer.io: %s\n", err)
			return 0
		}
		total += data.Rates[s2]
	}

	return total / 7
}*/