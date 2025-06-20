package ovh

import (
	"fmt"
	"log"

	"github.com/dann1/ovh-dynip/ipaddress"
	"github.com/ovh/go-ovh/ovh"
)

type DNSRecord struct {
	ID        int    `json:"id"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	FieldType string `json:"fieldType"`
	TTL       int    `json:"ttl"`
	Zone      string `json:"zone"`
}

var ovhclient *ovh.Client

func init() {
	var err error

	ovhclient, err = ovh.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to initialize OVH client: %v", err)
	}
}

func GetRecordIDs(domain string, subdomain string) ([]int, error) {
	url := fmt.Sprintf("/domain/zone/%s/record?fieldType=A&subDomain=%s", domain, subdomain)

	var recordIDs []int

	err := ovhclient.Get(url, &recordIDs)
	if err != nil {
		return nil, err
	}

	return recordIDs, nil
}

func GetRecord(domain string, recordID int) (*DNSRecord, error) {
	url := fmt.Sprintf("/domain/zone/%s/record/%d", domain, recordID)

	var record DNSRecord

	err := ovhclient.Get(url, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateRecord(domain string, recordID int, ip string) error {
	if !ipaddress.IsValidIPv4(ip) {
		return fmt.Errorf("invalid IPv4 address: %s", ip)
	}

	url := fmt.Sprintf("/domain/zone/%s/record/%d", domain, recordID)
	update := map[string]interface{}{
		"target": ip,
	}

	return ovhclient.Put(url, update, nil)
}

func CreateRecord(domain string, subdomain string, ip string) error {
	if !ipaddress.IsValidIPv4(ip) {
		return fmt.Errorf("invalid IPv4 address: %s", ip)
	}

	url := fmt.Sprintf("/domain/zone/%s/record", domain)
	payload := map[string]interface{}{
		"fieldType": "A",
		"subDomain": subdomain,
		"target":    ip,
		"ttl":       300,
	}

	return ovhclient.Post(url, payload, nil)
}

func RefreshZone(domain string) error {
	url := fmt.Sprintf("/domain/zone/%s/refresh", domain)
	return ovhclient.Post(url, nil, nil)
}

func GenerateConsumerKey(domain string) error {
	req := ovhclient.NewCkRequest()

	rules := []string{
		"/domain/zone/*",
		"/domain/zone",
		fmt.Sprintf("/domain/zone/%s/*", domain),
	}

	for _, rule := range rules {
		req.AddRules(ovh.ReadWrite, rule)
	}

	resp, err := req.Do()
	if err != nil {
		return fmt.Errorf("failed to request consumer key: %w", err)
	}

	fmt.Printf("Visit this URL to authorize: %s\n\nThen set this consumer key in your OVH config:\nconsumer_key=%s", resp.ValidationURL, resp.ConsumerKey)

	return nil
}
