package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"bytes"
	"io/ioutil"
)


func HandlePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload Payload
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if payload.BaseCurrency != "EUR"{
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

	///////////////////FIXER.IO///////////////////
	//	result, err := getFixer(payload.BaseCurrency, payload.TargetCurrency)
	//	if err != nil {
	//		http.Error(w, "Currency not found", http.StatusBadRequest)
	//	}

	//payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	payload.CurrentRate = ReadLatest(payload.TargetCurrency)
	if err != nil{
		http.Error(w, "Currency not found", http.StatusBadRequest)
		return
	}

	db := DatabaseCon()

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

	db := DatabaseCon()

	var payload Payload

	err := db.DB("cloudtech2").C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//	payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	payload.CurrentRate = ReadLatest(payload.TargetCurrency)
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

	db := DatabaseCon()

	err := db.DB("cloudtech2").C("webhooks").RemoveId(bson.ObjectIdHex(s))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	w.WriteHeader(http.StatusOK)
}

/*
func HandleInvoke(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Not a valid ID", http.StatusBadRequest)
		return
	}

	db := DatabaseCon()
	defer db.Close()
	var payload InvokedPayload

	err := db.DB("cloudtech2").C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//payload.CurrentRate, err = getFixer(payload.BaseCurrency, payload.TargetCurrency)
	payload.CurrentRate = ReadLatest(payload.TargetCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "%s", test)
	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(payload)

}
*/

func SendWebhook(url string, data []byte) {
	//var jsonStr= []byte(`{"content":"shit"}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
//	req, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println(err)
		fmt.Println(ioutil.ReadAll(resp.Body))
	}

}

func updateCurrencies(){

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

	for i := 0; i < count; i++{
		//newValue, err := GetFixer(payload[i].BaseCurrency, payload[i].TargetCurrency)
		newValue := ReadLatest(payload[i].TargetCurrency)
		if err != nil{
			fmt.Printf("Error: %s\n", err)
			break
		}
		payload[i].CurrentRate = newValue

		err = c.UpdateId(payload[i].ID, payload[i])
		if err != nil{
			fmt.Printf("Error: %s\n", err)
			break
		}
	//	fmt.Printf("Updated ID: %v\n %s\n", payload[i].ID.Hex(), http.StatusOK)
	}
}