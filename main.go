package main

import (
	"fmt"

	"github.com/davidleitw/gin-oauth2-example/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	s := backend.GetGoogleOauthURL(nil)
	fmt.Println("Server start running")
	fmt.Println("Generate google oauth url: ", s)
	// fmt.Println(s)

	server := gin.Default()
	server.GET("/callback", backend.GoogleCallBack)
	server.GET("/Hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})
	server.GET("/test", backend.GoogleOauthLogin)

	server.Run()

}
