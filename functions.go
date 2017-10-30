package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"bytes"
)

var dbURL = "mongodb://johan:123@ds227035.mlab.com:27035/cloudtech2"

func databaseCon() *mgo.Session {
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

func getFixer(s1 string, s2 string) (float64, error) {

	json1, err := http.Get("http://api.fixer.io/latest?base=" + s1) //+ "," + s2)
	if err != nil {
		fmt.Printf("fixer.io is not responding, %s\n", err)
		return 0, err
	}

	//data object
	var data Data

	//json decoder
	err = json.NewDecoder(json1.Body).Decode(&data)
	if err != nil { //err handler
		fmt.Printf("shit, %s\n", err)
		return 0, err
	}

	//number := data["rates"][s2].(float64)
	number := data.Rates[s2]
	return number, nil
}

func getFixerAverage(t time.Time, s1 string, s2 string) float64 {
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
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload Payload
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	///////////////////FIXER.IO///////////////////
	//	result, err := getFixer(payload.BaseCurrency, payload.TargetCurrency)
	//	if err != nil {
	//		http.Error(w, "Currency not found", http.StatusBadRequest)
	//	}
	payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	if err != nil{
		http.Error(w, "Currency not found", http.StatusBadRequest)
		return
	}

	db := databaseCon()

	payload.ID = bson.NewObjectId()

	test := db.DB("cloudtech2").C("webhooks").Insert(&payload)
	if test != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer db.Close()

	fmt.Fprintf(w, "%s", payload.ID.Hex())
}

func HandleGet(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := databaseCon()

	var payload Payload

	err := db.DB("cloudtech2").C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//	payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func HandleDelete(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := databaseCon()

	err := db.DB("cloudtech2").C("webhooks").RemoveId(bson.ObjectIdHex(s))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
}

func HandleInvoke(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := databaseCon()
	defer db.Close()
	var payload InvokedPayload

	err := db.DB("cloudtech2").C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "%s", test)
	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(payload)

}

func sendWebhook(url string, data []byte) {
	//var jsonStr= []byte(`{"content":"shit"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
//	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func updateCurrencies(w http.ResponseWriter){
	db := databaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("webhooks")
	count, _ := c.Count()

	var payload []Payload

	err := c.Find(nil).All(&payload)
	if err != nil {
		fmt.Printf("It's fucked: %s\n", err)
		return
	}

	for i := 0; i < count; i++{
		newValue, err := getFixer(payload[i].BaseCurrency, payload[i].TargetCurrency)
		if err != nil{
			fmt.Printf("It's fucked: %s\n", err)
			break
		}
		payload[i].CurrentRate = newValue

		err = c.UpdateId(payload[i].ID, payload[i])
		if err != nil{
			fmt.Printf("It's fucked: %s\n", err)
			break
		}
		fmt.Fprintf(w,"Updated ID: %v\n %s\n", payload[i].ID.Hex(), http.StatusOK)
	}
}