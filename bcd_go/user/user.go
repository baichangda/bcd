package user

import (
	"bcd_go/config"
	"bcd_go/message"
	"bcd_go/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"sync"
)

const CookieTokenName = "bcd-token"

const RedisUserKey = "user"

var SessionMap = sync.Map{}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func add(username string, password string) error {
	users, err := find(username)
	if err != nil {
		return errors.WithStack(err)
	} else {
		if len(users) == 0 {
			marshal, err := json.Marshal(User{
				Username: username,
				Password: password,
			})
			if err != nil {
				return errors.WithStack(err)
			}
			config.RedisClient.HSet(config.RedisCtx, RedisUserKey, username, marshal)
			util.Log.Infof("add succeed")
		} else {
			util.Log.Infof("add failed,user[%s] exist", username)
			return nil
		}
	}
	return nil
}

func del(username string) {
	config.RedisClient.HDel(config.RedisCtx, RedisUserKey, username)
}

func find(username string) ([]User, error) {
	if username == "" {
		all := config.RedisClient.HGetAll(config.RedisCtx, RedisUserKey)
		err := all.Err()
		if err != nil {
			if err == redis.Nil {
				return []User{}, nil
			} else {
				return nil, errors.WithStack(err)
			}
		} else {
			var users []User
			for _, v := range all.Val() {
				user := User{}
				err := json.Unmarshal([]byte(v), &user)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				users = append(users, user)
			}
			return users, nil
		}
	} else {
		get := config.RedisClient.HGet(config.RedisCtx, RedisUserKey, username)
		err := get.Err()
		if err != nil {
			if err == redis.Nil {
				return []User{}, nil
			} else {
				return nil, errors.WithStack(err)
			}
		} else {
			user := User{}
			err := json.Unmarshal([]byte(get.Val()), &user)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return []User{user}, nil
		}
	}

}

func list(_ctx *gin.Context) {
	users, err := find("")
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
	} else {
		message.ResponseSucceed_data(users, _ctx)
	}
}

func login(_ctx *gin.Context) {
	username := _ctx.PostForm("username")
	password := _ctx.PostForm("password")
	users, err := find(username)
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	if len(users) == 0 {
		message.GinError_msg_code(_ctx, "用户不存在", 101)
	} else {
		user := users[0]
		if user.Password == password {
			token := uuid.NewString()
			SessionMap.Store(token, user)
			_ctx.SetCookie(CookieTokenName, token, 86400, "/", "", false, true)
			message.ResponseSucceed_data(user, _ctx)
		} else {
			message.GinError_msg_code(_ctx, "密码错误", 102)
		}
	}
}

func logout(_ctx *gin.Context) {
	token, err := _ctx.Cookie(CookieTokenName)
	if err != nil {
		_ = _ctx.Error(errors.WithStack(err))
		return
	}
	SessionMap.Delete(token)
	_ctx.SetCookie(CookieTokenName, "", 0, "/", "", false, true)
	message.ResponseSucceed_msg("注销成功", _ctx)
}

func CheckLogin(_ctx *gin.Context) bool {
	path := _ctx.FullPath()
	if strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/api/user/login") {
		token, err := _ctx.Cookie(CookieTokenName)
		if err != nil {
			if err == http.ErrNoCookie {
				message.ResponseFailed(_ctx, "请先登陆", 401)
			} else {
				message.ResponseFailed_err(_ctx, err)
			}
			return false
		} else {
			user, ok := SessionMap.Load(token)
			if ok {
				_ctx.Set("user", user)
				return true
			} else {
				message.ResponseFailed(_ctx, "请先登陆", 401)
				return false
			}
		}
	}

	return true

}

func Route(engine *gin.Engine) {
	engine.GET("/api/user/list", list)
	engine.POST("/api/user/login", login)
	engine.POST("/api/user/logout", logout)
}

var list_username string

func listCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "查询用户",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := find("")
			if err != nil {
				util.Log.Errorf("%+v", err)
				return
			}
			if len(users) == 0 {
				util.Log.Info("no user")
			} else {
				builder := strings.Builder{}
				for _, v := range users {
					builder.WriteString("\n")
					marshal, err := json.Marshal(v)
					if err != nil {
						util.Log.Errorf("%+v", err)
						return
					}
					builder.Write(marshal)
				}
				util.Log.Info(builder.String())
			}
		},
	}
	cmd.Flags().StringVarP(&list_username, "username", "u", "", "用户名参数(精确匹配)")
	return &cmd
}

var add_username string
var add_password string

func addCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add",
		Short: "新增用户",
		Run: func(cmd *cobra.Command, args []string) {
			err := add(add_username, add_password)
			if err != nil {
				util.Log.Errorf("%+v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&add_username, "username", "u", "", "用户名")
	cmd.Flags().StringVarP(&add_password, "password", "p", "", "密码")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")
	return &cmd
}

var del_username string

func delCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "del",
		Short: "删除用户",
		Run: func(cmd *cobra.Command, args []string) {
			del(del_username)
		},
	}
	cmd.Flags().StringVarP(&del_username, "username", "u", "", "用户名")
	_ = cmd.MarkFlagRequired("username")
	return &cmd
}

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "user",
		Short: "用户管理",
	}
	cmd.AddCommand(listCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(delCmd())
	return &cmd
}
