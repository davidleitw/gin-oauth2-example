package backend

import (
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var github_config *oauth2.Config

func getGithubOauthURL() (*oauth2.Config, string) {
	options := CreateClientOptions("github", "https://ginoauth-example.herokuapp.com/callback/google")

	github_config = &oauth2.Config{
		ClientID:     options.getID(),
		ClientSecret: options.getSecret(),
		RedirectURL:  options.getRedirectURL(),
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}

	state := GenerateState()
	return github_config, state
}

func GithubOauthLogin(ctx *gin.Context) {
	config, state := getGithubOauthURL()
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

func GithubCallBack(ctx *gin.Context) {
	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	code := ctx.Query("code")
	token, err := github_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// redirect to islogin page, and add email, name into url's query string.
	redirectURL, err := url.Parse(IsLoginURL)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	query, err := url.ParseQuery(redirectURL.RawQuery)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// query.Add("email", user.Email)
	// query.Add("name", user.Name)
	query.Add("source", "google")
	redirectURL.RawQuery = query.Encode()

	// 跳轉登入成功畫面(顯示登入資訊)
	ctx.Redirect(http.StatusSeeOther, redirectURL.String())
}
