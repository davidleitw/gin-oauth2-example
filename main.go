package main

import (
	"example/backend"
	"fmt"
)

func main() {
	s := backend.GetGoogleOauthURL(nil)
	fmt.Println(s)
}
