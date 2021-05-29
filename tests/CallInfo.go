package tests

import (
	"fmt"

	"github.com/tofa-project/server-go/calls"
)

func CallInfo(uri string, token string) {
	err := calls.Info(uri, calls.Meta{
		"auth_token":  token,
		"description": "Harmless information notice from test remote application regarding A,S,D,etc",
	})

	fmt.Println(err)
}
