package calls

import (
	"fmt"
	"net/http"

	tofa_errors "github.com/tofa-project/server-go/errors"
	uri_handler "github.com/tofa-project/server-go/uri-handler"
)

// Ping to Tofa client
//
// uri: Tofa client URI
func Ping(uri string) error {

	dec, err := uri_handler.Decode(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PING", dec.ToUrl(), nil)
	if err != nil {
		return fmt.Errorf("could not create request! %s", err)
	}

	timeoutChan := make(chan bool)
	resChan := make(chan *http.Response)
	resErrChan := make(chan error)

	go fireReq(req, resChan, resErrChan)
	go startCountdown(CALL_CONNECT_TIMEOUT, timeoutChan)

	select {

	case <-timeoutChan:
		return &tofa_errors.ConnectTimedOut{}

	case res := <-resChan:
		defer res.Body.Close()

		if res.StatusCode == http.StatusNoContent {
			return nil
		} else {
			return tofa_errors.GetErrorByCode(res.StatusCode)
		}

	case e := <-resErrChan:
		return e
	}
}
