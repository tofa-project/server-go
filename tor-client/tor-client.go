package tor_client

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/proxy"
)

// Client instance will be reused when possible
var client *http.Client

// Avoid race!
var cMux = sync.Mutex{}

// To avoid race
func getClientMux() *http.Client {
	cMux.Lock()
	defer cMux.Unlock()

	if client == nil {
		panic("Client not initialized! Use tor_client.InitClient(addr)")
	}

	return client
}

// To avoid race
func setClientMux(c *http.Client) {
	cMux.Lock()
	defer cMux.Unlock()

	client = c
}

// Initializes client if not initialized already with
// custom SOCKS5 address to connect to
func InitClient(addr string) error {
	if client != nil {
		return fmt.Errorf("client already started")
	}

	// make tcp dialer
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		return err
	}

	// make transporter
	tr := &http.Transport{Dial: dialer.Dial}

	// make client
	setClientMux(&http.Client{
		Transport: tr,
	})

	return nil
}

// Returns the HTTP client used to perform requests via Tor
func GetClient() *http.Client {
	return getClientMux()
}
