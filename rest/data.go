package rest

import (
	"errors"
)

type ErrorCode int

const (
	EC_NONE              ErrorCode = 0
	EC_PARAMS_ERR                  = 10000
	EC_NETWORK_ERR                 = 10001
	EC_JSON_MARSHAL_ERR            = 10002
	EC_TOKEN_INVALID               = 10003
	EC_RESPONSE_DATA_ERR           = 10004
	EC_REQUEST_DATA_ERR            = 10005

	EC_PASSWORD_ERR       = 9000
	EC_MNEMONIC_INCORRECT = 9001
	EC_PASSWORD_INCORRECT = 9002
)

func (c ErrorCode) Code() (r int) {
	r = int(c)
	return
}

func (c ErrorCode) Error() (r error) {
	r = errors.New(c.String())
	return
}

func (c ErrorCode) String() (r string) {
	switch c {
	case EC_NONE:
		r = "SUCCESS"
	case EC_NETWORK_ERR:
		r = "请求错误|Request error"
	case EC_PARAMS_ERR:
		r = "参数错误|Params error"
	case EC_JSON_MARSHAL_ERR:
		r = "json格式异常|Json format exception"
	case EC_TOKEN_INVALID:
		r = "无效的Token|Invalid token"
	case EC_RESPONSE_DATA_ERR:
		r = "请重新登录|Please login again"
	case EC_REQUEST_DATA_ERR:
		r = "非法请求|Illegal request"

	case EC_PASSWORD_ERR:
		r = "密码长度必须至少为6个字符|password must be of minimum 6 characters length"
	case EC_MNEMONIC_INCORRECT:
		r = "助记词不正确|Mnemonic is incorrect"
	case EC_PASSWORD_INCORRECT:
		r = "旧密码不正确|Old password is incorrect"

	default:
	}
	return
}

func ErrorCodeMessage(c ErrorCode) (r string) {
	return c.String()
}

type CommonResp struct {
	RespCode int         `json:"respCode" form:"respCode"`
	RespDesc string      `json:"respDesc" form:"respDesc"`
	RespData interface{} `json:"respData,omitempty" form:"respData"`
}

//rest api request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mnemonic string `json:"mnemonic" binding:"required"`
}

type LoginRequest struct {
	Mnemonic   string `json:"mnemonic"`
	PrivateKey string `json:"private_key"`
	Password   string `json:"password" binding:"required"`
}

type FindPasswordRequest struct {
	PrivateKey string `json:"private_key" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ModifyPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
