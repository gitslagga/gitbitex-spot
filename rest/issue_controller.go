package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// Get /backend/issue/list
func BackendIssueListService(ctx *gin.Context) {
	out := CommonResp{}

	list, err := service.BackendIssueList()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = list

	ctx.JSON(http.StatusOK, out)
}

// Post /backend/issue/release
func BackendIssueReleaseService(ctx *gin.Context) {
	out := CommonResp{}

	err := service.BackendIssueRelease()
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = err.Error()
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()

	ctx.JSON(http.StatusOK, out)
}
