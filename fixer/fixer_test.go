package fixer

import (
	"testing"
	"fmt"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/JohanAAnesen/CloudTech_oblig2/handlers"
	"github.com/JohanAAnesen/CloudTech_oblig2/mongodb"
)

func TestGetFixer(t *testing.T) {

	GetFixer("EUR")

	timeTest := time.Now().Format("2006-01-02")

	db := mongodb.DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data handlers.Data

	err :=  c.Find(nil).Skip(dbSize-1).One(&data)
	if err != nil{
		fmt.Errorf("something went wrong reading mongodb: %s", err)
	}

	testValue := data.Date
	//testValue2 := ReadLatest(out[1])


	if testValue != timeTest{
		t.Fatalf("ERROR expected: %s but got: %s", timeTest, testValue)
	}

	//remove what got added to db
	c.Remove(bson.M{"date": timeTest})

}
/*
func TestGetFixerAverage(t *testing.T) {
	testTime := time.Now()
	var out = []string{"EUR", "NOK"}

	testAverage := GetFixerAverage(testTime, out[0], out[1])

	if testAverage <= 0{
		t.Fatalf("ERROR expected: Integer higher than 0, got: %v", testAverage)
	}
}
*/