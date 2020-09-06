package backend

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func GetGoogleOauthURL(c *ClientOption) string {
	if c == nil {
		c = CreateClientOptions()
	}

	config := &oauth2.Config{
		ClientID:     c.getID(),
		ClientSecret: c.getSecret(),
		RedirectURL:  "/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return config.AuthCodeURL("TheWorld")
}

func GoogleOauthLogin(ctx *gin.Context) {
	redirectURL := GetGoogleOauthURL(nil)
	log.Println("redirectURL = ", redirectURL)
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

func GoogleCallBack(ctx *gin.Context) {
	log.Println("Call back area. ")
	s := ctx.Query("TheWorld")
	if s != "state" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	code := ctx.Query("code")
	ctx.JSON(200, gin.H{
		"code": code,
	})
}
