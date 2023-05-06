package video

import (
	"bcd_go/config"
	"bcd_go/message"
	"bcd_go/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const videoListRedisKey = "video"

var videoDir = "/Users/baichangda/Downloads/testm3u8/temp"

func Route(engine *gin.Engine) {
	initEnv()
	engine.GET("/api/video/list", list)
	engine.GET("/api/video/downloadM3u8", downloadM3u8)
	engine.GET("/api/video/downloadTs", downloadTs)
}

type Video struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestFfmpeg() {
	err := Ffmpeg("/Users/baichangda/Downloads/testm3u8/test2.mp4", "/api/video/download?path=", "/Users/baichangda/Downloads/testm3u8")
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
}

var isInit = false

func initEnv() {
	if !isInit {
		viper.AutomaticEnv()
		getenv := viper.GetString("VIDEO_DIR")
		util.Log.Infof("get GO_VIDEO_DIR=%s", getenv)
		if len(getenv) > 0 {
			videoDir = getenv
		}
		isInit = true
	}
}

func Ffmpeg(sourcePath string, tsBase string, destDir string) error {
	cmdLine := fmt.Sprintf("-i %s -map 0 -c copy -f segment -segment_format mpegts -segment_time 2 -segment_list %s/output.m3u8 -segment_list_entry_prefix %s %s/output_%%06d.ts", sourcePath, destDir, tsBase, destDir)
	split := strings.Split(cmdLine, " ")
	command := exec.Command("ffmpeg", split...)
	util.Log.Infof("command:\n%s", command.String())
	output, err := command.CombinedOutput()
	util.Log.Infof("output:\n%s", string(output))
	if err != nil {
		return errors.Wrap(err, "CombinedOutput error")
	}
	return nil
}

func DeleteVideo(id string) error {
	//删除redis
	_ = config.RedisClient.HDel(config.RedisCtx, videoListRedisKey, id)
	//清空文件
	dir := path.Join(videoDir, id)
	util.Log.Infof("delete dir[%s]", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func AppendVideoFromDisk(name string, p string) (*Video, error) {
	keys := config.RedisClient.HKeys(config.RedisCtx, videoListRedisKey)
	ids, err := keys.Result()
	var id int64 = 0
	if err != nil && err != redis.Nil {
		if err == redis.Nil {
			for _, e := range ids {
				i, err := strconv.ParseInt(e, 10, 64)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				if i > id {
					id = i
				}
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}
	id += 1

	destDir := path.Join(videoDir, strconv.FormatInt(id, 10))
	err = os.MkdirAll(destDir, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "MkdirAll error")
	}
	tsBase := fmt.Sprintf("/api/video/downloadTs?id=%d&ts=", id)
	err = Ffmpeg(p, tsBase, destDir)
	if err != nil {
		return nil, err
	}

	video := Video{
		Id:   int(id),
		Name: name,
	}
	//存入redis
	marshal, err := json.Marshal(video)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal error")
	}
	config.RedisClient.HSet(config.RedisCtx, videoListRedisKey, id, marshal)
	return &video, nil
}

func downloadM3u8(_ctx *gin.Context) {
	id := _ctx.Query("id")
	_ctx.Set("content-type", "application/x-mpegURL")
	download(path.Join(videoDir, id, "/output.m3u8"), _ctx)
}

func downloadTs(_ctx *gin.Context) {
	id := _ctx.Query("id")
	ts := _ctx.Query("ts")
	p := path.Join(videoDir, id, ts)
	_ctx.Set("content-type", "application/octet-stream")
	download(p, _ctx)
}

func download(path string, _ctx *gin.Context) {
	_ctx.File(path)
}

func list(_ctx *gin.Context) {
	hGetAll := config.RedisClient.HGetAll(config.RedisCtx, videoListRedisKey)
	result, err := hGetAll.Result()
	_ctx.Set("content-type", "application/json;charset=utf-8")
	if err != nil {
		if err == redis.Nil {
			message.ResponseSucceed_data([]Video{}, _ctx)
		} else {
			message.GinError_err(_ctx, err)
		}
	} else {
		if len(result) == 0 {
			message.ResponseSucceed_data([]Video{}, _ctx)
		} else {
			var list []Video
			for _, v := range result {
				cur := Video{}
				err := json.Unmarshal([]byte(v), &cur)
				if err != nil {
					message.GinError_err(_ctx, err)
					return
				} else {
					list = append(list, cur)
				}
			}
			message.ResponseSucceed_data(list, _ctx)
		}
	}
}

func listCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "展示所有视频",
		Run: func(cmd *cobra.Command, args []string) {
			all := config.RedisClient.HGetAll(config.RedisCtx, videoListRedisKey)
			err := all.Err()
			if err != nil {
				if err == redis.Nil {
					util.Log.Info("no video")
					return
				} else {
					util.Log.Errorf("%+v", err)
					return
				}
			}
			if len(all.Val()) == 0 {
				util.Log.Info("no video")
			} else {
				builder := strings.Builder{}
				for _, v := range all.Val() {
					builder.WriteString("\n")
					builder.WriteString(v)
				}
				util.Log.Info(builder.String())
			}
		},
	}
	return &cmd
}

var add_name string
var add_path string

func addCmd() *cobra.Command {
	initEnv()
	cmd := cobra.Command{
		Use:   "add",
		Short: "添加视频",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := AppendVideoFromDisk(add_name, add_path)
			if err != nil {
				util.Log.Errorf("%+v", err)
				return
			}
			util.Log.Infof("add succeed %v", v)
		},
	}
	cmd.Flags().StringVarP(&add_name, "name", "n", "", "视频名称")
	cmd.Flags().StringVarP(&add_path, "path", "p", "", "视频路径")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("path")
	return &cmd
}

var del_id string

func delCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "del",
		Short: "删除视频",
		Run: func(cmd *cobra.Command, args []string) {
			err := DeleteVideo(del_id)
			if err != nil {
				util.Log.Errorf("%+v", err)
				return
			}
			util.Log.Infof("del succeed %s", args[2])
		},
	}
	cmd.Flags().StringVarP(&del_id, "id", "i", "", "id")
	_ = cmd.MarkFlagRequired("id")
	return &cmd
}

func delAllCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delAll",
		Short: "删除所有视频",
		Run: func(cmd *cobra.Command, args []string) {
			config.RedisClient.Del(config.RedisCtx, videoListRedisKey)
			err := os.RemoveAll(videoDir)
			if err != nil {
				util.Log.Errorf("%+v", err)
				return
			}
			util.Log.Infof("delAll succeed")
		},
	}
	return &cmd
}

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "video",
		Short: "视频",
	}
	cmd.AddCommand(listCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(delCmd())
	cmd.AddCommand(delAllCmd())
	return &cmd
}
