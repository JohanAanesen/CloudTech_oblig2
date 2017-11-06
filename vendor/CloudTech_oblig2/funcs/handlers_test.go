package funcs

import (
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandleDelete(t *testing.T) {

	type testLoad struct {
		ID   bson.ObjectId `json:"id" bson:"_id"`
		Base string        `json:"base" bson:"base"`
	}

	var payload testLoad
	payload.ID = bson.NewObjectId()
	payload.Base = "TEST"

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("webhooks")
	//dbSize, _ := c.Count()

	c.Insert(&payload)

	req, err := http.NewRequest("DELETE", "/"+payload.ID.Hex(), nil)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleMain)

	handler.ServeHTTP(httpTest, req)

	if status := httpTest.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandleGet(t *testing.T) {
	var testLoad Payload

	//known webhook data
	testID := "59edb79bbde1ab3fe0bd101f"
	testBase := "EUR"
	testTarget := "NOK"

	//sends known webhook id
	req, err := http.NewRequest("GET", "/"+testID, nil)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleMain)

	handler.ServeHTTP(httpTest, req)

	err = json.NewDecoder(httpTest.Body).Decode(&testLoad)
	if err != nil { //err handler
		t.Errorf("Something wrong with json decoder: %s", err)
	}

	if testLoad.BaseCurrency != testBase {
		t.Fatalf("ERROR expected: %s but got: %s", testBase, testLoad.BaseCurrency)
	}
	if testLoad.TargetCurrency != testTarget {
		t.Fatalf("ERROR expected: %s but got: %s", testTarget, testLoad.TargetCurrency)
	}

}

func TestHandleLatest(t *testing.T) {
	var testPayload LatestPayload

	testPayload.BaseCurrency = "EUR"
	testPayload.TargetCurrency = "NOK"

	testValue := ReadLatest("NOK")

	json1, _ := json.Marshal(testPayload)
	reader := bytes.NewReader(json1)

	req, err := http.NewRequest("POST", "/latest", reader)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleLatest)

	handler.ServeHTTP(httpTest, req)

	string1 := httpTest.Body.String()
	responseValue, _ := strconv.ParseFloat(string1, 64)

	if testValue != responseValue {
		t.Fatalf("ERROR expected: %v but got: %v", testValue, responseValue)
	}
}

func TestHandleAverage(t *testing.T) {
	var testPayload LatestPayload

	testPayload.BaseCurrency = "EUR"
	testPayload.TargetCurrency = "NOK"

	testValue := ReadAverage("NOK")

	json1, _ := json.Marshal(testPayload)
	reader := bytes.NewReader(json1)

	req, err := http.NewRequest("POST", "/average", reader)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleAverage)

	handler.ServeHTTP(httpTest, req)

	string1 := httpTest.Body.String()
	responseValue, _ := strconv.ParseFloat(string1, 64)

	if testValue != responseValue {
		t.Fatalf("ERROR expected: %v but got: %v", testValue, responseValue)
	}
}

func TestHandlePost(t *testing.T) {
	var testPay Payload

	testPay.WebhookURL = "TEST"
	testPay.BaseCurrency = "EUR"
	testPay.TargetCurrency = "NOK"
	testPay.MinTriggerValue = 1
	testPay.MaxTriggerValue = 1337

	json1, _ := json.Marshal(testPay)
	reader := bytes.NewReader(json1)

	req, err := http.NewRequest("POST", "/", reader)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)

	handler.ServeHTTP(httpTest, req)

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("webhooks")
	dbSize, _ := c.Count()

	var testData Payload
	c.Find(nil).Skip(dbSize - 1).One(&testData)

	if testPay.WebhookURL != testData.WebhookURL {
		t.Errorf("ERROR: got %v want %v", testData.WebhookURL, testPay.WebhookURL)
	}
	if testPay.TargetCurrency != testData.TargetCurrency {
		t.Errorf("ERROR: got %v want %v", testData.TargetCurrency, testPay.TargetCurrency)
	}
	if testPay.BaseCurrency != testData.BaseCurrency {
		t.Errorf("ERROR: got %v want %v", testData.BaseCurrency, testPay.BaseCurrency)
	}
	if testPay.MaxTriggerValue != testData.MaxTriggerValue {
		t.Errorf("ERROR: got %v want %v", testData.WebhookURL, testPay.WebhookURL)
	}

	c.Remove(bson.M{"webhookURL": "TEST"})

}

func TestHandlePost2(t *testing.T) {
	var testPay Payload

	testPay.WebhookURL = "test.url"
	testPay.BaseCurrency = "JPY"
	testPay.TargetCurrency = "NOK"
	testPay.MinTriggerValue = 1
	testPay.MaxTriggerValue = 1337

	json1, _ := json.Marshal(testPay)
	reader := bytes.NewReader(json1)

	req, err := http.NewRequest("POST", "/", reader)
	if err != nil {
		t.Fatal(err)
	}

	httpTest := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)

	handler.ServeHTTP(httpTest, req)

	if status := httpTest.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

}
