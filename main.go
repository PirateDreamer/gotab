package main

import (
	"net/http"

	"github.com/PirateDreamer/gotab/logic"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	// 创建默认配置
// 	option := &logic.Option{}
// 	cli.Execute(option)
// 	logic.HandlerOption(option)
// 	logic.Run(option)
// }

func main() {
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,                                                // 允许所有来源的请求
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},   // 允许的方法
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
	}))

	r.POST("/sql_proto", func(c *gin.Context) {
		type Request struct {
			Sql string `json:"sql" binding:"required"`
		}
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, Response{
				Code: 1,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		protoStr, err := logic.SqlToProto(req.Sql, false)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Code: 1,
				Msg:  err.Error(),
				Data: nil,
			})
		}
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  err.Error(),
			Data: protoStr,
		})
	})
	r.Run()
}
