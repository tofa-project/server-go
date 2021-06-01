package uri_handler

import (
	"encoding/base64"
	"strconv"
	"strings"

	tofa_errors "github.com/tofa-project/server-go/errors"
)

// decodes uri of type version 0
func decodeVersion0(splUri []string) (*Decoded, error) {
	if len(splUri) != 3 {
		return nil, &tofa_errors.BadURI{"splUri"}
	}

	splPortPath := strings.Split(splUri[2], "/")
	if len(splPortPath) < 2 {
		return nil, &tofa_errors.BadURI{"splPortPath"}
	}

	portInt, err := strconv.Atoi(splPortPath[0])
	if err != nil {
		return nil, &tofa_errors.BadURI{"portInt"}
	}

	return &Decoded{
		Onion:   splUri[1],
		Port:    portInt,
		Path:    splPortPath[1],
		Version: 0,
	}, nil
}

// Decodes an input URI and maps it into a Decoded struct
func Decode(iUri string) (*Decoded, error) {

	// decode b64
	decUriBytes, err := base64.StdEncoding.DecodeString(iUri)
	if err != nil {
		return nil, &tofa_errors.BadURI{"decUriBytes"}
	}

	// conver to string
	decUriStr := string(decUriBytes)

	// split uri string and check if it's valid
	decUriSpl := strings.Split(decUriStr, ":")
	if len(decUriSpl) < 2 {
		return nil, &tofa_errors.BadURI{"decUriSplt"}
	}

	// decode based on version
	switch decUriSpl[0] {
	case "0":
		return decodeVersion0(decUriSpl)
	default:
		return nil, &tofa_errors.UnsupportedURI{"unsupported URI version " + decUriSpl[0]}
	}

}
