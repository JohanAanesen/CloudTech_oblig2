package main

import (
	"time"
	"github.com/JohanAAnesen/CloudTech_oblig2/handlers"
)
func main() {
	for {
		//text := "Heroku timer test at: " + time.Now().String()
		delay := time.Hour * 24

		//sendDiscordLogEntry(text)
		handlers.UpdateCurrencies()

		time.Sleep(delay)
	}
}
