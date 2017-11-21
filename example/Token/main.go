package main

import (
	"fmt"
	"os"

	"github.com/spopezen/go-auth0/auth0"
)

func main() {
	c := auth0.NewClient(os.Getenv("AUTH0_TENANT"))
	token, err := c.GetTokenFromClientCreds(os.Getenv("AUTH0_CLIENT_ID"),
		os.Getenv("AUTH0_CLIENT_SECRET"),
		os.Getenv("AUTH0_API"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("TokenResponse:\n%+v\n", token)
}
