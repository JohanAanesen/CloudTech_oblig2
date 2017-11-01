package mongodb

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/JohanAAnesen/CloudTech_oblig2/handlers"
)

var dbURL = "mongodb://johan:123@ds227035.mlab.com:27035/cloudtech2"

func DatabaseCon() *mgo.Session {
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

func SaveData(data handlers.Data){
	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")

	err := c.Insert(data)
	if err != nil{
		fmt.Errorf("something went wrong writing to mongodb: %s", err)
	}
}

func ReadLatest(s string)float64{
	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data handlers.Data

	err :=  c.Find(nil).Skip(dbSize-1).One(&data)
	if err != nil{
		fmt.Errorf("something went wrong reading mongodb: %s", err)
	}

	return data.Rates[s]
}

func ReadAverage(s string)float64{
	db := DatabaseCon()
	defer db.Close()
	c := db.DB("cloudtech2").C("fixer")
	dbSize, _ := c.Count()

	var data1 handlers.Data
	var data2 handlers.Data
	var data3 handlers.Data
	var average float64

	err :=  c.Find(nil).Skip(dbSize-1).One(&data1)
	err =  c.Find(nil).Skip(dbSize-2).One(&data2)
	err =  c.Find(nil).Skip(dbSize-3).One(&data3)
	if err != nil{
		fmt.Errorf("something went wrong reading mongodb: %s", err)
		return 0
	}

	average = (data1.Rates[s] + data2.Rates[s] + data3.Rates[s])/3

	return average
}