package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// GET /admin/login
func AdminLoginService(ctx *gin.Context) {
	out := CommonResp{}
	var adminLogin AdminLoginRequest
	err := ctx.ShouldBindJSON(&adminLogin)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if len(adminLogin.Password) < 8 {
		out.RespCode = EC_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AdminLoginService request param: %v", adminLogin)

	admin, err := service.GetAdminByUsername(adminLogin.Username)
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] AdminLoginService GetAdminByUsername err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if admin == nil || admin.Password != encryptPassword(adminLogin.Password) {
		mylog.DataLogger.Info().Msgf("[Rest] AdminLoginService username or password error")
		out.RespCode = EC_USERNAME_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_USERNAME_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	token, err := service.CreateBackendToken(admin)
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] AdminLoginService CreateBackendToken err: %v", err)
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

// Get /admin/info
func AdminService(ctx *gin.Context) {
	out := CommonResp{}
	admin := GetCurrentAdmin(ctx)
	if admin == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = admin
	ctx.JSON(http.StatusOK, out)
}
