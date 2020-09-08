package backend

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var google_config *oauth2.Config

type googleUser struct {
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

func getGoogleOauthURL() string {
	options := CreateClientOptions("google")

	google_config = &oauth2.Config{
		ClientID:     options.getID(),
		ClientSecret: options.getSecret(),
		RedirectURL:  "https://ginoauth-example.herokuapp.com/callback/google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return google_config.AuthCodeURL("TheWorld")
}

func GoogleOauthLogin(ctx *gin.Context) {
	redirectURL := getGoogleOauthURL()
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

func GoogleCallBack(ctx *gin.Context) {
	state := ctx.Query("state")
	if state != "TheWorld" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
		return
	}

	// 藉由Authorization Code去跟google(resource)申請Access Token
	code := ctx.Query("code")
	token, err := google_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// 藉由獲得的Access Token去跟google申請資源
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

	var user googleUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	redirectURL, _ := url.Parse(IsLoginURL)
	query, _ := url.ParseQuery(redirectURL.RawQuery)
	query.Add("email", user.Email)
	query.Add("name", user.Name)
	redirectURL.RawQuery = query.Encode()
	log.Printf("name = %s, email = %s\n", user.Name, user.Email)
	ctx.Redirect(http.StatusSeeOther, redirectURL.String())
}
