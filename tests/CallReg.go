package tests

import (
	"fmt"

	"github.com/tofa-project/server-go/calls"
)

func CallReg(uri string) string {
	auth_token, err := calls.Reg(uri, calls.Meta{
		"name":        "Test remote app",
		"description": "Harmless test remote app description",
	})

	fmt.Println(auth_token, err)

	return auth_token
}
