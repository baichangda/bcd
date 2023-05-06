package photo

import (
	"bcd_go/message"
	"bcd_go/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"sort"
	"time"
)

var dir = "/Users/baichangda"

var isInit = false

func initEnv() {
	if !isInit {
		viper.AutomaticEnv()
		getenv := viper.GetString("PHOTO_DIR")
		util.Log.Infof("get PHOTO_DIR=%s", getenv)
		if len(getenv) > 0 {
			dir = getenv
		}
		isInit = true
	}
}

func Route(engine *gin.Engine) {
	initEnv()
	engine.POST("/api/photo/upload", upload)
	engine.GET("/api/photo/del", del)
	engine.GET("/api/photo/list", list)
	engine.Static("/api/photo/download", dir)
}

func del(_ctx *gin.Context) {
	name := _ctx.Query("name")
	if len(name) == 0 {
		_ = _ctx.Error(message.NewMyError("删除失败、必须有name参数", 1))
		return
	}
	err := os.Remove(dir + "/" + name)
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	message.ResponseSucceed_msg("删除成功", _ctx)
}

func upload(_ctx *gin.Context) {
	file, err := _ctx.FormFile("file")
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	err = _ctx.SaveUploadedFile(file, dir+"/"+time.Now().Format("20060102150405.000"))
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	message.ResponseSucceed_msg("上传成功", _ctx)
}

func list(_ctx *gin.Context) {
	readDir, err := os.ReadDir(dir)
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	var photos []string
	for _, f := range readDir {
		if !f.IsDir() {
			photos = append(photos, f.Name())
		}
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i] > photos[j]
	})

	message.ResponseSucceed_data(photos, _ctx)
}
