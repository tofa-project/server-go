package calls

import (
	"net/http"
	"time"

	tofa_errors "github.com/tofa-project/server-go/errors"
	tor_client "github.com/tofa-project/server-go/tor-client"
)

// Sends a http request and replies the results to dedicated channels
func fireReq(req *http.Request, resChan chan *http.Response, resErrChan chan error) {
	r, e := tor_client.GetClient().Do(req)
	if e == nil {
		resChan <- r
	} else {
		resErrChan <- &tofa_errors.RequestFailed{}
	}
}

// Starts countdown until timeout bacon is sent to timeoutChan
func startCountdown(sec uint, timeoutChan chan bool) {
	time.Sleep(time.Second * time.Duration(sec))

	timeoutChan <- true
}
