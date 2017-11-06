package funcs

import (
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"
)

func TestSaveData(t *testing.T) {
	//var testData Data

	testData := Data{
		Base:  "TEST",
		Date:  time.Now().Format("2006-01-02"),
		Rates: map[string]float64{"NOK": 1337, "USD": 69},
	}

	SaveData(testData)

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	//	dbSize, _ := c.Count()

	var testData2 Data
	c.Find(bson.M{"base": "TEST"}).One(&testData2)

	if testData2.Rates["NOK"] != 1337 {
		t.Fatalf("ERROR expected: %v but got: %v", testData.Rates["NOK"], testData2.Rates["NOK"])
	} else if testData2.Base != "TEST" {
		t.Fatalf("ERROR expected: %s but got: %s", testData.Base, testData2.Base)
	}

	c.Remove(bson.M{"base": "TEST"})

}

func TestReadLatest(t *testing.T) {
	out := "NOK"
	testValue := ReadLatest(out)

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data Data

	err := c.Find(nil).Skip(dbSize - 1).One(&data)
	if err != nil {
		t.Fatalf("something went wrong reading mongodb: %s", err)
	}

	if data.Rates[out] != testValue {
		t.Fatalf("ERROR expected: %v but got: %v", data.Rates[out], testValue)
	}
}

func TestReadAverage(t *testing.T) {
	out := "NOK"
	testAverage := ReadAverage(out)

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data1 Data
	var data2 Data
	var data3 Data
	var average float64

	err := c.Find(nil).Skip(dbSize - 1).One(&data1)
	err = c.Find(nil).Skip(dbSize - 2).One(&data2)
	err = c.Find(nil).Skip(dbSize - 3).One(&data3)
	if err != nil {
		t.Fatalf("something went wrong reading mongodb: %s", err)
		os.Exit(1)
	}

	average = (data1.Rates[out] + data2.Rates[out] + data3.Rates[out]) / 3

	if testAverage != average {
		t.Fatalf("ERROR expected: %v but got: %v", average, testAverage)
	}
}
