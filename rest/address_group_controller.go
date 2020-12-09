package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
)

// Get /address/groupInfo
func AddressGroupInfoService(ctx *gin.Context) {
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

	groups, err := service.GetAddressGroupByUserId(address.Id, before, after, limit)
	if groups == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(groups) > 0 {
		newBefore = groups[0].Id
		newAfter = groups[len(groups)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   groups,
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /address/group
func AddressGroupService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var addressGroup AddressGroupRequest
	err := ctx.ShouldBindJSON(&addressGroup)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressGroupService request param: %v", addressGroup)

	err = service.AddressGroup(address, addressGroup.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressGroupService AddressGroup err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /address/delegate
func AddressDelegateService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var addressGroup AddressGroupRequest
	err := ctx.ShouldBindJSON(&addressGroup)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressDelegateService request param: %v", addressGroup)

	err = service.AddressDelegate(address, addressGroup.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressDelegateService AddressDelegate err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /address/release
func AddressReleaseService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var addressGroup AddressGroupRequest
	err := ctx.ShouldBindJSON(&addressGroup)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressReleaseService request param: %v", addressGroup)

	err = service.AddressRelease(address, addressGroup.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressReleaseService AddressRelease err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// Get /address/taskRelease
func AddressTaskReleaseService(ctx *gin.Context) {
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

	releases, err := service.GetAddressReleaseByUserId(address.Id, before, after, limit)
	if releases == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(releases) > 0 {
		newBefore = releases[0].Id
		newAfter = releases[len(releases)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   releases,
	}
	ctx.JSON(http.StatusOK, out)
}
