package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
)

// 获取用户余额
// GET /accounts?currency=BTC&currency=USDT
func GetAccounts(ctx *gin.Context) {
	var accountVos []*AccountVo
	currencies := ctx.QueryArray("currency")
	if len(currencies) != 0 {
		for _, currency := range currencies {
			account, err := service.GetAccount(GetCurrentUser(ctx).Id, currency)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
				return
			}
			if account == nil {
				continue
			}

			accountVos = append(accountVos, newAccountVo(account))
		}
	} else {
		accounts, err := service.GetAccountsByUserId(GetCurrentUser(ctx).Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
			return
		}
		for _, account := range accounts {
			accountVos = append(accountVos, newAccountVo(account))
		}
	}
	ctx.JSON(http.StatusOK, accountVos)
}

// GET /account/address
func AccountAddressService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountAddress, err := service.AccountAddress(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = accountAddress

	ctx.JSON(http.StatusOK, out)
}

// POST /account/transfer
func AccountTransferService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var a AccountTransferRequest
	err := ctx.ShouldBindJSON(&a)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if a.From == a.To {
		out.RespCode = EC_THE_SAME_ACCOUNT
		out.RespDesc = ErrorCodeMessage(EC_THE_SAME_ACCOUNT)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if a.From == models.AccountShopTransfer {
		out.RespCode = EC_SHOP_ONLY_ENTER
		out.RespDesc = ErrorCodeMessage(EC_SHOP_ONLY_ENTER)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if a.Currency != models.AccountCurrencyBite && a.Currency != models.AccountCurrencyUsdt {
		out.RespCode = EC_CURRENCY_NOT_EXISTS
		out.RespDesc = ErrorCodeMessage(EC_CURRENCY_NOT_EXISTS)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if (a.From == models.AccountPoolTransfer || a.To == models.AccountPoolTransfer) && a.Currency != models.AccountCurrencyBite {
		out.RespCode = EC_POOL_ONLY_BITE
		out.RespDesc = ErrorCodeMessage(EC_POOL_ONLY_BITE)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.AccountTransfer(address.Id, a.From, a.To, a.Currency, a.Number)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AccountTransferService AccountTransfer err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// GET /account/transferInfo
func AccountTransferInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	before, err1 := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, err2 := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, err3 := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountTransfer, err := service.GetAccountTransferByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(accountTransfer) > 0 {
		newBefore = accountTransfer[0].Id
		newAfter = accountTransfer[len(accountTransfer)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   accountTransfer,
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /account/scan
func AccountScanService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var accountScan AccountScanRequest
	err := ctx.ShouldBindJSON(&accountScan)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.AccountScan(address.Id, accountScan.Url, accountScan.Number)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AccountScanService AccountScan err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// GET /account/scanInfo
func AccountScanInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	before, err1 := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, err2 := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, err3 := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountScan, err := service.GetAccountScanByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(accountScan) > 0 {
		newBefore = accountScan[0].Id
		newAfter = accountScan[len(accountScan)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   accountScan,
	}
	ctx.JSON(http.StatusOK, out)
}
