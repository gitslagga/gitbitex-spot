package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
	"strconv"
	"strings"
)

// GET /api/machine/info
func GetMachineService(ctx *gin.Context) {
	out := CommonResp{}

	machine, err := service.GetBuyMachine()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = machine

	ctx.JSON(http.StatusOK, out)
}

// POST /api/machine/buy
func BuyMachineService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var buyMachine BuyMachineRequest
	err := ctx.ShouldBindJSON(&buyMachine)
	if err != nil || buyMachine.MachineId == models.MachineGiveAwayId {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] BuyMachineService request param: %v", buyMachine)

	machine, err := service.GetMachineById(buyMachine.MachineId)
	if machine == nil || err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	err = service.BuyMachine(address, machine, strings.ToUpper(buyMachine.Currency))
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] BuyMachineService BuyMachine err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	go service.StartMachineLevel(address)

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// GET /api/machine/address
func AddressMachineService(ctx *gin.Context) {
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

	machineAddress, err := service.GetMachineAddressByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(machineAddress) > 0 {
		newBefore = machineAddress[0].Id
		newAfter = machineAddress[len(machineAddress)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   machineAddress,
	}
	ctx.JSON(http.StatusOK, out)
}

// GET /api/machine/log
func LogMachineService(ctx *gin.Context) {
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

	machineLog, err := service.GetMachineLogByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(machineLog) > 0 {
		newBefore = machineLog[0].Id
		newAfter = machineLog[len(machineLog)-1].Id
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   machineLog,
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /machine/convert
func MachineConvertService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var machineConvert MachineConvertRequest
	err := ctx.ShouldBindJSON(&machineConvert)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] MachineConvertService request param: %v", machineConvert)

	// ConvertType: 1-ytl兑换bite, 2-bite兑换ytl
	err = service.MachineConvert(address, machineConvert.ConvertType, machineConvert.Number)
	if err != nil {
		mylog.Logger.Error().Msgf("[Rest] MachineConvertService MachineConvert err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}

// GET /machine/convertInfo
func MachineConvertInfoService(ctx *gin.Context) {
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

	machineConvert, err := service.GetMachineConvertByUserId(address.Id, before, after, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var newBefore, newAfter int64 = 0, 0
	if len(machineConvert) > 0 {
		newBefore = machineConvert[0].Id
		newAfter = machineConvert[len(machineConvert)-1].Id
	}
	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = PageResp{
		Before: newBefore,
		After:  newAfter,
		List:   machineConvert,
	}
	ctx.JSON(http.StatusOK, out)
}

// GET /machine/level
func MachineLevelService(ctx *gin.Context) {
	out := CommonResp{}

	machineLevel, err := service.GetMachineLevel()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = machineLevel
	ctx.JSON(http.StatusOK, out)
}

// GET /machine/config
func MachineConfigService(ctx *gin.Context) {
	out := CommonResp{}

	config, err := service.GetMachineConfigs()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = config
	ctx.JSON(http.StatusOK, out)
}
