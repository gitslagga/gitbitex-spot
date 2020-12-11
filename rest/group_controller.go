package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
)

// Get /group/info
func GroupInfoService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	coin := ctx.Query("coin")
	before, _ := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, _ := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

	groups, err := service.GetGroupByCoin(coin, before, after, limit)
	if err != nil {
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

// Get /group/log
func GroupLogService(ctx *gin.Context) {
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

	groupLogs, err := service.GetGroupLogByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(groupLogs) > 0 {
		newBefore = groupLogs[0].Id
		newAfter = groupLogs[len(groupLogs)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   groupLogs,
	}
	ctx.JSON(http.StatusOK, out)
}

// Get /group/publicity
func GroupPublicityService(ctx *gin.Context) {
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

	groupLogs, err := service.GetGroupLogPublicity(before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(groupLogs) > 0 {
		newBefore = groupLogs[0].Id
		newAfter = groupLogs[len(groupLogs)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   groupLogs,
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /group/publish
func GroupPublishService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var groupRequest GroupRequest
	err := ctx.ShouldBindJSON(&groupRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] GroupPublishService request param: %v", groupRequest)

	group, err := service.GetGroupByUserIdCoin(address.Id, groupRequest.Coin)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if group != nil {
		out.RespCode = EC_GROUP_PUBLISH_EXISTS
		out.RespDesc = ErrorCodeMessage(EC_GROUP_PUBLISH_EXISTS)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.GroupPublish(address, groupRequest.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] GroupPublishService GroupPublish err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /group/join
func GroupJoinService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var groupJoinRequest GroupJoinRequest
	err := ctx.ShouldBindJSON(&groupJoinRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] GroupJoinService request param: %v", groupJoinRequest)

	group, err := service.GetGroupById(groupJoinRequest.GroupId)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if group == nil {
		out.RespCode = EC_GROUP_JOIN_NOT_EXISTS
		out.RespDesc = ErrorCodeMessage(EC_GROUP_JOIN_NOT_EXISTS)
		ctx.JSON(http.StatusOK, out)
		return
	}

	groupLog, err := service.GetGroupLogByGroupIdUserId(groupJoinRequest.GroupId, address.Id)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if groupLog != nil {
		out.RespCode = EC_GROUP_JOIN_REPEAT_ERR
		out.RespDesc = ErrorCodeMessage(EC_GROUP_JOIN_REPEAT_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.GroupJoin(address, group)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] GroupJoinService GroupJoin err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /group/delegate
func GroupDelegateService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var groupRequest GroupRequest
	err := ctx.ShouldBindJSON(&groupRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] GroupDelegateService request param: %v", groupRequest)

	err = service.GroupDelegate(address, groupRequest.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] GroupDelegateService GroupDelegate err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /group/release
func GroupReleaseService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var groupRequest GroupRequest
	err := ctx.ShouldBindJSON(&groupRequest)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] AddressReleaseService request param: %v", groupRequest)

	err = service.GroupRelease(address, groupRequest.Coin)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] AddressReleaseService GroupRelease err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
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

	before, _ := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, _ := strconv.ParseInt(ctx.Query("after"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

	releases, err := service.GetAddressReleaseByUserId(address.Id, before, after, limit)
	if err != nil {
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
