package rest

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// POST /api/address/mnemonic
func MnemonicService(ctx *gin.Context) {
	out := CommonResp{}

	mnemonic, err := service.CreateMnemonic()
	if err != nil {
		mylog.Frontend.Error().Msgf("[Rest] MnemonicService AddressRegister err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = mnemonic

	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/register
func RegisterService(ctx *gin.Context) {
	out := CommonResp{}
	var register RegisterRequest
	err := ctx.ShouldBindJSON(&register)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if len(register.Password) < 8 {
		out.RespCode = EC_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Frontend.Info().Msgf("[Rest] RegisterService request param: %v", register)

	_, err = service.AddressRegister(register.Username, encryptPassword(register.Password), register.Mnemonic)
	if err != nil {
		mylog.Frontend.Error().Msgf("[Rest] RegisterService AddressRegister err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/login
func LoginService(ctx *gin.Context) {
	out := CommonResp{}
	var login LoginRequest
	err := ctx.ShouldBindJSON(&login)
	if err != nil || (login.Mnemonic == "" && login.PrivateKey == "") {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if len(login.Password) < 8 {
		out.RespCode = EC_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Frontend.Info().Msgf("[Rest] LoginService request param: %v", login)

	address, err := service.AddressLogin(login.Mnemonic, login.PrivateKey, encryptPassword(login.Password))
	if err != nil {
		mylog.Frontend.Error().Msgf("[Rest] RegisterService AddressLogin err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	token, err := service.CreateFrontendToken(address)
	if err != nil {
		mylog.Frontend.Error().Msgf("[Rest] LoginService CreateFrontendToken err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	ctx.SetCookie("accessToken", token, 7*24*60*60, "/", "", false, false)

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = token
	ctx.JSON(http.StatusOK, out)
}

// DELETE /api/address/logout
func LogoutService(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", -1, "/", "", false, false)

	out := CommonResp{
		RespCode: EC_NONE.Code(),
		RespDesc: EC_NONE.String(),
	}
	ctx.JSON(http.StatusOK, out)
}

// GET /api/address/info
func AddressService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = address
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/findPassword
func FindPasswordService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var findPassword FindPasswordRequest
	err := ctx.ShouldBindJSON(&findPassword)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Frontend.Info().Msgf("[Rest] FindPasswordService request param: %v", findPassword)

	if address.PrivateKey != findPassword.PrivateKey {
		out.RespCode = EC_PRIVATE_KEY_INCORRECT
		out.RespDesc = ErrorCodeMessage(EC_PRIVATE_KEY_INCORRECT)
		ctx.JSON(http.StatusOK, out)
		return
	}

	address.Password = encryptPassword(findPassword.Password)
	err = service.UpdateAddress(address)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/modifyPassword
func ModifyPasswordService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var modifyPassword ModifyPasswordRequest
	err := ctx.ShouldBindJSON(&modifyPassword)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Frontend.Info().Msgf("[Rest] ModifyPasswordService request param: %v", modifyPassword)

	if address.Password != encryptPassword(modifyPassword.OldPassword) {
		out.RespCode = EC_PASSWORD_INCORRECT
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_INCORRECT)
		ctx.JSON(http.StatusOK, out)
		return
	}

	address.Password = encryptPassword(modifyPassword.NewPassword)
	err = service.UpdateAddress(address)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/activation
func ActivationService(ctx *gin.Context) {
	//out := CommonResp{}
	//address := GetCurrentAddress(ctx)
	//if address == nil {
	//	out.RespCode = EC_TOKEN_INVALID
	//	out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
	//	ctx.JSON(http.StatusOK, out)
	//	return
	//}
	//
	//var activation ActivationRequest
	//err := ctx.ShouldBindJSON(&activation)
	//if err != nil {
	//	out.RespCode = EC_PARAMS_ERR
	//	out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
	//	ctx.JSON(http.StatusOK, out)
	//	return
	//}
	//
	//mylog.Frontend.Info().Msgf("[Rest] ActivationService request param: %v", activation)
	//
	//err = service.UpdateAddress(address)
	//if err != nil {
	//	out.RespCode = EC_NETWORK_ERR
	//	out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
	//	ctx.JSON(http.StatusOK, out)
	//	return
	//}
	//
	//out.RespCode = EC_NONE.Code()
	//out.RespDesc = EC_NONE.String()
	//ctx.JSON(http.StatusOK, out)
}

func encryptPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}
