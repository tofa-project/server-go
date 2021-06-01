package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	tofa_errors "github.com/tofa-project/server-go/errors"
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
		return fmt.Errorf("invalid meta! %s", err)
	}

	// make request
	req, err := http.NewRequest("INFO", dec.ToUrl(), bytes.NewBuffer(metaJBytes))
	if err != nil {
		return fmt.Errorf("could not create request! %s", err)
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
		return &tofa_errors.CallTimedOut{}

	case res := <-resChan:
		defer res.Body.Close()

		if res.StatusCode == 200 {
			return nil
		} else {
			return tofa_errors.GetErrorByCode(res.StatusCode)
		}

	case e := <-resErrChan:
		return e
	}
}
