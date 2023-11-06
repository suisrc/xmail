package app

import (
	"github.com/gin-gonic/gin"
)

// Index 接口
type Index struct {
}

// Register 注册路由
func (aa *Index) Register(e *gin.Engine) {
	// e.LoadHTMLGlob("static/templates/*")
	e.StaticFile("/favicon.ico", "static/favicon.ico")
	e.Static("/ui", "static/ui")
}
