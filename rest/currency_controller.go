package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
)

// GET /address/config
func AddressConfigService(ctx *gin.Context) {
	out := CommonResp{}

	configs, err := service.GetValidAddressConfig()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = configs

	ctx.JSON(http.StatusOK, out)
}

// Get /address/depositInfo
func AddressDepositInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	before, err1 := strconv.ParseInt(ctx.DefaultQuery("before", "0"), 10, 64)
	after, err2 := strconv.ParseInt(ctx.DefaultQuery("after", "11"), 10, 64)
	limit, err3 := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	deposits, err := service.GetAddressDepositsByUserId(address.Id, before, after, limit)
	if deposits == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = deposits
	ctx.JSON(http.StatusOK, out)
}

// POST /address/withdraw
func AddressWithdrawService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var withdraw AddressWithdrawRequest
	err := ctx.ShouldBindJSON(&withdraw)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressWithdrawService request param: %v", withdraw)

	config, err := service.GetAddressConfigByCoin(withdraw.Coin)
	if config == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.AddressWithdraw(address, config, withdraw.Address, withdraw.Number)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// Get /address/withdrawInfo
func AddressWithdrawInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	before, err1 := strconv.ParseInt(ctx.DefaultQuery("before", "0"), 10, 64)
	after, err2 := strconv.ParseInt(ctx.DefaultQuery("after", "11"), 10, 64)
	limit, err3 := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	withdraws, err := service.GetAddressWithdrawsByUserId(address.Id, before, after, limit)
	if withdraws == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = withdraws
	ctx.JSON(http.StatusOK, out)
}

// Post /backend/address/withdraw
func BackendWithdrawService(ctx *gin.Context) {
	out := CommonResp{}

	var addressPassWithdraw AddressPassWithdrawRequest
	err := ctx.ShouldBindJSON(&addressPassWithdraw)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] BackendWithdrawService request param: %v", addressPassWithdraw)

	withdraw, err := service.GetAddressWithdrawsByOrderSN(addressPassWithdraw.OrderSN)
	if withdraw == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.BackendWithdraw(withdraw, addressPassWithdraw.Status)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}
