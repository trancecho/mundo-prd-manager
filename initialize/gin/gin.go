package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-prd-manager/config"
	"github.com/trancecho/mundo-prd-manager/initialize"
	"github.com/trancecho/mundo-prd-manager/initialize/router"
	"github.com/trancecho/mundo-prd-manager/server/api"
	"github.com/trancecho/mundo-prd-manager/server/middleware"
)

func GinInit() *gin.Engine {
	r := gin.Default()
	config.ConfigInit()
	initialize.DBInit()
	api.ClientInit()
	middleware.InitSecret()
	router.GenerateRouter(r)
	return r
}
