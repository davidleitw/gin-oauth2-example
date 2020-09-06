package backend

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var StateError = errors.New("state error.")
var google_config *oauth2.Config

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

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

	google_config = &oauth2.Config{
		ClientID:     c.getID(),
		ClientSecret: c.getSecret(),
		RedirectURL:  "https://ginoauth-example.herokuapp.com/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return google_config.AuthCodeURL("TheWorld")
}

func GoogleOauthLogin(ctx *gin.Context) {
	redirectURL := GetGoogleOauthURL(nil)
	log.Println("redirectURL = ", redirectURL)
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

func GoogleCallBack(ctx *gin.Context) {
	state := ctx.Query("state")
	if state != "TheWorld" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	// use code to get access token
	code := ctx.Query("code")
	token, err := google_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	client := google_config.Client(context.TODO(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer userInfo.Body.Close()

	info, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user User
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(200, gin.H{
		"Info": user,
	})
}
