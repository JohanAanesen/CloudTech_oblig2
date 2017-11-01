package main

import (
	"time"
	"github.com/JohanAanesen/CloudTech_oblig2/funcs"
)
func main() {
	for {
		//text := "Heroku timer test at: " + time.Now().String()
		delay := time.Hour * 24

		//sendDiscordLogEntry(text)
		funcs.UpdateCurrencies()

		time.Sleep(delay)
	}
}
