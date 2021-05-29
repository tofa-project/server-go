package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	uri_handler "github.com/tofa-project/server-go/uri-handler"
)

// Push information to Tofa client
//
// uri: Tofa client URI
//
// meta: must contain "auth_token" and "description"
func Info(uri string, meta Meta) error {
	// make sure host is reachable
	if err := Ping(uri); err != nil {
		return err
	}

	// decode URI
	dec, err := uri_handler.Decode(uri)
	if err != nil {
		return err
	}

	// make json
	metaJBytes, err := json.Marshal(Meta{"description": meta["description"]})
	if err != nil {
		return fmt.Errorf("Invalid meta! %s", err)
	}

	// make request
	req, err := http.NewRequest("INFO", dec.ToHttp(), bytes.NewBuffer(metaJBytes))
	if err != nil {
		return fmt.Errorf("Could not create request! %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+meta["auth_token"])

	timeoutChan := make(chan bool)
	resChan := make(chan *http.Response)
	resErrChan := make(chan error)

	go fireReq(req, resChan, resErrChan)
	go startCountdown(CALL_CONNECT_TIMEOUT, timeoutChan)

	// wait for timeout || res || resError
	select {

	case <-timeoutChan:
		return fmt.Errorf("Request timed out!")

	case res := <-resChan:
		defer res.Body.Close()

		// parse based on received code
		switch res.StatusCode {
		case 200:
			return nil

		default:
			return fmt.Errorf("Received code: %d, expected 200", res.StatusCode)
		}

	case e := <-resErrChan:
		return e
	}
}
