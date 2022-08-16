package util

import (
	ctx "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"gopa/gorm"
	"gopa/model"
	"path/filepath"
	"strings"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

// BuildQuery 构建OPA请求的query参数
// 如: data.perf_server.api.v1.bus.latestData
// 该方法会把-替换成_
func BuildQuery(referer string, general bool) string {
	query := "data"
	var querySuffix string
	if general {
		querySuffix = "any.allow"
	} else {
		querySuffix = ".allow"
	}
	s1 := strings.Split(referer, "?")[0]
	s2 := strings.Replace(s1, "-", "_", -1)
	query = query + strings.Replace(s2, "/", ".", -1) + querySuffix
	return query
}


// BuildPath 根据referer参数
// 构建OPA请求的mongo路径，用于查询mongo中的策略文件
// 如：/perf-server/api/v1/bus/latestData.rego
func BuildPath(referer string) string {
	s1 := strings.Split(referer, "?")[0]
	mongoPath := s1 + ".rego"
	return mongoPath
}


// BuildForm 根据请求头内的参数
// 构建请求OPA权限接口的表单
// 如果没有对应参数，该方法返回错误信息
func BuildForm(context *gin.Context, params []string) map[string]interface{} {
	input := map[string]interface{}{}
	for _, argument := range params {
		arg := context.GetHeader(argument)
		input[argument] = arg
	}
	return input
}


// ArgumentsParser 根据用户添加的rego策略文件内容
// 解析出input字段名并返回
func ArgumentsParser(fileContent string) []string {
	segments1 := strings.Split(fileContent, "\n")
	var arguments []string
	for _, segment := range segments1 {
		segment = strings.TrimSpace(segment)
		if strings.HasPrefix(segment, "input") {
			segment2 := strings.Replace(strings.Split(segment, " ")[0], "input.", "", -1)
			if !contains(arguments, segment2) {
				arguments = append(arguments, segment2)
			}
		}
	}
	return arguments
}

// contains 判断elem在不在list里
func contains(list []string, elem string) bool {
	for _, element := range list {
		if element == elem {
			return true
		}
	}
	return false
}


// FetchRegoByPath 根据path返回Mongo数据库中的rego文件
func FetchRegoByPath(path string) (model.RegoDocument, error) {
	var result model.RegoDocument
	filter := bson.D{{"path", path}}
	//print("请求mongo路径：", path)
	collection := gorm.Collections.RegoCollection
	err := collection.FindOne(ctx.TODO(), filter).Decode(&result)
	return result, err
}

// GeneralMatch 从referer的最高项开始向下查询，确认该用户
// 是否有任意路径下的通配权限，如果有返回true，无错误
// 如果该方法没有找到通配权限，返回false，无错误
// 如果在鉴权过程中出错，返回false和错误信息
func GeneralMatch(context *gin.Context, context1 ctx.Context, path string) (bool, error) {
	//
	dir := filepath.Dir(path)
	var root = ""
	var evalQuery string
	dirSegments := strings.Split(dir, "/")
	for _, dirSegment := range dirSegments {
		root = root + dirSegment + "/"
		evalQuery = BuildQuery(root, true)
		document, err := FetchRegoByPath(root + "any.rego")
		if err != nil {
			continue
		}
		module := document.Content
		arguments := document.Arguments
		form := BuildForm(context, arguments)
		query := rego.New(
			rego.Query(evalQuery),
			rego.Module(root, module),
			rego.Input(form),
		)
		results, err := query.Eval(context1)
		if err != nil {
			msg := fmt.Sprintf("Error when evaluating path: %s for general match.", root)
			return false, errors.New(msg)
		}
		if results.Allowed() {
			//fmt.Printf("在 %s 路径下有通配权限", root)
			return true, nil
		}
	}
	fmt.Printf("路径 [%s] 下没有通配权限\n", path)
	return false, nil
}