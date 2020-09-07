package main

import (
	"github.com/davidleitw/gin-oauth2-example/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/callback", backend.GoogleCallBack)
	server.GET("/Hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})
	server.GET("/test", backend.GoogleOauthLogin)

	server.Group("auth")
	{
		server.GET("/google", backend.GoogleOauthLogin)
		server.GET("/facebook", backend.FacebookOauthLogin)
	}
	server.Group("callback")
	{
		server.GET("/google", backend.GoogleCallBack)
		server.GET("/facebook", backend.FacebookCallBack)
	}

	_ = server.Run()
}
