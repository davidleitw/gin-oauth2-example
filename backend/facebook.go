package backend

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebook_config *oauth2.Config

type facebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getFacebookOauthURL() (*oauth2.Config, string) {
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
	state := GenerateState()
	return facebook_config, state
}

func FacebookOauthLogin(ctx *gin.Context) {
	config, state := getFacebookOauthURL()
	redirectURL := config.AuthCodeURL(state)

	session := sessions.Default(ctx)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

func FacebookCallBack(ctx *gin.Context) {
	if error_reason := ctx.Query("error_reason"); error_reason != "" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(error_reason))
		return
	}

	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
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

	// redirect to islogin page, and add email, name into url's query string.
	redirectURL, _ := url.Parse(IsLoginURL)
	query, _ := url.ParseQuery(redirectURL.RawQuery)
	query.Add("email", user.Email)
	query.Add("name", user.Name)
	query.Add("source", "facebook")
	redirectURL.RawQuery = query.Encode()
	ctx.Redirect(http.StatusSeeOther, redirectURL.String())
}
