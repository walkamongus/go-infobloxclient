package infobloxclient

import (
	"fmt"
	"log"
	"strings"
)

// InfobloxRecordResults represents one or more record results returned from API
type InfobloxRecordResults struct {
	Result []InfobloxRecord `json:"result"`
}

// InfobloxRecordResult represents a single record result returned from API
type InfobloxRecordResult struct {
	Result InfobloxRecord `json:"result"`
}

// InfobloxRecord represents generic data for CNAME, TXT, and Host records
// Some fields may be unpopulated depending on the record type
type InfobloxRecord struct {
	Ref             string     `json:"_ref,omitempty"`
	Canonical       string     `json:"canonical,omitempty"`
	Comment         string     `json:"comment,omitempty"`
	ConfigureForDNS bool       `json:"configure_for_dns,omitempty"`
	Name            string     `json:"name"`
	View            string     `json:"view,omitempty"`
	Text            string     `json:"text,omitempty"`
	Ipv4addrs       []Ipv4addr `json:"ipv4addrs,omitempty"`
	Ipv4addr        string     `json:"ipv4addr,omitempty"`
}

// Ipv4addr represents IPv4 settings for a single IP address
type Ipv4addr struct {
	Ref              string `json:"_ref,omitempty"`
	ConfigureForDhcp bool   `json:"configure_for_dhcp,omitempty"`
	Host             string `json:"host,omitempty"`
	Ipv4addr         string `json:"ipv4addr"`
}

// CreateRecord creates an Infoblox record of a specified type
func (c *InfobloxClient) CreateRecord(recordType string, data InfobloxRecord) (*InfobloxRecordResult, error) {
	rtype := strings.ToLower(recordType)
	response, err := c.R().
		SetBody(data).
		SetQueryParam("_return_fields+", "comment").
		SetResult(&InfobloxRecordResult{}).
		Post(fmt.Sprintf("/record:%s", rtype))
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case 401:
		return new(InfobloxRecordResult), fmt.Errorf("Unauthorized request. Infoblox API returned: %s", response.String())
	case 201:
		record := response.Result().(*InfobloxRecordResult)
		return record, nil
	default:
		return nil, fmt.Errorf("%s", response.String())
	}
}

// GetRecord retrieves a Infoblox record of a specified type
func (c *InfobloxClient) GetRecord(recordType string, name string) (*InfobloxRecordResults, error) {
	rtype := strings.ToLower(recordType)
	search := "name"
	if rtype == "cname" {
		search = "canonical"
	}
	response, err := c.R().
		SetResult(&InfobloxRecordResults{}).
		SetQueryParam("_return_fields+", "comment").
		Get(fmt.Sprintf("/record:%s?%s=%s", rtype, search, name))
	if err != nil {
		log.Fatal(err)
	}

	switch response.StatusCode() {
	case 401:
		return new(InfobloxRecordResults), fmt.Errorf("Unauthorized request. Infoblox API returned: %s", response.String())
	case 200:
		// length check is two here due to a successful response
		// that finds no records containing an "empty" slice with
		// two bytes in it
		if len(response.Body()) > 2 {
			record := response.Result().(*InfobloxRecordResults)
			return record, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("%s", response.String())
	}
}

// UpdateRecord creates an Infoblox record of a specified type
func (c *InfobloxClient) UpdateRecord(ref string, data InfobloxRecord) (*InfobloxRecordResult, error) {
	response, err := c.R().
		SetBody(data).
		SetResult(&InfobloxRecordResult{}).
		SetQueryParam("_return_fields+", "comment").
		Put(fmt.Sprintf("/%s", ref))
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case 401:
		return new(InfobloxRecordResult), fmt.Errorf("Unauthorized request. Infoblox API returned: %s", response.String())
	case 200:
		// length check is two here due to a successful response
		// that finds no records containing an "empty" slice with
		// two bytes in it
		if len(response.Body()) > 2 {
			record := response.Result().(*InfobloxRecordResult)
			return record, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("%s", response.String())
	}
}

// DeleteRecord deletes a Infoblox record by its ref
func (c *InfobloxClient) DeleteRecord(ref string) error {
	response, err := c.R().Delete(fmt.Sprintf("/%s", ref))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("%s", response.String())
	}
	return nil
}
