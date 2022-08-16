package middleware

import (
	ctx "context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
	"gopa/config"
	"gopa/handler"
	"gopa/util"
)

type Decision struct {
	DecisionID string `json:"decision_id"`
	Result     struct {
		Add struct {
			Allow bool `json:"allow"`
		} `json:"add"`
	} `json:"result"`
}

// GetPermission OPA权限中间件。
// @Summary      请求API时，中间件会返回该用户是否有对应权限。
// 如果具备相应权限，才会继续接下来的请求，否则终止请求。
// 每次执行完权限判断后，程序会删除指定目录下的rego文件
// @Description  GetPermission
// @Tags         middleware
func GetPermission() gin.HandlerFunc {
	return func(context *gin.Context) {
		context1 := ctx.Background()
		referer := context.Request.RequestURI
		mongoPath := util.BuildPath(referer)
		evalQuery := util.BuildQuery(referer, false)
		// 拿通配单独匹配判断
		ok, err := util.GeneralMatch(context, context1, mongoPath)
		if ok {
			context.Next()
			return
		}
		regoFileResult, err := util.FetchRegoByPath(mongoPath)
		if err != nil {
			msg := fmt.Sprintf("no rules specified for [%s]", referer)
			handler.SendResponse403(context, errors.New(msg), nil)
			context.Abort()
			return
		}
		module := regoFileResult.Content
		arguments := regoFileResult.Arguments
		form := util.BuildForm(context, arguments)
		if err != nil {
			handler.SendResponse403(context, err, nil)
			context.Abort()
			return
		}
		query := rego.New(
			rego.Query(evalQuery),
			rego.Module(mongoPath, module),
			rego.Input(form),
		)
		results, err := query.Eval(context1)
		if err != nil {
			handler.SendResponse403(context, err, nil)
			context.Abort()
			return
		}
		if results.Allowed() {
			context.Next()
		} else {
			handler.SendResponse403(context, errors.New("rejected"), nil)
			context.Abort()
			return
		}
	}
}


// Watch 监听对应文件夹的文件内容变化，并作出
// 相应响应
func Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file: ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error: ", err)
			}
		}
	}()
	err = watcher.Add(config.GetConfig().Opa.WatchDirectory)
	if err != nil {
		return
	}
	<-done
}

