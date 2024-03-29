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
	EC_WHITE_LIST_ERR              = 10006

	EC_USERNAME_ERR          = 9000
	EC_PASSWORD_ERR          = 9001
	EC_MNEMONIC_INCORRECT    = 9002
	EC_PRIVATE_KEY_INCORRECT = 9003
	EC_PASSWORD_INCORRECT    = 9004
	EC_USERNAME_PASSWORD_ERR = 9005
	EC_ACTIVATION_SELF_ERR   = 9006
	EC_USERNAME_EXISTS_ERR   = 9007

	EC_CLIENT_OID_ERR        = 8000
	EC_ORDER_NOT_EXISTS      = 8001
	EC_THE_SAME_ACCOUNT      = 8002
	EC_SHOP_ONLY_ENTER       = 8003
	EC_CURRENCY_NOT_EXISTS   = 8004
	EC_POOL_ONLY_BITE        = 8005
	EC_DAY_PROFIT_RELEASED   = 8006
	EC_GROUP_PUBLISH_EXISTS  = 8007
	EC_GROUP_JOIN_NOT_EXISTS = 8008
	EC_GROUP_JOIN_REPEAT_ERR = 8009
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
	case EC_WHITE_LIST_ERR:
		r = "不在白名单中|Not in whitelist"

	case EC_PASSWORD_ERR:
		r = "密码长度必须至少为6个字符|password must be of minimum 6 characters length"
	case EC_MNEMONIC_INCORRECT:
		r = "助记词不正确|Mnemonic is incorrect"
	case EC_PRIVATE_KEY_INCORRECT:
		r = "私钥不正确|Private key is incorrect"
	case EC_PASSWORD_INCORRECT:
		r = "旧密码不正确|Old password is incorrect"
	case EC_USERNAME_PASSWORD_ERR:
		r = "用户名或密码错误|Username or password error"
	case EC_ACTIVATION_SELF_ERR:
		r = "不能自己激活自己|Can not activate yourself"
	case EC_USERNAME_EXISTS_ERR:
		r = "用户名已存在:Username is already exists"

	case EC_CLIENT_OID_ERR:
		r = "无效的客户ID: %v|invalid client_oid"
	case EC_ORDER_NOT_EXISTS:
		r = "订单不存在|order not found"
	case EC_THE_SAME_ACCOUNT:
		r = "相同账户不能划转|The same account cannot be transferred"
	case EC_SHOP_ONLY_ENTER:
		r = "商城账户只进不出|Mall account can only enter but not exit"
	case EC_CURRENCY_NOT_EXISTS:
		r = "币种不存在|Currency does not exist"
	case EC_POOL_ONLY_BITE:
		r = "矿池账户只支持BITE|Mining pool accounts only support BITE"
	case EC_DAY_PROFIT_RELEASED:
		r = "今日收益已释放|Today's earnings have been released"
	case EC_GROUP_PUBLISH_EXISTS:
		r = "拼团已经存在|The group already exists"
	case EC_GROUP_JOIN_NOT_EXISTS:
		r = "该拼团不存在|The group does not exist"
	case EC_GROUP_JOIN_REPEAT_ERR:
		r = "不能重复参加|Cannot repeat join"

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

type PageResp struct {
	Before int64       `json:"before" form:"before"`
	After  int64       `json:"after" form:"after"`
	List   interface{} `json:"list,omitempty" form:"list"`
}

//rest api request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mnemonic string `json:"mnemonic" binding:"required"`
}

type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Mnemonic   string `json:"mnemonic"`
	PrivateKey string `json:"private_key"`
}

type AddressAddListRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mnemonic string `json:"mnemonic" binding:"required"`
}

type FindPasswordRequest struct {
	PrivateKey string `json:"private_key" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ModifyPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ActivationRequest struct {
	Address string  `json:"address" binding:"required"`
	Number  float64 `json:"number" binding:"required"`
}

type BuyMachineRequest struct {
	MachineId int64  `json:"machine_id" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
}

type MachineConvertRequest struct {
	ConvertType int     `json:"convert_type" binding:"required"`
	Number      float64 `json:"number" binding:"required"`
}

type AccountTransferRequest struct {
	From     int     `json:"from" binding:"required"`
	To       int     `json:"to" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
	Number   float64 `json:"number" binding:"required"`
}

type AccountScanRequest struct {
	Url    string  `json:"url" binding:"required"`
	Number float64 `json:"number" binding:"required"`
}

type AddressWithdrawRequest struct {
	Address string  `json:"address" binding:"required"`
	Coin    string  `json:"coin" binding:"required"`
	Number  float64 `json:"number" binding:"required"`
}

type AddressPassWithdrawRequest struct {
	OrderSN string `json:"order_sn" binding:"required"`
	Status  int    `json:"status" binding:"required"`
}

type AddressListRequest struct {
	Address string `json:"address" binding:"required"`
}

type GroupRequest struct {
	Coin string `json:"coin" binding:"required"`
}

type GroupJoinRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
}
