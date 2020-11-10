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

	frontend := r.Group("/api", checkFrontendToken())
	{
		frontend.GET("/address/info", AddressService)
		frontend.DELETE("/address/logout", LogoutService)
		frontend.POST("/address/findPassword", FindPasswordService)
		frontend.POST("/address/modifyPassword", ModifyPasswordService)
		frontend.POST("/address/activation", ActivationService)
	}

	r.GET("/api/config/info", GetConfigService)
	r.POST("/api/admin/login", AdminLoginService)
	backend := r.Group("/api", checkBackendToken())
	{
		backend.GET("/admin/info", AdminService)
		backend.POST("/config/update", UpdateConfigService)
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
