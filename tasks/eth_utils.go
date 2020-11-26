package tasks

import (
	"errors"
	"fmt"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"math/big"
	"strconv"
	"strings"
)

func GetBigFromHex(hexAmount string) (*big.Int, error) {
	s := hexAmount
	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	if len(s) == 0 {
		return big.NewInt(0), nil
	}

	n, boole := big.NewInt(0).SetString(s, 16)
	if boole == false {
		err := errors.New("invalid syntax")
		return nil, err
	}

	return n, nil
}

//eth wei -> eth count
func EthFromWei(bigS *big.Int) (strValue string, err error) {
	strValue = "0"
	strN := bigS.String()
	if strN == "0" {
		return
	}

	if len(strN) > EthDecimals {
		strValue = fmt.Sprintf("%v.%v", strN[:len(strN)-EthDecimals], strN[len(strN)-EthDecimals:len(strN)-EthDecimals+Erc20Decimals])
	} else if len(strN) > EthDecimals-Erc20Decimals {
		strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
		strValue = fmt.Sprintf(strFormat, strN[:len(strN)-(EthDecimals-Erc20Decimals)])
	}

	//mylog.DataLogger.Info().Msgf("HexParseEthValue: %v", strValue)
	return
}

//eth wei -> eth count
func FromWeiWithBigintAndToken(bigS *big.Int, token string) (strValue string, err error) {
	strValue = "0"
	var decimals int
	decimals, err = getDecimals(token)
	if err != nil {
		return
	}

	strN := bigS.String()
	if strN == "0" {
		return
	}

	if len(strN) > decimals {
		if decimals > Erc20Decimals {
			strValue = fmt.Sprintf("%v.%v", strN[:len(strN)-decimals], strN[len(strN)-decimals:len(strN)-decimals+Erc20Decimals])
		} else {
			strValue = fmt.Sprintf("%v.%v", strN[:len(strN)-decimals], strN[len(strN)-decimals:len(strN)-decimals+decimals])
		}
	} else {
		if decimals > Erc20Decimals {
			if len(strN) > decimals-Erc20Decimals {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
				strValue = fmt.Sprintf(strFormat, strN[:len(strN)-(decimals-Erc20Decimals)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strValue = fmt.Sprintf(strFormat, strN)
		}
	}

	//mylog.DataLogger.Info().Msgf("HexParseEthValue: %v", strValue)
	return
}

//eth wei -> eth count
func FromWeiWithDecimals(s string, decimals int) (strValue string, err error) {
	strValue = "0"

	if len(s) < 1 {
		err = errors.New("invalid syntax")
		return
	}

	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	if len(s) == 0 {
		return
	}

	n, boole := big.NewInt(0).SetString(s, 16)
	if boole == false {
		err = errors.New("invalid syntax")
		return
	}

	strN := n.String()
	if strN == "0" {
		return
	}

	if len(strN) > decimals {
		if decimals > Erc20Decimals {
			strValue = fmt.Sprintf("%v.%v", strN[:len(strN)-decimals], strN[len(strN)-decimals:len(strN)-decimals+Erc20Decimals])
		} else {
			strValue = fmt.Sprintf("%v.%v", strN[:len(strN)-decimals], strN[len(strN)-decimals:len(strN)-decimals+decimals])
		}
	} else {
		if decimals > Erc20Decimals {
			if len(strN) > decimals-Erc20Decimals {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
				strValue = fmt.Sprintf(strFormat, strN[:len(strN)-(decimals-Erc20Decimals)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strValue = fmt.Sprintf(strFormat, strN)
		}
	}

	return
}

func InputDataIsTransfer(s string) bool {
	if len(s) != (10 + 128) {
		return false
	}

	s1 := s[0:10]
	if s1 != EthTokenTransferHex {
		return false
	}

	return true
}

func ParseTransactionLog(transactionLog *TransactionReceiptLog, decimals int) (isTransfer bool, toAddress string, value string, err error) {
	var methodID string
	if len(transactionLog.Topics) < 3 {
		return
	}

	methodID = transactionLog.Topics[0]
	if methodID != EthLogTransferHex {
		return
	}

	if len(transactionLog.Topics[1]) != (2+64) || len(transactionLog.Topics[2]) != (2+64) {
		return
	}

	toAddress = "0x" + transactionLog.Topics[2][26:]

	//mylog.DataLogger.Info().Msgf("from wei, begin:%v", transactionLog.Data)
	value, err = FromWeiWithDecimals(transactionLog.Data, decimals)
	if err != nil {
		mylog.DataLogger.Error().Msgf("fromWei error, err: %v", err)
		return
	}
	//mylog.DataLogger.Info().Msgf("from wei, end:%v", value)

	isTransfer = true
	return
}

func getDecimals(token string) (decimals int, err error) {
	if token == EthName {
		decimals = EthDecimals
	} else {
		decimals = Erc20Decimals
	}

	return
}

func ToWei(sValue string, decimals int) (*big.Int, error) {
	bigDecimals := big.NewInt(10)
	bigDecimals = bigDecimals.Exp(bigDecimals, big.NewInt(int64(decimals)), nil)

	priceParts := strings.Split(sValue, ".")
	if len(priceParts) == 1 {
		count, err := strconv.ParseUint(sValue, 10, 64)
		if err != nil {
			return nil, err
		}

		iAmount := big.NewInt(int64(count))
		amount := iAmount.Mul(iAmount, bigDecimals)
		return amount, nil
	} else if len(priceParts) == 2 {
		uPrice1, err := strconv.ParseUint(priceParts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		if len(priceParts[1]) <= 0 {
			iAmount := big.NewInt(int64(uPrice1))
			amount := iAmount.Mul(iAmount, bigDecimals)
			return amount, nil
		}

		if len(priceParts[1]) > decimals {
			priceParts[1] = priceParts[1][:decimals]
		}
		uPrice2, err2 := strconv.ParseUint(priceParts[1], 10, 64)
		if err2 != nil {
			err = err2
			return nil, err
		}

		iAmount1 := big.NewInt(int64(uPrice1))
		iAmount1.Mul(iAmount1, bigDecimals)
		iAmount2 := big.NewInt(int64(uPrice2))
		iAmount2.Mul(iAmount2, bigDecimals.Exp(big.NewInt(10), big.NewInt(int64(decimals-len(priceParts[1]))), nil))
		amount := iAmount1.Add(iAmount1, iAmount2)
		return amount, nil
	}

	return nil, errors.New("input invalid")
}
