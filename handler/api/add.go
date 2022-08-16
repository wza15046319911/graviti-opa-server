package api

import (
	ctx "context"
	"github.com/gin-gonic/gin"
	"gopa/gorm"
	h "gopa/handler"
	"gopa/model"
	"gopa/util"
)

// ProjectAddForm 某用户发起添加新的组的请求时的结构体
type ProjectAddForm struct {
	AddedGroupName string `json:"project_name"`
}

// RoleAddForm 某用户发起添加新的角色的请求时的结构体
type RoleAddForm struct {
	AddedRoleName string `json:"role_name"`
}

// ProjectRolesAddForm 某用户发起向某个组添加新角色时发起请求的结构体
type ProjectRolesAddForm struct {
	ToProject     string `json:"to_project"`
	AddedRoleName string `json:"role_name"`
}

// UserAddForm 某用户发起添加新用户的请求时的结构体
type UserAddForm struct {
	Username string `json:"username"`
	Project  string `json:"project"`
	Role     string `json:"role"`
}

type ApplicationAddForm struct {
	AddedApplicationName string `json:"application_name"`
	ApplicationAddress   string `json:"application_address"`
	ApplicationRouters   string `json:"application_routers"`
}

// RegoForm 接口RegoAdd接受的表单数据
type RegoForm struct {
	Method      string `json:"method"`
	FilePath    string `json:"path"`
	Filename    string `json:"name"`
	FileContent string `json:"content"`
}

type ProjectResourceForm struct {
	ResourceName string `json:"resource_name"`
	ProjectName  string `json:"project_name"`
	RoleName     string `json:"role_name"`
}

// ProjectAdd 		api
// @Summary          ProjectAdd
// @Description    Add a project/group into the database
// @Tags               project
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project  header            string                      true  "Project"
// @Param            role     header            string                      true  "Role"
// @Param                     form              body        ProjectAddForm        true    "form"
// @Success          200              {object}          handler.Response
// @Failure          400              {object}          handler.Response
// @Router                    /api/v1/project/add [post]
func ProjectAdd(context *gin.Context) {
	var form ProjectAddForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	addedGroupName := form.AddedGroupName
	group := model.GopaProjects{ProjectName: addedGroupName}
	db := gorm.DB.Self
	res := db.Create(&group)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// RoleAdd 			api
// @Summary          RoleAdd
// @Description    Add a role into the database
// @Tags               role
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project        header            string                   true  "Project"
// @Param            role           header            string                   true  "Role"
// @Param                     form          body                RoleAddForm    true    "form"
// @Success          200                    {object}          handler.Response
// @Failure          400                    {object}          handler.Response
// @Router                    /api/v1/role/add [post]
func RoleAdd(context *gin.Context) {
	var form RoleAddForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	addedRoleName := form.AddedRoleName
	role := model.GopaRoles{RoleName: addedRoleName}
	db := gorm.DB.Self
	res := db.Create(&role)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// ProjectRoleAdd 	api
// @Summary          ProjectRoleAdd
// @Description    Add a role to the given project into the database
// @Tags               projectRole
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project        header            string                                 true  "Project"
// @Param            role           header            string                                 true  "Role"
// @Param                     form          body                ProjectRolesAddForm  true    "form"
// @Success          200                    {object}          handler.Response
// @Failure          400                    {object}          handler.Response
// @Router                    /api/v1/projectRole/add [post]
func ProjectRoleAdd(context *gin.Context) {
	var form ProjectRolesAddForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	req := model.GopaProjectRoles{
		ProjectName: form.ToProject,
		RoleName:    form.AddedRoleName,
	}
	db := gorm.DB.Self
	res := db.Create(&req)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// UserAdd 			api
// @Summary          UserAdd
// @Description    Add a user into the database
// @Tags               user
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project        header            string      true         "Project"
// @Param            role           header            string      true         "Role"
// @Param                     form                    body        UserAddForm  true    "form"
// @Success          200                    {object}          handler.Response
// @Failure          400                    {object}          handler.Response
// @Router                    /api/v1/user/add [post]
func UserAdd(context *gin.Context) {
	var form UserAddForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	req := model.GopaMembers{
		Username: form.Username,
		Project:  form.Project,
		Role:     form.Role,
	}
	db := gorm.DB.Self
	res := db.Create(&req)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// ApplicationAdd 	api
// @Summary          ApplicationAdd
// @Description    Add an application into the database
// @Tags               application
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project  header            string                            true  "Project"
// @Param            role     header            string                            true  "Role"
// @Param                     form              body        ApplicationAddForm    true    "form"
// @Success          200              {object}          handler.Response
// @Failure          400              {object}          handler.Response
// @Router                    /api/v1/application/add [post]
func ApplicationAdd(context *gin.Context) {
	var form ApplicationAddForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	req := model.GopaApplication{
		ResourceName: form.AddedApplicationName,
		Address:      form.ApplicationAddress,
		Routers:      form.ApplicationRouters,
	}
	db := gorm.DB.Self
	res := db.Create(&req)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// RegoAdd api
// @Summary          RegoAdd
// @Description    Add a new rego file by given path to the database. If the path already exists, do not insert.
// @Tags               rego
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project  header            string      true      "Project"
// @Param            role     header            string      true      "Role"
// @Param                     form              body        RegoForm    true    "form"
// @Success          200              {object}          handler.Response
// @Failure          400              {object}          handler.Response
// @Router                    /api/v1/rego/add [post]
func RegoAdd(context *gin.Context) {
	var form RegoForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	newRegoDocument := model.RegoDocument{
		Method:    form.Method,
		Path:      form.FilePath,
		Name:      form.Filename,
		Content:   form.FileContent,
		Arguments: util.ArgumentsParser(form.FileContent),
	}
	insertResult, err := gorm.Collections.RegoCollection.InsertOne(ctx.TODO(), newRegoDocument)
	if err != nil {
		return
	}
	h.SendResponse(context, nil, insertResult)
}


// ProjectResourceAdd 	api
// @Summary            ProjectResourceAdd
// @Description      Add a group and role to the given resource.
// @Tags                 projectResource
// @Accept             application/json
// @Produce            application/json
// @Security           Token
// @Param              project  header            string                                 true  "Project"
// @Param              role     header            string                                 true  "Role"
// @Param                       form    body                ProjectResourceForm  true  "form"
// @Success            200              {object}          handler.Response
// @Failure            400              {object}          handler.Response
// @Router                      /api/v1/projectResource/add [post]
func ProjectResourceAdd(context *gin.Context) {
	var form ProjectResourceForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	addedResourceName := form.ResourceName
	addedProjectName := form.ProjectName
	addedRoleName := form.RoleName
	resource := model.ProjectResources{
		GopaProjectRoles: model.GopaProjectRoles{
			ProjectName: addedProjectName,
			RoleName:    addedRoleName,
		},
		ResourceRouter: addedResourceName,
	}
	db := gorm.DB.Self
	res := db.Create(&resource)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}
