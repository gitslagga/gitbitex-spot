package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
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
		out.RespDesc = err.Error()
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

	if a.From == a.To || (a.Currency != models.AccountCurrencyBite && a.Currency != models.AccountCurrencyUsdt) ||
		((a.From == models.AccountPoolTransfer || a.To == models.AccountPoolTransfer) && a.Currency != models.AccountCurrencyBite) {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.AccountTransfer(address.Id, a.From, a.To, a.Currency, a.Number)
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

	accountTransfer, err := service.GetAccountTransferByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = accountTransfer
	ctx.JSON(http.StatusOK, out)
}
