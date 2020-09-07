package backend

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebook_config *oauth2.Config

func getFacebookOauthURL() string {
	options := CreateClientOptions("facebook")

	facebook_config = &oauth2.Config{
		ClientID:     options.getID(),
		ClientSecret: options.getSecret(),
		RedirectURL:  "https://ginoauth-example.herokuapp.com/callback/facebook",
		Scopes: []string{
			"email",
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
	if state != "Facebook" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	code := ctx.Query("code")
	token, err := facebook_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	fmt.Println("token = ", token)
}
