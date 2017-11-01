package main

import (
	"testing"
	"time"
	"fmt"
)

func TestGetFixer(t *testing.T) {
	var out = []string{"EUR", "NOK"}
//	testValue := GetFixer(out[0], out[1])

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data Data

	err :=  c.Find(nil).Skip(dbSize-1).One(&data)
	if err != nil{
		fmt.Errorf("something went wrong reading mongodb: %s", err)
	}

	testValue := data.Rates[out[1]]
	testValue2 := ReadLatest(out[1])


	if testValue != testValue2{
		t.Fatalf("ERROR expected: %s but got: %s", testValue2, testValue)
	}

}

func TestGetFixerAverage(t *testing.T) {
	testTime := time.Now()
	var out = []string{"EUR", "NOK"}

	testAverage := GetFixerAverage(testTime, out[0], out[1])

	if testAverage <= 0{
		t.Fatalf("ERROR expected: Integer higher than 0, got: %v", testAverage)
	}
}
