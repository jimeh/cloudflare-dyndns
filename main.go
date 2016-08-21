package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimeh/cloudflare-dyndns/updater"
)

func main() {
	var email string
	var apiKey string
	var host string

	flag.StringVar(
		&email, "email", os.Getenv("CF_EMAIL"),
		"Cloudflare email account. Can be specified via CF_EMAIL env var.",
	)
	flag.StringVar(
		&apiKey, "apikey", os.Getenv("CF_APIKEY"),
		"Cloudflare API key. Can be specified via CF_APIKEY env var.",
	)
	flag.StringVar(
		&host, "host", os.Getenv("CF_HOST"),
		"DNS entry to update. Can be specified via CF_HOST env var.",
	)

	flag.Parse()

	if email == "" || apiKey == "" || host == "" {
		fmt.Println("usage: cloudflare-dyndns <options>\n\noptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	updater := updater.New(email, apiKey)

	stop, err := updater.UpdateLoop(host)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	<-stop // wait forever
}
