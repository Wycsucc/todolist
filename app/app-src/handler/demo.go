package handler

import (
	"github.com/gin-gonic/gin"
)

/*
	Demo ...AllRoute
*/
func Demo(c *gin.Context) {
	c.JSON(200, gin.H{
		"messgae": "demo",
	})
}
