## Golang Adapter
Use this package if you're developing web apps in Go.

## Include
`go get github.com/tofa-project/server-go@e7030394e4b80b14024c08440ebabe6a0ee320d0`

## Use
Flow of use:
1. Initialize with Tor socks5 proxy address
2. Use adapter methods to communicate with Tofa Clients

Calls are synchronous.

```go
package main

import (
  "fmt"

  tor_client "github.com/tofa-project/server-go/tor-client"
  
  "github.com/tofa-project/server-go/calls"
  "github.com/tofa-project/server-go/errors"
)

func main(){
  /**
  * First initialize adapter with Tor proxy address. Usually it's 127.0.0.1:9050
  * Adapter instance is shared across application, so init should be called only once.
  */
  tor_client.InitClient("127.0.0.1:9050")
  
  
  /**
  * Attempts to register with Tofa Client. 
  * It requires Client URI, and metadata so human can recognize your service.
  * Metadata must contain "name" and "description" (both strings).
  * 
  * @returns: the authentication token which is mandatory when performing ASK and INFO calls
  *           Also an error, if it occurred.
  *
  * Registration process must occur only once, and authentication token
  * must be stored in a database and re-used for eternity.
  */
  auth_token, err := calls.Reg(uri, calls.Meta{
    "name":        "Test remote app",
    "description": "Harmless test remote app description",
  })

  fmt.Println(auth_token, err)
  

  /**
  * Attempts to ask for confirmation form Tofa Client amid an action. 
  * It requires Client URI, and metadata so human can recognize the action.
  * Metadata must contain a comprehensive "description" and the "auth_token" (both strings).
  * 
  * @returns: true/false whether human allowed the action or not and error if it occurred
  */
  accepted, err := calls.Ask(uri, calls.Meta{
    "auth_token":  auth_token,
    "description": "Harmless request from test remote application regarding action Z ",
  })

  fmt.Println(accepted, err)
  
  /**
  * Attempts to send an INFO call. This is only a notification sent to the Client.
  * It requires Client URI, and metadata so human can recognize your service.
  * Metadata must contain "name" and "auth_token" (both strings).
  * 
  * @returns: error, if it occurred
  */
  err := calls.Info(uri, calls.Meta{
    "auth_token":  auth_token,
    "description": "Harmless information notice from test remote application regarding A,S,D,etc",
  })

  fmt.Println(err)
  
  /**
  * Errors are splitted based on error case.
  * You can take actions based on which error occurred using type assertion.
  *
  * A full documented list can be browsed within IDE at github.com/tofa-project/server-go/errors
  */
  if inst, ok := err.(errors.BadCall); ok {
    // do stuff
  }
}
```


