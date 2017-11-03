package funcs

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestGetFixer(t *testing.T) {

	GetFixer("EUR")

	hour, _, _ := time.Now().Clock()
	timeTest := time.Now()
	if hour < 17 {
		timeTest = timeTest.AddDate(0, 0, -1)
	}
	timeTestString := timeTest.Format("2006-01-02")

	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data Data

	err := c.Find(nil).Skip(dbSize - 1).One(&data)
	if err != nil {
		t.Fatalf("something went wrong reading mongodb: %s", err)
	}

	testValue := data.Date
	//testValue2 := ReadLatest(out[1])

	if testValue != timeTestString {
		t.Fatalf("ERROR expected: %s but got: %s", timeTest, testValue)
	}

	//remove what got added to db
	c.Remove(bson.M{"date": timeTestString})

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
