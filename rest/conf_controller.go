package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// GET /configs
func GetConfigs(ctx *gin.Context) {
	configs, err := service.GetConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	m := map[string]string{}
	for _, config := range configs {
		m[config.Key] = config.Value
	}

	ctx.JSON(http.StatusOK, m)
}

// Get /config/info
func GetConfigService(ctx *gin.Context) {
	out := CommonResp{}
	config, err := service.GetConfigs()
	if err == nil {
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

// POST /config/update
func UpdateConfigService(ctx *gin.Context) {
	out := CommonResp{}
	var updateConfig UpdateConfigRequest
	err := ctx.ShouldBindJSON(&updateConfig)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] UpdateConfigService request param: %v", updateConfig)

	config, err := service.GetConfigById(updateConfig.Id)
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] UpdateConfigService GetConfigById err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if config == nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	config.Value = updateConfig.Value
	err = service.UpdateConfig(config)

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}
