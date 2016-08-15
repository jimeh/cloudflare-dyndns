package main

import (
	"log"
	"os"

	"github.com/jimeh/cloudflare-dyndns/updater"
)

func main() {
	var email = os.Getenv("CF_EMAIL")
	var apiKey = os.Getenv("CF_API")
	var host = os.Getenv("CF_HOST")

	updater := updater.New(email, apiKey)

	err := updater.Update(host)
	if err != nil {
		log.Fatal(err)
	}
}
