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

	before, _ := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, _ := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

	deposits, err := service.GetAddressDepositsByUserId(address.Id, before, after, limit)
	if deposits == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(deposits) > 0 {
		newBefore = deposits[0].Id
		newAfter = deposits[len(deposits)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   deposits,
	}
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
		mylog.Logger.Error().Msgf("[Rest] AddressWithdrawService AddressWithdraw err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
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

	before, _ := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, _ := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

	withdraws, err := service.GetAddressWithdrawsByUserId(address.Id, before, after, limit)
	if withdraws == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(withdraws) > 0 {
		newBefore = withdraws[0].Id
		newAfter = withdraws[len(withdraws)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   withdraws,
	}
	ctx.JSON(http.StatusOK, out)
}
