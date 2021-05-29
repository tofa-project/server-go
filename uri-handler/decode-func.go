package uri_handler

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// decodes uri of type version 0
func decodeVersion0(splUri []string) (*Decoded, error) {
	if len(splUri) != 3 {
		return nil, fmt.Errorf("Invalid URI!")
	}

	splPortPath := strings.Split(splUri[2], "/")
	if len(splPortPath) < 2 {
		return nil, fmt.Errorf("Invalid URI!")
	}

	portInt, err := strconv.Atoi(splPortPath[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid URI! %s", err)
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
		return nil, fmt.Errorf("Invalid URI! %s", err)
	}

	// conver to string
	decUriStr := string(decUriBytes)

	// split uri string and check if it's valid
	decUriSpl := strings.Split(decUriStr, ":")
	if len(decUriSpl) < 2 {
		return nil, fmt.Errorf("Invalid URI!")
	}

	// decode based on version
	switch decUriSpl[0] {
	case "0":
		return decodeVersion0(decUriSpl)
	default:
		return nil, fmt.Errorf("Unsupported URI version!")
	}

}
