package main

import (
	"github.com/davidleitw/gin-oauth2-example/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	// server.GET("/callback", backend.GoogleCallBack)
	// server.GET("/Hello", func(ctx *gin.Context) {
	// 	ctx.String(200, "Hello world")
	// })
	// server.GET("/test", backend.GoogleOauthLogin)

	oauth := server.Group("oauth")
	{
		oauth.GET("/google", backend.GoogleOauthLogin)
		oauth.GET("/facebook", backend.FacebookOauthLogin)
	}
	callback := server.Group("callback")
	{
		callback.GET("/google", backend.GoogleCallBack)
		callback.GET("/facebook", backend.FacebookCallBack)
	}

	_ = server.Run()
}
