package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/rs/xid"
)

var StateError = errors.New("state error.")

const IsLoginURL = "/islogin" // "https://ginoauth-example.herokuapp.com/Islogin"

type ClientOption struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func createClientOptions(company, redirectURL string) *ClientOption {
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
		redirectURL:  redirectURL,
	}
}

func CreateClientOptions(company string, redirectURL string) *ClientOption {
	return createClientOptions(company, redirectURL)
}

func CreateClientOptionsWithString(ID, Secret, RedirectURL string) *ClientOption {
	c := new(ClientOption)
	c.setID(ID)
	c.setSecret(Secret)
	c.setRedirectURL(RedirectURL)
	return c
}

func (c *ClientOption) setRedirectURL(URL string) {
	c.redirectURL = URL
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

func (c *ClientOption) getRedirectURL() string {
	return c.redirectURL
}

func GenerateState() string {
	return xid.New().String()
}

func __debug__printJSON(js []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, js, "", "\n")

	result := string(prettyJSON.Bytes())

	if err == nil {
		log.Println(result)
	} else {
		log.Println("Println Json error = ", err)
	}
}
