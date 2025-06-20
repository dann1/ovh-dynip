package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/dann1/ovh-dynip/ipaddress"
	sovh "github.com/dann1/ovh-dynip/ovh"
)

func main() {
	fqdn := flag.String("update", "", "FQDN to update A record for current public IP")
	generateKey := flag.Bool("generate-key", false, "Generate a consumer key with required permissions")
	publicIP := flag.Bool("public-ip", false, "Get current public IP")

	flag.Parse()

	switch {
	case *fqdn != "":
		updateFQDN(*fqdn)
	case *publicIP:
		fmt.Printf("Public IP: %s\n", getPublicIP())
	case *generateKey:
		err := sovh.GenerateConsumerKey(*fqdn)
		if err != nil {
			log.Fatalf("Failed to generate consumer key: %v", err)
		}
	default:
		flag.Usage()
	}
}

func updateFQDN(fqdn string) {
	domains := strings.Split(fqdn, ".")
	if len(domains) < 3 {
		log.Fatalf("Invalid FQDN: %s", fqdn)
	}

	domain := fmt.Sprintf("%s.%s", domains[1], domains[2])
	subdomain := domains[0]
	publicIP := getPublicIP()

	log.Printf("Updating %s -> %s", fqdn, publicIP)

	recordIDs, err := sovh.GetRecordIDs(domain, subdomain)
	if err != nil {
		log.Fatalf("Could not get record IDs: %v", err)
	}

	if len(recordIDs) == 0 {
		err = sovh.CreateRecord(domain, subdomain, publicIP)
		if err != nil {
			log.Fatalf("Failed to create record: %v", err)
		}
		log.Printf("Created record for %s.%s", subdomain, domain)
	} else {
		recordID := recordIDs[0]
		record, err := sovh.GetRecord(domain, recordID)
		if err != nil {
			log.Fatalf("Could not get record: %v", err)
		}

		if record.Target != publicIP {
			err = sovh.UpdateRecord(domain, recordID, publicIP)
			if err != nil {
				log.Fatalf("Failed to update record: %v", err)
			}
			log.Printf("Updated record %d to %s", recordID, publicIP)
		} else {
			log.Println("No update needed")
		}
	}

	err = sovh.RefreshZone(domain)
	if err != nil {
		log.Fatalf("Failed to refresh zone: %v", err)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Printf("To update the A record with current public IP %s a consumer key with proper permissions is needed:\nGo to https://api.ovh.com/createToken/ to get Application Key and Secret\nThen set them in ~/.ovh.conf as in https://github.com/ovh/go-ovh?tab=readme-ov-file#application-keyapplication-secret\nYou can generate the consumer key with -generate-key", getPublicIP())
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
}

func getPublicIP() (publicIP string) {
	publicIP, err := ipaddress.GetPublicIP()
	if err != nil {
		log.Fatalf("Failed to get public IP: %v", err)
	}

	return publicIP
}
