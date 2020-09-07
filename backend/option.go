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

func createClientOptions(company string) *ClientOption {
	var ID, Secret string
	switch company {
	case "google":
		ID = os.Getenv("GoogleID")
		Secret = os.Getenv("GoogleSecret")
	case "facebook":
		ID = os.Getenv("FacebookID")
		Secret = os.Getenv("FacebookSecret")
	case "github":
		ID = os.Getenv("GithubID")
		Secret = os.Getenv("GithubSecret")
	case "twitter":
		ID = os.Getenv("TwitterID")
		Secret = os.Getenv("TwitterSecret")
	default:
		ID = ""
		Secret = ""
	}

	return &ClientOption{
		clientID:     ID,
		clientSecret: Secret,
	}
}

func CreateClientOptions(company string) *ClientOption {
	return createClientOptions(company)
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
