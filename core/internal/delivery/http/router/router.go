package router

import (
	"github.com/gin-gonic/gin"
)

func Init() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())

	return router, nil
}
