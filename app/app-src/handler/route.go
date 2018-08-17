package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
	此处添加路由
*/
func AllRoute(router *gin.Engine, prefix string) {
	// demo
	demo := fmt.Sprintf("%s/%s", prefix, "demo")
	router.GET(demo, Demo)
}
