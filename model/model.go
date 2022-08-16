package model

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel config
type BaseModel struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type GopaProjects struct {
	BaseModel
	ProjectName string
}

type GopaRoles struct {
	BaseModel
	RoleName string
}

type GopaProjectRoles struct {
	BaseModel
	ProjectName string
	RoleName    string
}

type GopaMembers struct {
	BaseModel
	Username string
	Password string
	Project  string
	Role     string
}

type GopaApplication struct {
	BaseModel
	ResourceName string
	Address      string
	Routers      string
}

type ProjectResources struct {
	GopaProjectRoles
	ResourceRouter string
}

// TableName 结构体映射表名称
func (GopaProjects) TableName() string {
	return "gopa_projects"
}

// TableName 结构体映射表名称
func (GopaRoles) TableName() string {
	return "gopa_roles"
}

// TableName 结构体映射表名称
func (GopaProjectRoles) TableName() string {
	return "gopa_project_roles"
}

// TableName 结构体映射表名称
func (GopaMembers) TableName() string {
	return "gopa_members"
}

// RegoDocument MongoDB
type RegoDocument struct {
	//ID          string
	Method    string
	Path      string
	Name      string
	Content   string
	Arguments []string
}

func (GopaApplication) TableName() string {
	return "gopa_applications"
}

func (ProjectResources) TableName() string {
	return "project_resources"
}

// UserRole 用户角色
type UserRole struct {
	ID     string `json:"id"`      // 唯一标识
	UserID string `json:"user_id"` // 用户ID
	RoleID string `json:"role_id"` // 角色ID
}

// UserRoles 角色菜单列表
type UserRoles []*UserRole

// User 用户对象
type User struct {
	ID        string    `json:"id"`                                    // 唯一标识
	UserName  string    `json:"user_name" binding:"required"`          // 用户名
	RealName  string    `json:"real_name" binding:"required"`          // 真实姓名
	Password  string    `json:"password"`                              // 密码
	Phone     string    `json:"phone"`                                 // 手机号
	Email     string    `json:"email"`                                 // 邮箱
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator   string    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // 角色授权
}
