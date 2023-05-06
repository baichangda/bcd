package main

import (
	"bcd_go/message"
	"bcd_go/photo"
	"bcd_go/tools/ocr"
	"bcd_go/user"
	"bcd_go/video"
	"bcd_go/wechat"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"strings"
)

func startHttpServer() {
	g := gin.Default()
	g.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(func(_ctx *gin.Context) {
		//response content-type类型处理
		fullPath := _ctx.FullPath()
		if !strings.HasPrefix(fullPath, "/api/photo/download") && !strings.HasPrefix(fullPath, "/api/video") {
			_ctx.Header("content-type", "application/json")
		}

		//错误捕获
		defer func() {
			errs := _ctx.Errors
			if len(errs) > 0 {
				message.FromGinError(errs[0]).Response(_ctx)
			}
		}()

		_ctx.Next()
	})

	user.Route(g)
	wechat.Route(g)
	video.Route(g)
	photo.Route(g)
	ocr.Route(g)

	err := g.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

var (
	port string
)

func ServerCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "server",
		Short: "启动web服务",
		Run: func(cmd *cobra.Command, args []string) {
			startHttpServer()
		},
	}
	cmd.Flags().StringVarP(&port, "port", "p", "80", "http服务端口")
	return &cmd
}
