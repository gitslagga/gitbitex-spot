package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

const (
	keyCurrentUser    = "__current_user"
	keyCurrentAddress = "__current_address"
)

func checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if len(token) == 0 {
			var err error
			token, err = c.Cookie("accessToken")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, newMessageVo(errors.New("token not found")))
				return
			}
		}

		user, err := service.CheckToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, newMessageVo(err))
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, newMessageVo(errors.New("bad token")))
			return
		}

		c.Set(keyCurrentUser, user)
		c.Next()
	}
}

func GetCurrentUser(ctx *gin.Context) *models.User {
	val, found := ctx.Get(keyCurrentUser)
	if !found {
		return nil
	}
	return val.(*models.User)
}

//development new
func checkFrontendToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		out := CommonResp{}
		token := c.GetHeader("token")
		if len(token) == 0 {
			var err error
			token, err = c.Cookie("accessToken")
			if err != nil {
				out.RespCode = EC_NETWORK_ERR
				out.RespDesc = err.Error()
				c.AbortWithStatusJSON(http.StatusOK, out)
				return
			}
		}

		address, err := service.CheckFrontendToken(token)
		if err != nil {
			out.RespCode = EC_NETWORK_ERR
			out.RespDesc = err.Error()
			c.AbortWithStatusJSON(http.StatusOK, out)
			return
		}

		if address == nil {
			out.RespCode = EC_TOKEN_INVALID
			out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
			c.AbortWithStatusJSON(http.StatusOK, out)
			return
		}

		c.Set(keyCurrentAddress, address)
		c.Next()
	}
}

func GetCurrentAddress(ctx *gin.Context) *models.Address {
	val, found := ctx.Get(keyCurrentAddress)
	if !found {
		return nil
	}
	return val.(*models.Address)
}

func BackendAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		out := CommonResp{}

		// 定义ip白名单
		whiteList := map[string]bool{
			"127.0.0.1": true,
		}

		if _, ok := whiteList[c.ClientIP()]; !ok {
			out.RespCode = EC_WHITE_LIST_ERR
			out.RespDesc = ErrorCodeMessage(EC_WHITE_LIST_ERR)
			c.AbortWithStatusJSON(http.StatusOK, out)
			return
		}

		c.Next()
	}
}
