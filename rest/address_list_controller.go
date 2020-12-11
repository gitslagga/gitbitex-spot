package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// Get /address/list
func AddressListService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	addressLists, err := service.AddressListService(address)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if addressLists != nil {
		out.RespData = addressLists
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /address/addList
func AddressAddListService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var login LoginRequest
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
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
	if len(login.Username) < 3 {
		out.RespCode = EC_USERNAME_ERR
		out.RespDesc = ErrorCodeMessage(EC_USERNAME_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressAddListService request param: %v", login)

	err = service.AddressAddList(address.Id, login.Username, encryptPassword(login.Password), login.Mnemonic, login.PrivateKey)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressAddListService AddressAddList err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// DELETE /address/delList
func AddressDelListService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var addressListRequest AddressListRequest
	err := ctx.ShouldBindJSON(&addressListRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressDelListService request param: %v", addressListRequest)

	err = service.DeleteAddressList(addressListRequest.Address)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressDelListService DeleteAddressList err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// POST /address/switchList
func AddressSwitchListService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var addressListRequest AddressListRequest
	err := ctx.ShouldBindJSON(&addressListRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressSwitchListService request param: %v", addressListRequest)

	addressList, err := service.GetAddressListByAddress(addressListRequest.Address)
	if err != nil || addressList == nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	address, err = service.AddressSwitchList(address, addressList)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressSwitchListService AddressSwitchList err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	token, err := service.CreateFrontendToken(address)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressSwitchListService CreateFrontendToken err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = token

	ctx.JSON(http.StatusOK, out)
}
