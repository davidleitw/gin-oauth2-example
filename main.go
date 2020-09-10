package main

import (
	"github.com/davidleitw/gin-oauth2-example/backend"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	server.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})
	server.Static("/img", "./frontend/img")
	server.Static("login", "./frontend/login")
	server.Static("islogin", "./frontend/Islogin")

	// oauth2 group, if you want to login, visit these routes.
	oauth := server.Group("oauth")
	{
		oauth.GET("/google", backend.GoogleOauthLogin)
		oauth.GET("/facebook", backend.FacebookOauthLogin)
		oauth.GET("/github", backend.GithubOauthLogin)
	}

	// callback group
	callback := server.Group("callback")
	{
		callback.GET("/google", backend.GoogleCallBack)
		callback.GET("/facebook", backend.FacebookCallBack)
		callback.GET("/github", backend.GithubCallBack)
	}

	_ = server.Run()
}
