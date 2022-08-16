package api

import (
	ctx "context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopa/gorm"
	h "gopa/handler"
	"gopa/model"
	"gopa/pkg/errno"
)

type GopaProjects struct {
	ProjectName string `json:"project_name"`
}

type GopaRoles struct {
	RoleName string `json:"role_name"`
}

type GopaProjectRoles struct {
	GopaProjects
	GopaRoles
}

type GopaMembers struct {
	Username string `json:"username"`
	Project  string `json:"projecr"`
	Role     string `json:"role"`
}

type GopaApplication struct {
	ResourceName string `json:"resource_name"`
	Address      string `json:"address"`
	Routers      string `json:"routers"`
}

type ProjectResource struct {
	ProjectName  string `json:"project_name"`
	RoleName     string `json:"role_name"`
	ResourceRouter string `json:"resource_router"`
}

// ProjectsList 		api
// @Summary          ProjectsList
// @Description    list all projects of OPA
// @Tags               project
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project            header    string    true  "Project"
// @Param            role               header    string    true  "Role"
// @Success          200      {object}  handler.Response
// @Failure          400      {object}  handler.Response
// @Router                    /api/v1/project/list [get]
func ProjectsList(context *gin.Context) {
	var data []GopaProjects
	var result []string
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	for _, project := range data {
		result = append(result, project.ProjectName)
	}
	h.SendResponse(context, errno.OK, result)
}

// RolesList 		api
// @Summary          RolesList
// @Description    list all roles of OPA
// @Tags               role
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project            header    string    true  "Project"
// @Param            role               header    string    true  "Role"
// @Success          200      {object}  handler.Response
// @Failure          400      {object}  handler.Response
// @Router                    /api/v1/role/list [get]
func RolesList(context *gin.Context) {
	var data []GopaRoles
	var result []string
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	for _, role := range data {
		result = append(result, role.RoleName)
	}
	h.SendResponse(context, errno.OK, result)
}

// ProjectRoleList 	api
// @Summary          ProjectRoleList
// @Description    list all projects and corresponding roles of OPA
// @Tags                projectRole
// @Accept            application/json
// @Produce           application/json
// @Security          Token
// @Param             project            header    string    true  "Project"
// @Param             role               header    string    true  "Role"
// @Success           200      {object}  handler.Response
// @Failure           400      {object}  handler.Response
// @Router                     /api/v1/projectRole/list [get]
func ProjectRoleList(context *gin.Context) {
	var data []GopaProjectRoles
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, errno.OK, data)
}

// UserList 		api
// @Summary          UserList
// @Description    list all users of OPA
// @Tags               user
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project            header    string    true  "Project"
// @Param            role               header    string    true  "Role"
// @Success          200      {object}  handler.Response
// @Failure          400      {object}  handler.Response
// @Router       /api/v1/user/list [get]
func UserList(context *gin.Context) {
	var data []GopaMembers
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, errno.OK, data)
}

// ApplicationList 	api
// @Summary          ApplicationList
// @Description    list all applications of OPA
// @Tags               application
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project            header    string    true  "Project"
// @Param            role               header    string    true  "Role"
// @Success          200      {object}  handler.Response
// @Failure          400      {object}  handler.Response
// @Router       /api/v1/application/list [get]
func ApplicationList(context *gin.Context) {
	var data []GopaApplication
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, errno.OK, data)
}

// RegoList 		api
// @Summary          RegoList
// @Description    Return a rego file by given path currently in database. If filepath is not specified, API returns all rego files.
// @Tags               rego
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project            header            string    true    "Project"
// @Param            role               header            string    true    "Role"
// @Param                     filepath          query               string    false    "specify a rego path"
// @Success          200                        {object}          handler.Response
// @Failure          400                        {object}          handler.Response
// @Router                    /api/v1/rego/list [get]
func RegoList(context *gin.Context) {
	path, ok := context.GetQuery("filepath")
	var filter bson.D
	var result model.RegoDocument
	if !ok {
		var results []model.RegoDocument
		filter = bson.D{}
		cur, err := gorm.Collections.RegoCollection.Find(ctx.TODO(), filter)
		if err != nil {
			h.SendResponse400(context, err, nil)
			return
		}
		for cur.Next(ctx.TODO()) {
			err := cur.Decode(&result)
			if err != nil {
				h.SendResponse400(context, err, nil)
				return
			}
			results = append(results, result)
		}
		h.SendResponse(context, nil, results)
	} else {
		filter = bson.D{{"path", path}}
		err := gorm.Collections.RegoCollection.FindOne(ctx.TODO(), filter).Decode(&result)
		if err != nil {
			h.SendResponse400(context, err, nil)
			return
		}
		h.SendResponse(context, nil, result)
	}
}

// ProjectResourceList 	api
// @Summary            ProjectResourceList
// @Description      Return all groups and roles that belong to one resource.
// @Tags                 projectResource
// @Accept             application/json
// @Produce            application/json
// @Success            200  {object}  handler.Response
// @Failure            400  {object}  handler.Response
// @Router                  /api/v1/projectResource/list [get]
func ProjectResourceList(context *gin.Context) {
	var data []ProjectResource
	db := gorm.DB.Self
	res := db.Find(&data)
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, errno.OK, data)
}
