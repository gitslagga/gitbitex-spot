package tasks

import (
	"errors"
	"fmt"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"math/big"
	"strconv"
	"strings"
)

func GetBigFromHex(hexamount string) (*big.Int, error) {
	s := hexamount
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
func EthFromWei(bigs *big.Int) (strvalue string, err error) {

	strvalue = "0"
	decimals := EthDecimals

	strn := bigs.String()
	if strn == "0" {
		return
	}

	if len(strn) > decimals {
		if decimals > Erc20Decimals {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+Erc20Decimals])
		} else {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+decimals])
		}
	} else {
		if decimals > Erc20Decimals {
			if len(strn) > decimals-Erc20Decimals {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
				strvalue = fmt.Sprintf(strFormat, strn[:len(strn)-(decimals-Erc20Decimals)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strvalue = fmt.Sprintf(strFormat, strn)
		}
	}

	//mylog.DataLogger.Info().Msgf("HexParseEthvalue: %v", strvalue)
	return
}

//eth wei -> eth count
func FromWeiWithBigintAndToken(bigs *big.Int, token string) (strvalue string, err error) {

	strvalue = "0"
	var decimals int
	decimals, err = getDecimals(token)
	if err != nil {
		return
	}

	strn := bigs.String()
	if strn == "0" {
		return
	}

	if len(strn) > decimals {
		if decimals > Erc20Decimals {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+Erc20Decimals])
		} else {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+decimals])
		}
	} else {
		if decimals > Erc20Decimals {
			if len(strn) > decimals-Erc20Decimals {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
				strvalue = fmt.Sprintf(strFormat, strn[:len(strn)-(decimals-Erc20Decimals)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strvalue = fmt.Sprintf(strFormat, strn)
		}
	}

	//mylog.DataLogger.Info().Msgf("HexParseEthvalue: %v", strvalue)
	return
}

//eth wei -> eth count
func FromWeiWithDecimals(s string, decimals int) (strvalue string, err error) {

	strvalue = "0"

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

	strn := n.String()
	if strn == "0" {
		return
	}

	if len(strn) > decimals {
		if decimals > Erc20Decimals {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+Erc20Decimals])
		} else {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+decimals])
		}
	} else {
		if decimals > Erc20Decimals {
			if len(strn) > decimals-Erc20Decimals {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", Erc20Decimals, "s")
				strvalue = fmt.Sprintf(strFormat, strn[:len(strn)-(decimals-Erc20Decimals)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strvalue = fmt.Sprintf(strFormat, strn)
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

func ToWei(svalue string, decimals int) (*big.Int, error) {
	bigdecimals := big.NewInt(10)
	bigdecimals = bigdecimals.Exp(bigdecimals, big.NewInt(int64(decimals)), nil)

	priceparts := strings.Split(svalue, ".")
	if len(priceparts) == 1 {
		count, err := strconv.ParseUint(svalue, 10, 64)
		if err != nil {
			return nil, err
		}

		iamount := big.NewInt(int64(count))
		amount := iamount.Mul(iamount, bigdecimals)
		return amount, nil
	} else if len(priceparts) == 2 {
		uprice1, err := strconv.ParseUint(priceparts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		if len(priceparts[1]) <= 0 {
			iamount := big.NewInt(int64(uprice1))
			amount := iamount.Mul(iamount, bigdecimals)
			return amount, nil
		}

		if len(priceparts[1]) > decimals {
			priceparts[1] = priceparts[1][:decimals]
		}
		uprice2, err2 := strconv.ParseUint(priceparts[1], 10, 64)
		if err2 != nil {
			err = err2
			return nil, err
		}

		iamount1 := big.NewInt(int64(uprice1))
		iamount1.Mul(iamount1, bigdecimals)
		iamount2 := big.NewInt(int64(uprice2))
		iamount2.Mul(iamount2, bigdecimals.Exp(big.NewInt(10), big.NewInt(int64(decimals-len(priceparts[1]))), nil))
		amount := iamount1.Add(iamount1, iamount2)
		return amount, nil
	}

	return nil, errors.New("input invalid")
}
