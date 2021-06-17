package infobloxclient

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// InfobloxClient represents Infoblox API client
type InfobloxClient struct {
	*resty.Client
}

// NewInfobloxClient instantiates REST client for Infoblox
func NewInfobloxClient(url string, version string, username string, password string, insecure bool, debug bool) InfobloxClient {
	client := resty.New()
	if debug {
		client.SetDebug(true)
	}
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecure})
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	infobloxURL := fmt.Sprintf("%s/wapi/%s/", url, version)
	client.SetHostURL(infobloxURL)
	client.SetBasicAuth(username, password)
	client.SetQueryParam("_return_as_object", "1")
	return InfobloxClient{client}
}
