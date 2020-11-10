package rest

import (
	"github.com/gin-gonic/gin"
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
