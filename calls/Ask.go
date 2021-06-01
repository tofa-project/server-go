package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	tofa_errors "github.com/tofa-project/server-go/errors"
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
		return false, fmt.Errorf("invalid meta! %s", err)
	}

	// crate request
	req, err := http.NewRequest("ASK", dec.ToUrl(), bytes.NewBuffer(metaJBytes))
	if err != nil {
		return false, fmt.Errorf("could not create request! %s", err)
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
		return false, &tofa_errors.CallTimedOut{}

	case res := <-resChan:
		defer res.Body.Close()

		if res.StatusCode == 270 {
			return true, nil
		} else {
			return false, tofa_errors.GetErrorByCode(res.StatusCode)
		}

	case e := <-resErrChan:
		return false, e
	}
}
