package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebook_config *oauth2.Config

type facebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email`
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
	fmt.Println("token = ", token)

	client := facebook_config.Client(context.TODO(), token)
	fmt.Println("client = ", client)

	userInfo, err := client.Get("https://graph.facebook.com/v8.0/me?fields=id,name,email")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	info, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user facebookUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, gin.H{
		"Info": user,
	})
}
