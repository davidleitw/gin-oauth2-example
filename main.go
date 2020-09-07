package main

import (
	"github.com/davidleitw/gin-oauth2-example/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// oauth2 group, if you want to login, visit these routes.
	oauth := server.Group("oauth")
	{
		oauth.GET("/google", backend.GoogleOauthLogin)
		oauth.GET("/facebook", backend.FacebookOauthLogin)
	}

	// callback group
	callback := server.Group("callback")
	{
		callback.GET("/google", backend.GoogleCallBack)
		callback.GET("/facebook", backend.FacebookCallBack)
	}

	_ = server.Run()
}
