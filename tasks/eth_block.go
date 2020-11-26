package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitslagga/gitbitex-spot/utils"
	"math/big"
	"net/http"
)

var StatusCodeError = errors.New("status code error")

const (
	BaseUrl                 = "https://mainnet.infura.io/v3/eae18f8a2d38404f96848b6c5b64bcbd"
	EthName                 = "eth"
	ERC20                   = "erc20"
	EthTokenTransferHex     = "0xa9059cbb"
	EthTokenTransferFromHex = "0x23b872dd"
	EthTokenDecimalsHex     = "0x313ce567"
	EthTokenBalanceOfHex    = "0x70a08231"
	EthTokenTotalSupplyHex  = "0x18160ddd"

	EthLogTransferHex = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" //eth log,中transfer id
)

const (
	EthDecimals   = 18
	Erc20Decimals = 8
	EthGasLimit   = 40000
	Erc20GasLimit = 90000
	//ETH_GAS_PRICE   = 2e10
)

func getEthJsonStr(method string, params string) (string, error) {
	var jsonStr string

	if len(method) == 0 {
		return jsonStr, errors.New("method is null")
	}

	jsonStr = fmt.Sprintf(`{"jsonrpc": "2.0", "id":"exchange", "method": "%v", "params": [%v]}`, method, params)

	return jsonStr, nil
}

func EthBlockNumber() (uint64, error) {
	var blockHeight uint64

	jsonStr, err := getEthJsonStr("eth_blockNumber", "")
	if err != nil {
		return blockHeight, err
	}

	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("eth_blockNumber, respBody=%v", string(respBody))
	if statusCode == 200 {
		var resp BlockCountResp
		err = json.Unmarshal(respBody, &resp)
		if err != nil {
			return blockHeight, err
		}

		if len(resp.Result) <= 2 {
			return blockHeight, errors.New("result error")
		}

		blockNumber := big.NewInt(0)
		block, boole := blockNumber.SetString(resp.Result[2:], 16)
		if boole == false {
			return blockHeight, errors.New("big.setString error")
		}

		blockHeight = block.Uint64()
	} else {
		return blockHeight, StatusCodeError
	}

	return blockHeight, nil
}

func GetBlockByNumber(curblock uint64) ([]*Transaction, error) {

	var err error

	var jsonStr = fmt.Sprintf(`{"jsonrpc": "2.0", "id":"exchange", "method": "eth_getBlockByNumber", "params": ["0x%x",true]}`, curblock)
	//mylog.DataLogger.Info().Msgf("getBlockByNumber, block: %v", curblock)
	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	if statusCode != 200 {
		return nil, StatusCodeError
	}
	var resp BlockTrueResp
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	//mylog.DataLogger.Info().Msgf("block:%d, transactions count:%d", curblock, len(resp.Result.Transactions))

	return resp.Result.Transactions, nil
}

func GetTransactionReceipt(txid string) (*RowTransactionReceipt, error) {
	var err error

	var jsonStr = fmt.Sprintf(`{"jsonrpc": "2.0", "id":"exchange", "method": "eth_getTransactionReceipt", "params": ["%v"]}`, txid)
	//mylog.DataLogger.Info().Msgf("eth_getTransactionReceipt, txid: %v", txid)
	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("eth_getTransactionReceipt, respBody: %v", string(respBody))
	if statusCode != 200 {
		return nil, StatusCodeError
	}
	var resp RowTransactionReceipt
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

//tag: the default block parameter, HEX string, or string "earliest", or string "latest", or string "pending"
func GetBalance(address string, tag string) (string, error) {

	var amount string = "0"

	sparams := fmt.Sprintf(`"%v", "%v"`, address, tag)
	jsonStr, err := getEthJsonStr("eth_getBalance", sparams)
	if err != nil {
		return amount, err
	}

	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("eth_getBalance, respBody=%v", string(respBody))
	if statusCode == 200 {
		var resp TokenBalanceResp
		err = json.Unmarshal(respBody, &resp)
		if err != nil {
			return amount, err
		}

		if len(resp.Amount) > 0 {
			amount = resp.Amount
		}

	} else {
		return amount, errors.New("eth_getBalance error")
	}

	return amount, nil
}

func TokenGetBalance(address string, contractAddress string) (string, error) {
	var amount string = "0"

	jsondata := fmt.Sprintf("%v%064s", EthTokenBalanceOfHex, address[2:])

	var jsonStr = fmt.Sprintf(`{"jsonrpc": "2.0", "id":"exchange", "method": "eth_call", "params": [{"to":"%v", "data":"%v"},"latest"]}`, contractAddress, jsondata)
	//mylog.DataLogger.Info().Msgf("token_getBalance, reqbody=%v", jsonStr)
	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("token_getBalance, respBody=%v", string(respBody))

	if statusCode == 200 {
		var resp TokenBalanceResp
		err := json.Unmarshal(respBody, &resp)
		if err != nil {
			return amount, err
		}

		if len(resp.Amount) > 0 {
			amount = resp.Amount
		}
	} else {
		return amount, errors.New("eth_getBalance error")
	}

	return amount, nil
}

func EthGetTransactionCount(address string) (uint64, error) {
	var count uint64

	parasStr := fmt.Sprintf(`"%v", "pending"`, address)

	jsonStr, err := getEthJsonStr("eth_getTransactionCount", parasStr)
	if err != nil {
		return count, err
	}

	//mylog.DataLogger.Info().Msgf("eth_getTransactionCount, jsonStr:%v", jsonStr)
	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("eth_getTransactionCount, respBody=%v", string(respBody))
	if statusCode == 200 {
		var resp EthTransactionCountResp
		err = json.Unmarshal(respBody, &resp)
		if err != nil {
			return count, err
		}

		if len(resp.Error) > 0 {
			return count, errors.New(resp.Error)
		}

		if len(resp.Result) <= 2 {
			return count, errors.New("result error")
		}

		tmp := big.NewInt(0)
		bigCount, boole := tmp.SetString(resp.Result[2:], 16)
		if boole == false {
			return count, errors.New("big.setString error")
		}

		count = bigCount.Uint64()
	} else {
		return count, StatusCodeError
	}

	return count, nil
}

func EthSendRawTransaction(raw string) (string, error) {
	var txid string

	parasStr := fmt.Sprintf(`"%v"`, raw)

	jsonStr, err := getEthJsonStr("eth_sendRawTransaction", parasStr)
	if err != nil {
		return txid, err
	}

	respBody, _, statusCode := utils.SharedProxy().PostJson("", BaseUrl, []byte(jsonStr), func(*http.Request) {})
	//mylog.DataLogger.Info().Msgf("eth_sendRawTransaction, respBody=%v", string(respBody))
	if statusCode == 200 {
		var resp SendRawTransactionResp
		err = json.Unmarshal(respBody, &resp)
		if err != nil {
			return txid, err
		}

		if resp.Error != nil {
			return txid, errors.New(resp.Error.Message)
		}

		txid = resp.Txid
	} else {
		return txid, StatusCodeError
	}

	return txid, nil
}
