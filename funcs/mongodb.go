package funcs

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var dbURL = "mongodb://johan:123@ds227035.mlab.com:27035/cloudtech2"

//DatabaseCon connects to database and returns session
func DatabaseCon() *mgo.Session {
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

//SaveData stores data in database
func SaveData(data Data) {
	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")

	err := c.Insert(data)
	if err != nil {
		fmt.Printf("something went wrong writing to mongodb: %s", err)
	}
}

//ReadLatest retrieves the latest currencies from database
func ReadLatest(s string) float64 {
	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data Data

	err := c.Find(nil).Skip(dbSize - 1).One(&data)
	if err != nil {
		fmt.Printf("something went wrong reading mongodb: %s", err)
	}

	return data.Rates[s]
}

//ReadAverage retrieves the average of the past 3 days from database
func ReadAverage(s string) float64 {
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
		fmt.Printf("something went wrong reading mongodb: %s", err)
		return 0
	}

	average = (data1.Rates[s] + data2.Rates[s] + data3.Rates[s]) / 3

	return average
}
