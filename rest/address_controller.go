// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/service"
	"net/http"
)

// POST /api/address/mnemonic
func MnemonicService(ctx *gin.Context) {
	out := CommonResp{}

	mnemonic, err := service.CreateMnemonic()
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] MnemonicService CreateAddress err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = mnemonic

	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/register
func RegisterService(ctx *gin.Context) {
	out := CommonResp{}
	var register RegisterRequest
	err := ctx.ShouldBindJSON(&register)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if len(register.Password) < 8 {
		out.RespCode = EC_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] RegisterService request param: %v", register)

	_, err = service.CreateAddress(register.Username, encryptPassword(register.Password), register.Mnemonic)
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] RegisterService CreateAddress err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/login
func LoginService(ctx *gin.Context) {
	out := CommonResp{}
	var login LoginRequest
	err := ctx.ShouldBindJSON(&login)
	if err != nil || (login.Mnemonic == "" && login.PrivateKey == "") {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	if len(login.Password) < 8 {
		out.RespCode = EC_PASSWORD_ERR
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] LoginService request param: %v", login)

	address, err := service.UpdateAddress(login.Mnemonic, login.PrivateKey, encryptPassword(login.Password))
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] RegisterService CreateAddress err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	token, err := service.CreateJwtToken(address)
	if err != nil {
		mylog.DataLogger.Info().Msgf("[Rest] LoginService CreateJwtToken err: %v", err)
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	ctx.SetCookie("accessToken", token, 7*24*60*60, "/", "", false, false)

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = token
	ctx.JSON(http.StatusOK, out)
}

// DELETE /api/address/logout
func LogoutService(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", -1, "/", "", false, false)

	out := CommonResp{
		RespCode: EC_NONE.Code(),
		RespDesc: EC_NONE.String(),
	}
	ctx.JSON(http.StatusOK, out)
}

// GET /api/address/info
func AddressService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	out.RespData = address
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/findPassword
func FindPasswordService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var findPassword FindPasswordRequest
	err := ctx.ShouldBindJSON(&findPassword)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] FindPasswordService request param: %v", findPassword)

	if address.PrivateKey != findPassword.PrivateKey {
		out.RespCode = EC_MNEMONIC_INCORRECT
		out.RespDesc = ErrorCodeMessage(EC_MNEMONIC_INCORRECT)
		ctx.JSON(http.StatusOK, out)
		return
	}

	address.Password = encryptPassword(findPassword.Password)
	err = service.UpdateAddressByAddr(address)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

// POST /api/address/modifyPassword
func ModifyPasswordService(ctx *gin.Context) {
	out := CommonResp{}
	address := GetCurrentAddress(ctx)
	if address == nil {
		out.RespCode = EC_TOKEN_INVALID
		out.RespDesc = ErrorCodeMessage(EC_TOKEN_INVALID)
		ctx.JSON(http.StatusOK, out)
		return
	}

	var modifyPassword ModifyPasswordRequest
	err := ctx.ShouldBindJSON(&modifyPassword)
	if err != nil {
		out.RespCode = EC_PARAMS_ERR
		out.RespDesc = ErrorCodeMessage(EC_PARAMS_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	mylog.Logger.Info().Msgf("[Rest] ModifyPasswordService request param: %v", modifyPassword)

	if address.Password != encryptPassword(modifyPassword.OldPassword) {
		out.RespCode = EC_PASSWORD_INCORRECT
		out.RespDesc = ErrorCodeMessage(EC_PASSWORD_INCORRECT)
		ctx.JSON(http.StatusOK, out)
		return
	}

	address.Password = encryptPassword(modifyPassword.NewPassword)
	err = service.UpdateAddressByAddr(address)
	if err != nil {
		out.RespCode = EC_NETWORK_ERR
		out.RespDesc = ErrorCodeMessage(EC_NETWORK_ERR)
		ctx.JSON(http.StatusOK, out)
		return
	}

	out.RespCode = EC_NONE.Code()
	out.RespDesc = EC_NONE.String()
	ctx.JSON(http.StatusOK, out)
}

func encryptPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}
