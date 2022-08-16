package api

import (
	ctx "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopa/gorm"
	h "gopa/handler"
	"gopa/model"
	"gopa/schema"
)

type UpdatePermissionForm struct {
	schema.ActionUserForm
	Username string `json:"username"`
	Project    string `json:"project"`
	Role     string `json:"role"`
}

type UpdatedRegoDocumentForm struct {
	FilePath   string `json:"path"`
	NewContent string `json:"content"`
}

type UpdatedProjectResourceForm struct {
	ResourceName string `json:"resource_name"`
	ProjectName  string `json:"project_name"`
	RoleName     string `json:"role_name"`
}

// UserUpdate 	api
// @Summary        UserUpdate
// @Description  Update the group or role of a given user
// @Tags             user
// @Accept         application/json
// @Produce        application/json
// @Security       Token
// @Param          project  header            string                                    true  "Project"
// @Param          role     header            string                                    true  "Role"
// @Param                   form              body        UpdatePermissionForm    true        "form"
// @Success        200              {object}          handler.Response
// @Failure        400              {object}          handler.Response
// @Router                  /api/v1/user/update [post]
func UserUpdate(context *gin.Context) {
	var form UpdatePermissionForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	db := gorm.DB.Self
	var user model.GopaMembers
	user.Username = form.Username
	res := db.Model(&user).Where("username = ?", form.Username).Updates(model.GopaMembers{Project: form.Project, Role: form.Role})
	if res.Error != nil {
		h.SendResponse400(context, res.Error, nil)
		return
	}
	h.SendResponse(context, nil, "success")
}

// RegoUpdate 	api
// @Summary        RegoUpdate
// @Description  Update the content of a rego file by given path
// @Tags             rego
// @Accept         application/json
// @Produce        application/json
// @Security       Token
// @Param          project  header            string                                         true  "Project"
// @Param          role     header            string                                         true  "Role"
// @Param                   form              body        UpdatedRegoDocumentForm    true    "form"
// @Success        200              {object}          handler.Response
// @Failure        400              {object}          handler.Response
// @Router                  /api/v1/rego/update [post]
func RegoUpdate(context *gin.Context) {
	var form UpdatedRegoDocumentForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	filter := bson.D{{"path", form.FilePath}}
	update := bson.D{
		{"$set", bson.D{{"content", form.NewContent}}},
	}
	updateResult, err := gorm.Collections.RegoCollection.UpdateOne(ctx.TODO(), filter, update)
	if err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	fmt.Println("Matched updated content: ", updateResult.ModifiedCount)
	h.SendResponse(context, nil, updateResult.ModifiedCount)
}

// ProjectResourceUpdate 	api
// @Summary              ProjectResourceUpdate
// @Description        Update the permission of a group with its related role to a resource
// @Tags                   projectResource
// @Accept               application/json
// @Produce              application/json
// @Security             Token
// @Param                project  header            string                                            true  "Project"
// @Param                role     header            string                                            true  "Role"
// @Param                         form    body                UpdatedProjectResourceForm    true    "form"
// @Success              200              {object}          handler.Response
// @Failure              400              {object}          handler.Response
// @Router                        /api/v1/projectResource/update [post]
func ProjectResourceUpdate(context *gin.Context) {
	var form UpdatedProjectResourceForm
	if err := context.BindJSON(&form); err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	filter := bson.D{{"resouce_router", form.ResourceName}}
	update := bson.D{{"project_name", form.ProjectName}, {"role_name", form.RoleName}}
	updateResult, err := gorm.Collections.RegoCollection.UpdateOne(ctx.TODO(), filter, update)
	if err != nil {
		h.SendResponse400(context, err, nil)
		return
	}
	fmt.Println("Matched updated content: ", updateResult.ModifiedCount)
	h.SendResponse(context, nil, updateResult.ModifiedCount)
}
