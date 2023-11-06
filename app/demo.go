package app

import (
	"github.com/gin-gonic/gin"
)

// Demo 接口
type Demo struct {
}

// Register 注册路由
func (aa *Demo) Register(r gin.IRouter) {
	r.GET("hello1", aa.hello)
}

func (aa *Demo) hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "demo hello"})
}
