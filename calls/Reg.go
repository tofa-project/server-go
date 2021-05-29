package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	uri_handler "github.com/tofa-project/server-go/uri-handler"
)

// Handles registration with the application.
//
// uri: Tofa client URI
//
// meta: must contain "name" and "description"
//
// Returns authorization key or error according to case
func Reg(uri string, meta Meta) (string, error) {

	// make sure host is reachable
	if err := Ping(uri); err != nil {
		return "", err
	}

	// decode URI
	dec, err := uri_handler.Decode(uri)
	if err != nil {
		return "", err
	}

	// make json
	metaJBytes, err := json.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("Invalid meta! %s", err)
	}

	// make request
	req, err := http.NewRequest("REG", dec.ToHttp(), bytes.NewBuffer(metaJBytes))
	if err != nil {
		return "", fmt.Errorf("Could not create request! %s", err)
	}

	timeoutChan := make(chan bool)
	resChan := make(chan *http.Response)
	resErrChan := make(chan error)

	go fireReq(req, resChan, resErrChan)
	go startCountdown(CALL_RESPONSE_TIMEOUT, timeoutChan)

	// wait for: timeout || response || response error
	select {

	case <-timeoutChan:
		return "", fmt.Errorf("Request timed out!")

	case res := <-resChan:
		defer res.Body.Close()

		// parse based on received code
		switch res.StatusCode {

		case 270:
			resByte, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return "", fmt.Errorf("Could not parse response! %s", err)
			}

			resJson := make(Meta)
			err = json.Unmarshal(resByte, &resJson)
			if err != nil {
				return "", fmt.Errorf("Could not parse response! %s", err)
			}

			return resJson["auth_token"], nil

		case 570:
			return "", fmt.Errorf("Client denied request!")

		default:
			return "", fmt.Errorf("Received code: %d expected 270 || 570", res.StatusCode)
		}

	case err := <-resErrChan:
		return "", err

	}

}
