package tests

import (
	"fmt"

	"github.com/tofa-project/server-go/calls"
)

func CallAsk(uri string, token string) {
	accepted, err := calls.Ask(uri, calls.Meta{
		"auth_token":  token,
		"description": "Harmless request from test remote application regarding action Z which leads to A,S,D,etc",
	})

	fmt.Println(accepted, err)
}
