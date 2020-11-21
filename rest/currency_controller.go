package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
)

// GET /currency/config
func CurrencyConfigService(ctx *gin.Context) {
	out := CommonResp{}

	currencies, err := service.GetValidCurrencies()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = currencies

	ctx.JSON(http.StatusOK, out)
}

// Get /currency/depositInfo
func CurrencyDepositInfoService(ctx *gin.Context) {
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

	machine, err := service.GetCurrencyDepositsByUserId(address.Id, before, after, limit)
	if machine == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /currency/withdraw
func CurrencyWithdrawService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var withdraw CurrencyWithdrawRequest
	err := ctx.ShouldBindJSON(&withdraw)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] CurrencyWithdrawService request param: %v", withdraw)

	currency, err := service.GetCurrencyByCoin(withdraw.Coin)
	if currency == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.CurrencyWithdraw(address, currency, withdraw.Address, withdraw.Number)
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

// Get /currency/withdrawInfo
func CurrencyWithdrawInfoService(ctx *gin.Context) {
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

	machine, err := service.GetCurrencyWithdrawsByUserId(address.Id, before, after, limit)
	if machine == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}
