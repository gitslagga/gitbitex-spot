package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/service"
	"github.com/gitslagga/gitbitex-spot/utils"
	"net/http"
	"strconv"
)

// GET /products
func GetProducts(ctx *gin.Context) {
	products, err := service.GetProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	var productVos []*ProductVo
	for _, product := range products {
		productVos = append(productVos, newProductVo(product))
	}

	ctx.JSON(http.StatusOK, productVos)
}

// GET /products/<product-id>/trades
func GetProductTrades(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var tradeVos []*tradeVo
	trades, _ := service.GetTradesByProductId(productId, 50)
	for _, trade := range trades {
		tradeVos = append(tradeVos, newTradeVo(trade))
	}

	ctx.JSON(http.StatusOK, tradeVos)
}

// GET /products/<product-id>/candles
func GetProductCandles(ctx *gin.Context) {
	productId := ctx.Param("productId")
	granularity, _ := utils.AToInt64(ctx.Query("granularity"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "1000"))
	if limit <= 0 || limit > 10000 {
		limit = 1000
	}

	//[
	//    [ time, low, high, open, close, volume ],
	//    [ 1415398768, 0.32, 4.2, 0.35, 4.2, 12.3 ],
	//]
	var tickVos [][6]float64
	ticks, err := service.GetTicksByProductId(productId, granularity/60, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}
	for _, tick := range ticks {
		tickVos = append(tickVos, [6]float64{float64(tick.Time), utils.DToF64(tick.Low), utils.DToF64(tick.High),
			utils.DToF64(tick.Open), utils.DToF64(tick.Close), utils.DToF64(tick.Volume)})
	}

	ctx.JSON(http.StatusOK, tickVos)
}

// GET /product/info
func GetProductService(ctx *gin.Context) {
	out := CommonResp{}
	products, err := service.GetProducts()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = products
	ctx.JSON(http.StatusOK, out)
}

// GET /product/trade/:productId
func GetProductTradeService(ctx *gin.Context) {
	out := CommonResp{}

	productId := ctx.Param("productId")
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	trades, err := service.GetTradesByProductId(productId, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = trades
	ctx.JSON(http.StatusOK, out)
}

// GET /candle/:productId
func GetProductCandleService(ctx *gin.Context) {
	out := CommonResp{}

	productId := ctx.Param("productId")
	granularity, err1 := utils.AToInt64(ctx.DefaultQuery("granularity", "60"))
	limit, err2 := strconv.Atoi(ctx.DefaultQuery("limit", "1000"))
	if err1 != nil || err2 != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}
	if limit <= 0 || limit > 10000 {
		limit = 1000
	}

	ticks, err := service.GetTicksByProductId(productId, granularity/60, limit)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = ticks
	ctx.JSON(http.StatusOK, out)
}
