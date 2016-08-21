package updater

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// DefaultIPCheckURL is the default URL used to figure out the public IP.
const DefaultIPCheckURL = "http://whatismyip.akamai.com/"

// DefaultInterval is the default number of seconds to wait before each update.
const DefaultInterval = 30

// New creates a new DynDNS instance.
func New(email string, apiKey string) *Updater {
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		log.Fatal(err)
	}

	return &Updater{
		API:        api,
		IPCheckURL: DefaultIPCheckURL,
		Interval:   DefaultInterval,
	}
}

// Updater deals with updating the IP address for a DNS record
type Updater struct {
	API        *cloudflare.API
	IPCheckURL string
	Interval   int
}

// UpdateLoop performs a the full update sequence.
func (u *Updater) UpdateLoop(host string) (<-chan bool, error) {
	stop := make(chan bool, 1)

	fmt.Printf("Looking up record for %s...\n", host)
	record, err := u.RecordByHost(host)
	if err != nil {
		return stop, err
	}

	fmt.Printf("Found %s (Zone: %s)\n", record.Name, record.ZoneName)
	go func() {
		fmt.Printf("Starting IP check (repeats every %d seconds)\n", u.Interval)

		_, err = u.UpdateRecord(record)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}

		for {
			select {
			case <-stop:
				close(stop)
				return
			case <-time.After(time.Second * time.Duration(u.Interval)):
				_, err = u.UpdateRecord(record)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				}
			}
		}
	}()

	return stop, nil
}

// Update performs a single update
func (u *Updater) Update(host string) error {
	record, err := u.RecordByHost(host)
	if err != nil {
		return err
	}

	_, err = u.UpdateRecord(record)
	return err
}

// UpdateRecord updates a cloudflare.DNSRecord.
func (u *Updater) UpdateRecord(record *cloudflare.DNSRecord) (*cloudflare.DNSRecord, error) {
	currentIP, err := u.WhatIsMyIP()
	if err != nil {
		return nil, err
	}

	record, err = u.Record(record.ZoneID, record.ID)
	if err != nil {
		return nil, err
	}

	if currentIP != record.Content {
		fmt.Printf(
			"Updating %s to %s (was %s)\n",
			record.Name, currentIP, record.Content,
		)
		record.Content = currentIP
		err = u.API.UpdateDNSRecord(record.ZoneID, record.ID, *record)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}

// Record fetches a cloudflare.DNSRecord from the given Zone and Record IDs.
func (u *Updater) Record(zoneID string, recordID string) (*cloudflare.DNSRecord, error) {
	record, err := u.API.DNSRecord(zoneID, recordID)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// RecordByHost fetches a cloudflare.DNSRecord from the host given.
func (u *Updater) RecordByHost(host string) (*cloudflare.DNSRecord, error) {
	zoneID, err := u.ZoneID(host)
	if err != nil {
		return nil, err
	}

	recordID, err := u.RecordID(host, zoneID)
	if err != nil {
		return nil, err
	}

	record, err := u.Record(zoneID, recordID)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// WhatIsMyIP fetches the public IP via http://whatismyip.akamai.com/
func (u *Updater) WhatIsMyIP() (string, error) {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(u.IPCheckURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(
			"Got a %d response from %s",
			resp.StatusCode, u.IPCheckURL,
		)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ZoneID finds the zone ID for the relevant Host.
func (u *Updater) ZoneID(host string) (string, error) {
	zones, err := u.API.ListZones()
	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if strings.HasSuffix(host, zone.Name) {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("No zone found for \"%s\"", host)
}

// RecordID finds the host's DNS record ID.
func (u *Updater) RecordID(host string, zoneID string) (string, error) {
	records, err := u.API.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return "", err
	}

	for _, r := range records {
		if r.Type == "A" && r.Name == host {
			return r.ID, nil
		}
	}

	return "", fmt.Errorf("No A type record found for \"%s\"", host)
}
