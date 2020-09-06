package backend

import (
	"errors"
	"os"
)

var StateError = errors.New("state error.")

type ClientOption struct {
	clientID     string
	clientSecret string
}

func CreateClientOptions() *ClientOption {
	return &ClientOption{
		clientID:     os.Getenv("GoogleID"),
		clientSecret: os.Getenv("GoogleSecret"),
	}
}

func CreateClientOptionsWithString(ID, Secret string) *ClientOption {
	c := new(ClientOption)
	c.setID(ID)
	c.setSecret(Secret)
	return c
}

func (c *ClientOption) setID(ID string) {
	c.clientID = ID
}

func (c *ClientOption) setSecret(Secret string) {
	c.clientSecret = Secret
}

func (c *ClientOption) getID() string {
	return c.clientID
}

func (c *ClientOption) getSecret() string {
	return c.clientSecret
}
