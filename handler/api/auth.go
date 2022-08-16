package api

import (
	ctx "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
	"gopa/handler"
	"gopa/schema"
	"gopa/util"
	"strings"
)

// Auth api
// @Summary            Auth
// @Description    Determine if a user has permission on some resource.
// @Tags               auth
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project    header            string    true  "Project"
// @Param            role       header            string    true  "Role"
// @Success          200                {object}          handler.Response
// @Failure          400                {object}          handler.Response
// @Failure          403                {object}          handler.Response
// @Router                    /api/v1/auth [get]
func Auth(context *gin.Context) {
	var forbiddenError = errors.New(schema.AccessForbiddenMessage)
	context1 := ctx.Background()
	referer := context.GetHeader("referer")
	if referer == "" {
		handler.SendResponse400(context, errors.New("no referer"), nil)
		return
	}
	referer = strings.TrimSuffix(referer, "/")
	// 暂时对swagger文档一律放行
	if strings.Contains("swagger", referer) {
		handler.SendResponse(context, nil, "allowed")
		return
	}
	mongoPath := util.BuildPath(referer)
	evalQuery := util.BuildQuery(referer, false)
	// 拿通配单独匹配判断
	fmt.Println("寻找通配规则")
	fmt.Println()
	ok, err := util.GeneralMatch(context, context1, mongoPath)
	if ok {
		handler.SendResponse(context, nil, "general match")
		return
	}
	// 在运行通配权限时出错
	if err != nil {
		handler.SendResponse403(context, err, nil)
		return
	}
	regoFileResult, err := util.FetchRegoByPath(mongoPath)
	// 没有通配权限
	// 也没有找到该路径对应的策略文件
	if err != nil {
		msg := fmt.Sprintf("No rules specified for %s", mongoPath)
		handler.SendResponse403(context, errors.New(msg), nil)
		return
	}
	module := regoFileResult.Content
	arguments := regoFileResult.Arguments
	form := util.BuildForm(context, arguments)
	// Evaluate Permission
	fmt.Println("Evaluating")
	query := rego.New(
		rego.Query(evalQuery),
		rego.Module(mongoPath, module),
		rego.Input(form),
	)
	results, err := query.Eval(context1)
	if err != nil {
		handler.SendResponse400(context, err, nil)
		return
	}
	if results.Allowed() {
		handler.SendResponse(context, nil, "allowed")
	} else {
		handler.SendResponse403(context, forbiddenError, "rejected")
	}
}
