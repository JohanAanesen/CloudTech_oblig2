package main

import (
	"gopkg.in/mgo.v2"
)

var dbURL = "mongodb://johan:123@ds227035.mlab.com:27035/cloudtech2"

func databaseCon()(*mgo.Session){
	session, err := mgo.Dial(dbURL)
	if err != nil{
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}