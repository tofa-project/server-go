package tests

import (
	"fmt"

	"github.com/tofa-project/server-go/calls"
)

func CallPing(uri string) {
	err := calls.Ping(uri)

	fmt.Println(err)
}
