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

	r.POST("/api/user/signup", SignUp)
	r.POST("/api/user/signin", SignIn)
	r.POST("/api/user/token", GetToken)
	r.GET("/api/product", GetProducts)
	r.GET("/api/product/:productId/trades", GetProductTrades)
	r.GET("/api/product/:productId/book", GetProductOrderBook)
	r.GET("/api/product/:productId/candles", GetProductCandles)

	private := r.Group("/api", checkToken())
	{
		private.GET("/api/configs", GetConfigs)
		private.POST("/api/config", UpdateConfig)

		private.GET("/order", GetOrders)
		private.POST("/order", PlaceOrder)
		private.DELETE("/order/:orderId", CancelOrder)
		private.DELETE("/order", CancelOrders)
		private.GET("/account", GetAccounts)
		private.GET("/user/self", GetUsersSelf)
		private.POST("/user/password", ChangePassword)
		private.DELETE("/user/accessToken", SignOut)

		private.GET("/wallet/:currency/address", GetWalletAddress)
		private.GET("/wallet/:currency/transactions", GetWalletTransactions)
		private.POST("/wallet/:currency/withdrawal", Withdrawal)
	}

	//development new
	r.POST("/api/address/mnemonic", MnemonicService)
	r.POST("/api/address/register", RegisterService)
	r.POST("/api/address/login", LoginService)

	personal := r.Group("/api", checkJwtToken())
	{
		personal.GET("/address/info", AddressService)
		personal.DELETE("/address/logout", LogoutService)
		personal.POST("/address/findPassword", FindPasswordService)
		personal.POST("/address/modifyPassword", ModifyPasswordService)
		personal.POST("/address/activation", ActivationService)
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
