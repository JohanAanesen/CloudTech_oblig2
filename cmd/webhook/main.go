package main

import "time"

func main() {
	for {
		//text := "Heroku timer test at: " + time.Now().String()
		delay := time.Minute * 15

		//sendDiscordLogEntry(text)
		UpdateCurrencies()

		time.Sleep(delay)
	}
}
