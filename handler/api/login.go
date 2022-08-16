package api

import (
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gopa/config"
	"gopa/gorm"
	"gopa/handler"
	"gopa/model"
	"time"
)

// LoginForm 用户登录表单
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string
	Group    string
	Role     string
}

var identityKey = "username"

// JwtAuth 			api
// @Summary            Login handler
// @Description    Login to obtain token, project and role
// @Tags               auth
// @Accept           application/json
// @Produce          application/json
// @Param                 form            body      LoginForm        true  "form"
// @Success          200        {object}        handler.Response
// @Failure          400        {object}        handler.Response
// @Failure          403        {object}        handler.Response
// @Router                /api/sso-login [post]
func JwtAuth() *jwt.GinJWTMiddleware {
	secretKey := config.GetConfig().Service.JwtSecret
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(secretKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			v, ok := data.(*User)
			if ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					"group": v.Group,
					"role": v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form LoginForm
			var project string
			var role string
			var defaultUser = &User{
				Username: "",
				Group:    "guest",
				Role:     "guest",
			}
			if err := c.BindJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			if err := LdapAuth(form.Username, form.Password); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			var data GopaMembers
			defaultUser.Username = form.Username
			db := gorm.DB.Self
			res := db.Where("username = ?", form.Username).First(&data)
			c.Set("username", form.Username)
			if res.Error != nil {
				fmt.Println("This user does not exist. Now create a new user.")
				req := model.GopaMembers{
					Username: form.Username,
					Project:  "guest",
					Role:     "guest",
				}
				res := db.Create(&req)
				if res.Error != nil {
					return defaultUser, nil
				}
				project = "guest"
				role = "guest"
			} else {
				project = data.Project
				role = data.Role
			}
			c.Set("group", project)
			c.Set("role", role)
			return &User{
				Username: form.Username,
				Group: project,
				Role: role,
			}, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			username, _ := c.Get("username")
			group, _ := c.Get("group")
			role, _ := c.Get("role")
			c.JSON(code, gin.H{
				"token": message,
				"expires": time,
				"username": username,
				"group": group,
				"role": role,
			})
		},
	})

	if err != nil {
		//handler.SendResponse400(context, err, nil)
		return nil
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		//handler.SendResponse400(context, err, nil)
		return nil
	}
	return authMiddleware

}

func Test(context *gin.Context) {
	fmt.Println(jwt.ExtractClaims(context)["username"])
	fmt.Println(jwt.ExtractClaims(context)["group"])
	fmt.Println(jwt.ExtractClaims(context)["role"])
	handler.SendResponse(context, nil, "OK")
}
