package calls

import (
	"fmt"
	"net/http"

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

	req, err := http.NewRequest("PING", dec.ToHttp(), nil)
	if err != nil {
		return fmt.Errorf("Could not create request! %s", err)
	}

	timeoutChan := make(chan bool)
	resChan := make(chan *http.Response)
	resErrChan := make(chan error)

	go fireReq(req, resChan, resErrChan)
	go startCountdown(CALL_CONNECT_TIMEOUT, timeoutChan)

	select {

	case <-timeoutChan:
		return fmt.Errorf("Request timed out!")

	case res := <-resChan:
		defer res.Body.Close()

		// parse based on received code
		switch res.StatusCode {

		case http.StatusNoContent:
			return nil

		default:
			return fmt.Errorf("Received code: %d, expected %d", res.StatusCode, http.StatusNoContent)

		}

	case e := <-resErrChan:
		return e
	}
}
