package tasks

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/utils"
	"math/big"
	"net/http"
	"strconv"
)

var StatusCodeError = errors.New("status code error")

const (
	BaseUrl         = "https://mainnet.infura.io/v3/eae18f8a2d38404f96848b6c5b64bcbd"
	ChainID         = 1
	EthMainAddress  = "0x1C49b8CCee62b15F750C8D96e0258fB09B109F50"
	EthMainMnemonic = "poem police guide flip drip scout clutch now surround vacuum share page carpet demand alone"
	EthMainPrivate  = "3e55c9310f04d55c74dbe0911951845c9fbe0fca5768fe696ba38a5d5af304b4"

	UsdtName            = "USDT"
	UsdtContractAddress = "0xdac17f958d2ee523a2206206994597c13d831ec7"
)

const (
	EthName                 = "ETH"
	ERC20                   = "ERC20"
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
	EthGasPrice   = 500000000000
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
func EthGetBalance(address string, tag string) (string, error) {

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

//--sign transaction
func EthPersonalSendTransactionToBlockWithFee(fromAddress string, password string, toAddress string, bigAmount *big.Int) (string, error) {
	//减去eth转币需要的费用
	bigFee := big.NewInt(0).Mul(big.NewInt(EthGasPrice), big.NewInt(EthGasLimit))
	subFeeAmount := bigAmount.Sub(bigAmount, bigFee)
	if subFeeAmount.Int64() < 0 {
		mylog.Logger.Error().Msgf("[ETH_PersonalSendTransactionToBlockWithFee] Sub fee amount < 0")
		return "", errors.New("sub fee amount < 0")
	}

	nonce, err := EthGetTransactionCount(fromAddress)
	if err != nil {
		mylog.Logger.Error().Msgf("[ETH_PersonalSendTransactionToBlockWithFee] EthGetTransactionCount err:%v", err)
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), subFeeAmount, uint64(EthGasLimit), big.NewInt(EthGasPrice), nil)
	raw, err := SignTransaction(ChainID, tx, password)
	if err != nil {
		mylog.Logger.Error().Msgf("[ETH_PersonalSendTransactionToBlockWithFee] SignTransaction err:%v", err)
		return "", err
	}

	txid, err := EthSendRawTransaction(raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func Erc20PersonalSendTransactionToBlock(fromAddress string, password string, toAddress string, token string, bigAmount *big.Int) (string, error) {
	nonce, err := EthGetTransactionCount(fromAddress)
	if err != nil {
		mylog.Logger.Error().Msgf("[ERC20_PersonalSendTransactionToBlock] EthGetTransactionCount err:%v", err)
		return "", err
	}

	config, err := mysql.SharedStore().GetAddressConfigByCoin(token)
	if err != nil {
		mylog.Logger.Error().Msgf("[ERC20_PersonalSendTransactionToBlock] SelectContractAddressFromMysql error, err:%v, token: %v", err, token)
		return "", err
	}

	data, err := MakeERC20TransferData(toAddress, bigAmount)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(config.ContractAddress), big.NewInt(0), uint64(Erc20GasLimit), big.NewInt(EthGasPrice), data)
	raw, err := SignTransaction(ChainID, tx, password)
	if err != nil {
		mylog.Logger.Error().Msgf("[ERC20_PersonalSendTransactionToBlock] SignTransaction err:%v", err)
		return "", err
	}

	txId, err := EthSendRawTransaction(raw)
	if err != nil {
		return "", err
	}

	return txId, nil
}

func PersonalSendTransactionToBlock(fromAddress, toAddress, token, amount string) (string, error) {
	raw, err := PersonalSignTransaction(fromAddress, toAddress, token, amount)
	if err != nil {
		mylog.Logger.Error().Msgf("[PersonalSendTransactionToBlock] PersonalSignTransaction err:%v", err)
		return "", err
	}

	txid, err := EthSendRawTransaction(raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func checkMainAddressTokenIsEnough(amount, contractAddress, token string, decimals int) bool {
	mainHexBalance, err := TokenGetBalance(EthMainAddress, contractAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[checkMainAddressTokenIsEnough] TokenGetBalance err:%v", err)
		return false
	}

	//判断main address token是否充足
	sMainBalance, err := FromWeiWithDecimals(mainHexBalance, decimals)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[checkMainAddressTokenIsEnough] FromWeiWithDecimals err:%v", err)
		return false
	}

	fMainBalance, err := strconv.ParseFloat(sMainBalance, 64)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[checkMainAddressTokenIsEnough] ParseFloat err:%v", err)
		return false
	}
	fAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[checkMainAddressTokenIsEnough] ParseFloat err:%v", err)
		return false
	}

	mylog.DataLogger.Info().Msgf("[checkMainAddressTokenIsEnough] %s main balance:%s, amount:%s", token, sMainBalance, amount)

	if fMainBalance < fAmount {
		mylog.DataLogger.Warn().Msgf("[checkMainAddressTokenIsEnough] main balance not enough, token:%s, main:%s, amount:%s", token, sMainBalance, amount)
		return false
	}

	return true
}

func PersonalSignTransaction(fromAddress, toAddress, coin, amount string) (string, error) {
	nonce, err := EthGetTransactionCount(fromAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[PersonalSignTransaction] EthGetTransactionCount err:%v", err)
		return "", err
	}

	if coin == EthName {
		hexAmount, err := ToWei(amount, EthDecimals)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[PersonalSignTransaction] getAmountFromPrice err:%v, value: %v", err, amount)
			return "", err
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), hexAmount, uint64(EthGasLimit), big.NewInt(EthGasPrice), nil)
		return SignTransaction(ChainID, tx, EthMainPrivate)
	}

	config, err := mysql.SharedStore().GetAddressConfigByCoin(coin)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[PersonalSignTransaction] SelectTokenInfoFromDB error, err:%v, coin: %v", err, coin)
		return "", err
	}

	isEnough := checkMainAddressTokenIsEnough(amount, config.ContractAddress, coin, config.Decimals)
	mylog.DataLogger.Info().Msgf("[PersonalSignTransaction] checkMainAddressTokenIsEnough isEnougn:%v", isEnough)
	if isEnough == false {
		err = errors.New(fmt.Sprintf("main balance not enough, coin:%s", coin))
		return "", err
	}

	hexAmount, err := ToWei(amount, config.Decimals)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[PersonalSignTransaction] getAmountFromPrice err:%v, value: %v", err, amount)
		return "", err
	}

	data, err := MakeERC20TransferData(toAddress, hexAmount)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(config.ContractAddress), big.NewInt(0), uint64(Erc20GasLimit), big.NewInt(EthGasPrice), data)
	return SignTransaction(ChainID, tx, EthMainPrivate)
}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignTransaction(chainID int64, tx *types.Transaction, privateKeyStr string) (string, error) {
	privateKey, err := StringToPrivateKey(fmt.Sprintf("0x%v", privateKeyStr))
	if err != nil {
		return "", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%v", hex.EncodeToString(b)), nil
}

func MakeERC20TransferData(toAddress string, amount *big.Int) ([]byte, error) {
	methodId := crypto.Keccak256Hash([]byte("transfer(address,uint256)"))
	var data []byte
	data = append(data, methodId[:4]...)
	paddedAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...)
	return data, nil
}
