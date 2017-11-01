package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"bytes"
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
	db := DatabaseCon()
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
		newValue, err := GetFixer(payload[i].BaseCurrency, payload[i].TargetCurrency)
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