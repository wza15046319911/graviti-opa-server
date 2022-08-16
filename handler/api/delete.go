package api

import (
	ctx "context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopa/gorm"
	h "gopa/handler"
	"gopa/model"
)

type ProjectDeleteForm struct {
	DeletedGroupName string `json:"project_name"`
}

type DeletedRoleForm struct {
	DeletedRoleName string `json:"role_name"`
}

type RoleDeleteFromGroupForm struct {
	DeletedGroupName string `json:"project_name"`
	DeletedRoleName  string `json:"role_name"`
}

type UserDeleteForm struct {
	Username string `json:"username"`
}

type ApplicationDeleteForm struct {
	DeletedApplicationName string `json:"application_name"`
}

type DeletedRegoDocumentForm struct {
	FilePath string `json:"path"`
}


// ProjectDelete 	api
// @Summary          ProjectDelete
// @Description    Delete a project
// @Tags               project
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project  header            string                           true  "Project"
// @Param            role     header            string                           true  "Role"
// @Param                     form              body        ProjectDeleteForm    true  "form"
// @Success          200              {object}          handler.Response
// @Failure          400              {object}          handler.Response
// @Router                    /api/v1/project/delete [post]
func ProjectDelete(context *gin.Context) {
	var form ProjectDeleteForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	group := form.DeletedGroupName
	db := gorm.DB.Self
	res := db.Unscoped().Where("project_name = ?", group).Delete(&model.GopaProjects{ProjectName: group})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// RoleDelete 	api
// @Summary        RoleDelete
// @Description  Delete a role
// @Tags             role
// @Accept         application/json
// @Produce        application/json
// @Security       Token
// @Param          project  header            string                         true  "Project"
// @Param          role     header            string                         true  "Role"
// @Param                   form              body        DeletedRoleForm    true  "form"
// @Success        200              {object}          handler.Response
// @Failure        400              {object}          handler.Response
// @Router                  /api/v1/role/delete [post]
func RoleDelete(context *gin.Context) {
	var form DeletedRoleForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	role := form.DeletedRoleName
	db := gorm.DB.Self
	res := db.Unscoped().Where("role_name = ?", role).Delete(&model.GopaRoles{RoleName: role})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// ProjectRoleDelete 	api
// @Summary            ProjectRoleDelete
// @Description      Delete a role from a given project
// @Tags                 projectRole
// @Accept             application/json
// @Produce            application/json
// @Security           Token
// @Param              project  header            string                                         true  "Project"
// @Param              role     header            string                                         true  "Role"
// @Param                       form              body        RoleDeleteFromGroupForm    true    "form"
// @Success            200              {object}          handler.Response
// @Failure            400              {object}          handler.Response
// @Router                      /api/v1/projectRole/delete [post]
func ProjectRoleDelete(context *gin.Context) {
	var form RoleDeleteFromGroupForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	group := form.DeletedGroupName
	role := form.DeletedRoleName
	db := gorm.DB.Self
	res := db.Unscoped().Where("project_name = ?", group).Where("role_name = ?", role).Delete(
		&model.GopaProjectRoles{
			ProjectName: group,
			RoleName: role,
		})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}

// UserDelete 	api
// @Summary        UserDelete
// @Description  Delete a user
// @Tags             user
// @Accept         application/json
// @Produce        application/json
// @Security       Token
// @Param          project  header            string                      true  "Project"
// @Param          role     header            string                      true  "Role"
// @Param                   form              body        UserDeleteForm        true  "form"
// @Success        200              {object}          handler.Response
// @Failure        400              {object}          handler.Response
// @Router                  /api/v1/user/delete [post]
func UserDelete(context *gin.Context) {
	var form UserDeleteForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	username := form.Username
	db := gorm.DB.Self
	res := db.Unscoped().Where("username = ?", username).Delete(
		&model.GopaMembers{
			Username: username,
		})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}



// RegoDelete 	api
// @Summary        RegoDelete
// @Description  Delete a rego file by given path
// @Tags             rego
// @Accept         application/json
// @Produce        application/json
// @Security       Token
// @Param          project  header            string                                         true  "Project"
// @Param          role     header            string                                         true  "Role"
// @Param                   form              body        DeletedRegoDocumentForm    true    "form"
// @Success        200              {object}          handler.Response
// @Failure        400              {object}          handler.Response
// @Router                  /api/v1/rego/delete [post]
func RegoDelete(context *gin.Context) {
	var form DeletedRegoDocumentForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	filter := bson.D{{"path", form.FilePath}}
	deleteResult, err := gorm.Collections.RegoCollection.DeleteOne(ctx.TODO(), filter)
	if err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	h.SendResponse(context, nil, deleteResult.DeletedCount)
}

// ApplicationDelete api
// @Summary          ApplicationDelete
// @Description    Delete this application from OPA. This method will delete ALL RELATED INFORMATION from database.
// @Tags               application
// @Accept           application/json
// @Produce          application/json
// @Security         Token
// @Param            project  header            string                                     true  "Project"
// @Param            role     header            string                                     true  "Role"
// @Param                     form              body        ApplicationDeleteForm    true        "form"
// @Success          200              {object}          handler.Response
// @Failure          400              {object}          handler.Response
// @Router                    /api/v1/application/delete [post]
func ApplicationDelete(context *gin.Context) {
	var form ApplicationDeleteForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	name := form.DeletedApplicationName
	db := gorm.DB.Self
	res := db.Unscoped().Where("resource_name = ?", name).Delete(
		&model.GopaApplication{
			ResourceName: name,
		})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, res.RowsAffected)
}


// ProjectResourceDelete 	api
// @Summary              ProjectResourceDelete
// @Description        Delete a projectRole from a resource.
// @Tags                   projectResource
// @Accept               application/json
// @Produce              application/json
// @Security             Token
// @Param                project  header            string                                     true  "Project"
// @Param                role     header            string                                     true  "Role"
// @Param                         form              body        ApplicationDeleteForm    true        "form"
// @Success              200              {object}          handler.Response
// @Failure              400              {object}          handler.Response
// @Router                        /api/v1/projectResource/delete [post]
func ProjectResourceDelete(context *gin.Context) {

}