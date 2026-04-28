package router

import (
	"os"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/middleware"

	// docs/swagger 包由 `swag init -o docs/swagger` 生成。通过匿名 import
	// 触发 init()，将 swagger spec 注册到 swaggo 运行时中。
	_ "github.com/QuantumNous/new-api/docs/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetSwaggerRouter 挂载 Swagger UI。
//
// 访问：GET /swagger/index.html
//
// 门控策略：
//   - 非 release 模式默认开启；
//   - release 模式默认关闭，可通过环境变量 ENABLE_SWAGGER=true 强制开启（生产调试用）。
func SetSwaggerRouter(router *gin.Engine) {
	if !swaggerEnabled() {
		return
	}

	swaggerGroup := router.Group("/swagger")
	swaggerGroup.Use(middleware.RouteTag("swagger"))
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	common.SysLog("Swagger UI enabled at /swagger/index.html")
}

func swaggerEnabled() bool {
	if strings.EqualFold(os.Getenv("ENABLE_SWAGGER"), "true") {
		return true
	}
	if strings.EqualFold(os.Getenv("ENABLE_SWAGGER"), "false") {
		return false
	}
	return gin.Mode() != gin.ReleaseMode
}
