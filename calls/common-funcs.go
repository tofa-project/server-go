package calls

import (
	"fmt"
	"net/http"
	"time"

	tor_client "github.com/tofa-project/server-go/tor-client"
)

// Sends a http request and replies the results to dedicated channels
func fireReq(req *http.Request, resChan chan *http.Response, resErrChan chan error) {
	r, e := tor_client.GetClient().Do(req)
	if e == nil {
		resChan <- r
	} else {
		resErrChan <- fmt.Errorf("Could not perform request! %s", e)
	}
}

// Starts countdown until timeout bacon is sent to timeoutChan
func startCountdown(sec uint, timeoutChan chan bool) {
	time.Sleep(time.Second * time.Duration(sec))

	timeoutChan <- true
}
