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

	accountCurrency, err := service.GetAccountsByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountAsset, err := service.GetAccountsAssetByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountPool, err := service.GetAccountsPoolByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountShop, err := service.GetAccountsShopByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = map[string]interface{}{
		"accountCurrency": accountCurrency,
		"accountAsset":    accountAsset,
		"accountPool":     accountPool,
		"accountShop":     accountShop,
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /account/convert
func AccountConvertService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var accountConvert AccountConvertRequest
	err := ctx.ShouldBindJSON(&accountConvert)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.AccountConvert(address, accountConvert.Number)
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

// GET /account/convertInfo
func AccountConvertInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	accountConvert, err := service.GetAccountConvertByUserId(address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = accountConvert
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

	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	if a.From == a.To || (a.From == 1 && a.To == 2 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 1 && a.To == 3 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) ||
		(a.From == 1 && a.To == 4 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) ||
		(a.From == 2 && a.To == 1 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 2 && a.To == 3 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 2 && a.To == 4 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 3 && a.To == 1 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) ||
		(a.From == 3 && a.To == 2 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 3 && a.To == 4 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) ||
		(a.From == 4 && a.To == 1 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) ||
		(a.From == 4 && a.To == 2 && a.Currency != models.CURRENCY_BITC) ||
		(a.From == 4 && a.To == 3 && a.Currency != models.CURRENCY_BITC && a.Currency != models.CURRENCY_USDT) {
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
