package backend

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebook_config *oauth2.Config

type facebookUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}

func getFacebookOauthURL() string {
	options := CreateClientOptions("facebook")

	facebook_config = &oauth2.Config{
		ClientID:     options.getID(),
		ClientSecret: options.getSecret(),
		RedirectURL:  "https://ginoauth-example.herokuapp.com/callback/facebook",
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}

	return facebook_config.AuthCodeURL("FaceBook")
}

func FacebookOauthLogin(ctx *gin.Context) {
	redirectURL := getFacebookOauthURL()
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

func FacebookCallBack(ctx *gin.Context) {
	if error_reason := ctx.Query("error_reason"); error_reason != "" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(error_reason))
		return
	}

	state := ctx.Query("state")
	if state != "FaceBook" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	code := ctx.Query("code")
	token, err := facebook_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	client := facebook_config.Client(context.TODO(), token)
	fmt.Println("client = ", client)

	userEmail, err := client.Get("https://graph.facebook.com/v8.0/me/messages")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("user info = ", userEmail)
}
