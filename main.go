package main

import (
	"vkc/vpp"
	"vkc/wira"

	"github.com/urfave/cli/v2"
)

var (
	NAME    = "SKY"
	Usage   = "GIN"
	VERSION = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Usage = Usage
	app.Version = VERSION
	run := vpp.CreateRunServe(wira.NewApp)
	vpp.RunApp(app, run)
}

// func main() {
// 	r := gin.Default()
// 	r.GET("/ping", ping)
// 	r.Run("0.0.0.0:80")
// }

// // 健康检测
// func ping(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "pong",
// 	})
// }
