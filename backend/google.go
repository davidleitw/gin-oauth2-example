package backend

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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

func GetGoogleOauthURL(c *ClientOption) string {
	if c == nil {
		c = CreateClientOptions()
	}

	config := &oauth2.Config{
		ClientID:     c.getID(),
		ClientSecret: c.getSecret(),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return config.AuthCodeURL("state")
}

func GoogleOauthLogin(ctx *gin.Context) {
	redirectURL := GetGoogleOauthURL(nil)
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}
