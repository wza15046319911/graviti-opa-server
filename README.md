# gopa-server
GOPA（Graviti OPA）：Graviti内部权限管理中台。该中台定义了组和角色的概念，规定了某个组的某个角色可以使用某些资源。资源可以理解为是一个具体的项目，也可以理解为该项目的请求路由。

举个🌰：我是`infra-cloud`组的`admin`，我可以使用eva发布平台发布一款应用。我们假设eva发布的api为`/api/v1/release`，则其他人无法访问该API（无法发布）。

该中台基于golang的gin web框架，使用[gin-jwt](https://github.com/appleboy/gin-jwt)中间件处理登录以及生成token的逻辑，使用[Open Policy Agent](https://github.com/open-policy-agent/opa)控制访问权限，并接入kong网关。

### Swagger user guide
项目使用swagger管理API
* 安装 `swagger`
```
   $ sudo mkdir -p $GOPATH/src/github.com/swaggo
   $ cd $GOPATH/src/github.com/swaggo
   $ git clone https://github.com/swaggo/swag
   $ cd swag/cmd/swag/
   $ go install -v
```
如果不知道gopath在哪里，或者gopath不在`PATH`内，可以运行：
```
   $ go env
```
该指令会显示gopath
* 下载 `gin-swagger`
```
   $ cd $GOPATH/src/github.com/swaggo
   $ git clone https://github.com/swaggo/gin-swagger
```
* 生成swagger文档
```
   $ cd xxx/pathto/gopa-server/
   $ swag init
```
* API注释示例
```
   // @Summary       api
   // @Description   Add a new user
   // @Tags          user
   // @Accept        json
   // @Produce       json
   // @Param         env   path     string          true "dev/fat/uat/pro"
   // @Param         user  body     model.UserInfo  true "Create a new user"
   // @Success       200   {object} handler.Response 
   // @Router        /user/{env} [post]
   func Create(c *gin.Context) {
       ...
   }
```


详情见[swagger文档](https://github.com/swaggo/swag/blob/master/README.md)

### How to start

* git clone ssh://git@phabricator.graviti.cn:2224/source/gopa-server.git
* swag init
* go build .
* go run main.go