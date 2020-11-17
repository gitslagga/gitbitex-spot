package rest

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func (server *HttpServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.Use(setCROSOptions)

	r.GET("/api/configs", GetConfigs)
	r.POST("/api/users", SignUp)
	r.POST("/api/users/accessToken", SignIn)
	r.POST("/api/users/token", GetToken)
	r.GET("/api/products", GetProducts)
	r.GET("/api/products/:productId/trades", GetProductTrades)
	r.GET("/api/products/:productId/book", GetProductOrderBook)
	r.GET("/api/products/:productId/candles", GetProductCandles)

	private := r.Group("/", checkToken())
	{
		private.GET("/api/orders", GetOrders)
		private.POST("/api/orders", PlaceOrder)
		private.DELETE("/api/orders/:orderId", CancelOrder)
		private.DELETE("/api/orders", CancelOrders)
		private.GET("/api/accounts", GetAccounts)
		private.GET("/api/users/self", GetUsersSelf)
		private.POST("/api/users/password", ChangePassword)
		private.DELETE("/api/users/accessToken", SignOut)
		private.GET("/api/wallets/:currency/address", GetWalletAddress)
		private.GET("/api/wallets/:currency/transactions", GetWalletTransactions)
		private.POST("/api/wallets/:currency/withdrawal", Withdrawal)
	}

	//development new
	r.POST("/api/address/mnemonic", MnemonicService)
	r.POST("/api/address/register", RegisterService)
	r.POST("/api/address/login", LoginService)

	r.GET("/api/product/info", GetProductService)
	r.GET("/api/trade/:productId", GetProductTradeService)
	r.GET("/api/book/:productId", GetProductOrderBookService)
	r.GET("/api/candle/:productId", GetProductCandleService)

	frontend := r.Group("/api", checkFrontendToken())
	{
		frontend.GET("/api/config/info", GetConfigService)
		frontend.GET("/address/info", AddressService)
		frontend.DELETE("/address/logout", LogoutService)
		frontend.POST("/address/findPassword", FindPasswordService)
		frontend.POST("/address/modifyPassword", ModifyPasswordService)
		frontend.POST("/address/activation", ActivationService)

		frontend.GET("/order/info", GetOrderService)
		frontend.POST("/order/place", PlaceOrderService)
		frontend.DELETE("/order/cancel/:orderId", CancelOrderService)
		frontend.DELETE("/order/cancelAll", CancelAllOrderService)

		frontend.GET("/machine/info", GetMachineService)
		frontend.POST("/machine/buy", BuyMachineService)
		frontend.GET("/machine/address", AddressMachineService)
		frontend.GET("/machine/log", LogMachineService)

		frontend.GET("/account/address", AccountAddressService)
		frontend.POST("/account/convert", AccountConvertService)
		frontend.GET("/account/convertInfo", AccountConvertInfoService)
		//frontend.POST("/account/transfer", AccountTransferService)
		//frontend.POST("/account/scan", AccountScanService)

	}

	err := r.Run(server.addr)
	if err != nil {
		panic(err)
	}
}

func setCROSOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")
}
