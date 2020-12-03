package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"time"
)

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
		mylog.Logger.Error().Msgf("[Rest] BackendWithdrawService BackendWithdraw err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// Get /backend/issue/list
func BackendIssueListService(ctx *gin.Context) {
	out := CommonResp{}

	list, err := service.BackendIssueList()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = list

	ctx.JSON(http.StatusOK, out)
}

// Post /backend/issue/start
func BackendIssueStartService(ctx *gin.Context) {
	out := CommonResp{}

	err := service.BackendIssueStart()
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] BackendIssueStartService BackendIssueStart err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// Get /backend/holding/list
func BackendHoldingListService(ctx *gin.Context) {
	out := CommonResp{}

	list, err := service.BackendHoldingList()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = list

	ctx.JSON(http.StatusOK, out)
}

// Post /backend/holding/start
func BackendHoldingStartService(ctx *gin.Context) {
	out := CommonResp{}

	addressHolding, err := service.GetLastAddressHolding()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if addressHolding != nil {
		currentTime := time.Now()
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 00, 00, 00, 00, currentTime.Location())
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

		if addressHolding.CreatedAt.After(startTime) && addressHolding.CreatedAt.Before(endTime) {
			out.RespCode = EC_DAY_PROFIT_RELEASED
			out.RespDesc = ErrorCodeMessage(EC_DAY_PROFIT_RELEASED)
			ctx.JSON(http.StatusOK, out)
			return
		}
	}

	err = service.BackendHoldingStart()
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] BackendHoldingStartService BackendHoldingStart err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}
