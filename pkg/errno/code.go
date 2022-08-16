package errno

var (
	// Common errors
	OK                = &Errno{Code: 200, Message: "OK"}
	InstancePullInOK  = &Errno{Code: 200, Message: "Instance success pull in."}
	InstancePullOutOK = &Errno{Code: 200, Message: "Instance success pull out."}

	InternalServerError = &Errno{Code: 10001, Message: "Internal server error."}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	// Eureka errors
	ErrEurekaClient           = &Errno{Code: 21001, Message: "Error occurred while call the eureka server."}
	ErrEurekaInstanceNotFound = &Errno{Code: 21002, Message: "Instance not found in eureka server."}
	ErrHealthCheck            = &Errno{Code: 21003, Message: "HealthCheck Faild."}
	ErrEurekaPullIn           = &Errno{Code: 21004, Message: "Instance pull in faild."}
	ErrEurekaPullOut          = &Errno{Code: 21005, Message: "Instance pull out faild."}
	ErrEurekaPullOutReject    = &Errno{Code: 22006, Message: "AppGroup has only one instance, cannot pull out."}

	// AliCloud errors
	ErrAliCloudSlb                     = &Errno{Code: 22001, Message: "Error occurred while call the alicloud slb api."}
	ErrAliCloudSlbNotFound             = &Errno{Code: 22002, Message: "Not found the instance belong to any slb."}
	ErrAliCloudSlbInstancePortNotFound = &Errno{Code: 22003, Message: "Not found the instance port active in any slb."}
	ErrAliCloudPullIn                  = &Errno{Code: 22004, Message: "Instance pull in faild."}
	ErrAliCloudPullOut                 = &Errno{Code: 22005, Message: "Instance pull out faild."}
	ErrAliCloudSlbReject               = &Errno{Code: 22006, Message: "vGroup has only one instance, cannot pull out."}

	// user errors
	NotOK             = &Errno{Code: 20001, Message: "Not OK"}
	ErrValidation     = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase       = &Errno{Code: 20002, Message: "Database error."}
	ErrToken          = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}
	ErrInstanceStatus = &Errno{Code: 20004, Message: "Instance status must be up or down."}

	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
	ErrEncrypt           = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}
)
