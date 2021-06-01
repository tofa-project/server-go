package uri_handler

import (
	"strconv"
	"strings"
)

// After decoding a string URI this type will be returned
type Decoded struct {
	// Onion ID to contact amid calls
	Onion string

	// Port to connect to amid calls
	Port int

	// Path to use amid calls
	Path string

	// URI version for appropriate treatment
	Version uint
}

// Converts to http url
func (d *Decoded) ToUrl() string {
	return strings.Join([]string{
		"http://", d.Onion, ".onion:", strconv.Itoa(d.Port), "/", d.Path,
	}, "")
}
