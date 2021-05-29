package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	uri_handler "github.com/tofa-project/server-go/uri-handler"
)

// Push ask request to Tofa client
//
// uri: Tofa client URI
//
// meta: must contain "auth_token" and "description"
func Ask(uri string, meta Meta) (bool, error) {

	// make sure host is reachable
	if err := Ping(uri); err != nil {
		return false, err
	}

	// decode uri
	dec, err := uri_handler.Decode(uri)
	if err != nil {
		return false, err
	}

	// make json to send
	metaJBytes, err := json.Marshal(Meta{"description": meta["description"]})
	if err != nil {
		return false, fmt.Errorf("Invalid meta! %s", err)
	}

	// crate request
	req, err := http.NewRequest("ASK", dec.ToHttp(), bytes.NewBuffer(metaJBytes))
	if err != nil {
		return false, fmt.Errorf("Could not create request! %s", err)
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
		return false, fmt.Errorf("Request timed out!")

	case res := <-resChan:
		defer res.Body.Close()

		// parse based on received code
		switch res.StatusCode {
		case 270:
			return true, nil

		case 570:
			return false, nil

		default:
			return false, fmt.Errorf("Received code: %d, expected 270 || 570", res.StatusCode)
		}

	case e := <-resErrChan:
		return false, e
	}

}
