package schema


const (
	InvalidMethodMessage     = "invalid method"
	NoMethodSpecifiedMessage = "you need to specify a method for this APi"
	AccessForbiddenMessage = "you don't have access to this resource or perform such an action"
)


// ActionUserForm 用户进行的所有操作均需要进行鉴权
// 该结构体包含鉴权所需的所有字段
type ActionUserForm struct {
	Project string `json:"project"`
	Role     string `json:"role"`
	//FromGroup    string `json:"from_group"`
}

// RequestDecision 用户信息发送到OPA获取鉴权结果的表单结构体
type RequestDecision struct {
	Input ActionUserForm `json:"input"`
}

type RequestResult struct {
	RowsAffected int64
	StatusCode   int64
	Message      string
}
