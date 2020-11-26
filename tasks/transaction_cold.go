package tasks

import (
	"container/list"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type ethAddressItem struct {
	FromAddress string
	Count       int
}

type tokenAddressItem struct {
	FromAddress   string
	Token         string
	GasPrice      int64
	GasLimit      int64
	SendTimestamp int64
	Count         int
}

var ethToMainChan chan *ethAddressItem = make(chan *ethAddressItem, 1000)
var tokenToMainChan chan *tokenAddressItem = make(chan *tokenAddressItem, 10000)

var tokenToMainListMutex = sync.Mutex{}
var tokenToMainList = list.New()

func mutexPushBackTokenToMainList(addressItem *tokenAddressItem) {
	tokenToMainListMutex.Lock()
	tokenToMainList.PushBack(addressItem)
	tokenToMainListMutex.Unlock()
}

func mutexPushFrontTokenToMainList(addressItem *tokenAddressItem) {
	tokenToMainListMutex.Lock()
	tokenToMainList.PushFront(addressItem)
	tokenToMainListMutex.Unlock()
}

func mutexRemoveTokenFromMainList(item *list.Element) {
	tokenToMainListMutex.Lock()
	tokenToMainList.Remove(item)
	tokenToMainListMutex.Unlock()
}

func StartSendToMainTask() {
	go startSendEthToMainTask()
	go startSendTokenToMainTask1()
	go startSendTokenToMainTask2()
}

func startSendEthToMainTask() {
	for {
		select {
		case addressItem, ok := <-ethToMainChan:
			if !ok {
				time.Sleep(time.Second * 1)
				continue
			}

			err := SendRowEthToMainWallet(addressItem.FromAddress)
			if err != nil {
				mylog.DataLogger.Error().Msgf("[sendAllToMainWallet] SendRowEthToMainWallet err:%v, address:%v", err, addressItem.FromAddress)

				ethToMainChan <- addressItem

				addressItem.Count += 1

				time.Sleep(time.Minute * 3)
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func startSendTokenToMainTask1() {
	for {
		select {
		case addressItem, ok := <-tokenToMainChan:
			if !ok {
				time.Sleep(time.Second * 1)
				continue
			}

			t := time.Now().Unix()
			if t > addressItem.SendTimestamp {
				mutexPushFrontTokenToMainList(addressItem)
			} else {
				mutexPushBackTokenToMainList(addressItem)
			}
		}
	}
}

func startSendTokenToMainTask2() {
	timer := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-timer.C:

			tmpItem := tokenToMainList.Front()
			if tmpItem == nil {
				continue
			}
			addressItem := tmpItem.Value.(*tokenAddressItem)

			if addressItem.SendTimestamp > time.Now().Unix() {
				continue
			} else {
				mutexRemoveTokenFromMainList(tmpItem)

				isSuccess, err := SendRowTokenToMainWallet(addressItem)
				if isSuccess == false {
					if err != nil {
						mylog.DataLogger.Error().Msgf("[COLD] SendRowTokenToMainWallet err:%v", err)
					}

					addressItem.Count += 1

					if addressItem.Count <= 21 {
						//手续费不够时，打入ETH手续费。打入的手续费大于最低充币数量时，会被充币轮询。
						//所以配置的最小充币数量，应该大于代币归集手续费。
						addressItem.SendTimestamp = time.Now().Unix() + 300*int64(addressItem.Count*3) //300s后重试
						tokenToMainChan <- addressItem
					} else {
						mylog.DataLogger.Error().Msgf("[COLD] SendRowTokenToMainWallet failed count:%d, address:%v", addressItem.Count, addressItem.FromAddress)
					}
				}
			}

		}
	}
}

func SendRowEthToMainWallet(address string) error {
	bigAmount, err := GetBalance(address, EthName)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] GetBalance rows scan error! err:%v", err)
		return err
	}

	strAmount, err := FromWeiWithBigintAndToken(bigAmount, EthName)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] FromWeiWithToken error! err:%v", err)
		return err
	}
	fAmount, err := decimal.NewFromString(strAmount)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] decimal NewFromString err:%v", err)
		return err
	}
	if fAmount.LessThan(minDeposit) {
		return nil
	}

	mylog.DataLogger.Info().Msgf("[COLD] sendToMainWallet eth, address: %v", address)

	account, err := mysql.SharedStore().GetAddressByAddress(address)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SelectPrivkeyFromMysql error! err:%v", err)
		return err
	}

	//发送到main address
	txid, err := EthPersonalSendTransactionToBlockWithFee(address, account.PrivateKey, ethColdAddress1, bigAmount)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SendTransactionToMainWallet error! err:%v", err)
		return err
	}

	mylog.DataLogger.Info().Msgf("[COLD] sendToMainWallet success eth, txid: %v, token: %v, value: %v", txid, EthName, strAmount)

	return nil
}

//发送eth矿工费到拥有token的地址上
func sendEthFeeToTokenAddress(address, token, ethFee string) error {

	mylog.DataLogger.Info().Msgf("[COLD] SendEthFee, token:%v, address:%v", token, address)

	//发送eth手续费给此地址
	txid, err := PersonalSendTransactionToBlock(EthMainAddress, address, EthName, ethFee)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] PersonalSendTransactionToBlock err:%v", err)
		return err
	}

	mylog.DataLogger.Info().Msgf("[COLD] SendEthFeeToTokenaddress, address:%v, fee:%v, txid:%v", address, ethFee, txid)

	return nil
}

func isNeedSendToMainWallet(address, token string) (bool, error) {
	bigAmount, err := GetBalance(address, token)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] isNeedSendToMainWallet GetBalance rows scan error! err:%v", err)
		return false, err
	} else if bigAmount.Uint64() == 0 {
		return false, nil
	}

	config, err := mysql.SharedStore().GetAddressConfigByCoin(token)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] isNeedSendToMainWallet SelectTokenInfoFromDB err:%v", err)
		return false, err
	}
	if config == nil {
		return false, nil
	}

	depositMin, _ := config.MinDeposit.Float64()

	balanceStr, err := FromWeiWithBigintAndToken(bigAmount, token)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] isNeedSendToMainWallet FromWei err:%v", err)
		return false, err
	}

	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] isNeedSendToMainWallet ParseFloat %s err:%v", balanceStr, err)
		return false, err
	}

	if balance < depositMin {
		return false, nil
	}

	return true, nil
}

func SendRowTokenToMainWallet(addressItem *tokenAddressItem) (bool, error) {
	if addressItem.GasPrice == 0 {
		addressItem.GasPrice = EthGasPrice
	}
	if addressItem.GasLimit == 0 {
		addressItem.GasLimit = Erc20GasLimit
	}

	bigEthFee := big.NewInt(addressItem.GasPrice * addressItem.GasLimit)
	ethFee, err := EthFromWei(bigEthFee)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SendRowTokenToMainWallet FromWei err:%v", err)
		return false, err
	}

	bigEthBalance, err := GetBalance(addressItem.FromAddress, EthName)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] GetBalance err:%v", err)
		return false, err
	}
	strEthBalance, err := FromWeiWithBigintAndToken(bigEthBalance, EthName)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] FromWeiWithToken err:%v", err)
		return false, err
	}
	fEthBalance, _ := strconv.ParseFloat(strEthBalance, 64)
	fEthFee, _ := strconv.ParseFloat(ethFee, 64)
	if fEthBalance < fEthFee {
		needSend, err := isNeedSendToMainWallet(addressItem.FromAddress, addressItem.Token)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[COLD] isNeedSendToMainWallet err:%v", err)
			return false, err
		} else if needSend == false {
			return true, nil
		}

		err = sendEthFeeToTokenAddress(addressItem.FromAddress, addressItem.Token, ethFee)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[COLD] SendEthFeeToTokenAddress err:%v", err)
			return false, err
		}

		return false, nil
	}

	bigAmount, err := GetBalance(addressItem.FromAddress, addressItem.Token)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] GetBalance rows scan error! err:%v", err)
		return false, err
	} else if bigAmount.Int64() == 0 {
		return true, nil
	}

	mylog.DataLogger.Info().Msgf("[COLD] sendToMainWallet, token: %v, address: %v", addressItem.Token, addressItem.FromAddress)

	account, err := mysql.SharedStore().GetAddressByAddress(addressItem.FromAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SelectPrivkeyFromMysql error! err:%v", err)
		return false, err
	}

	//发送到main address
	txId, err := Erc20PersonalSendTransactionToBlock(addressItem.FromAddress, account.PrivateKey, ethColdAddress2, addressItem.Token, bigAmount)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SendTransactionToMainWallet error! err:%v", err)
		return false, err
	}

	//db
	amountStr, _ := FromWeiWithBigintAndToken(bigAmount, addressItem.Token)
	value, _ := decimal.NewFromString(amountStr)

	err = mysql.SharedStore().AddAddressCollect(&models.AddressCollect{
		UserId:      account.Id,
		Coin:        addressItem.Token,
		TxId:        txId,
		FromAddress: addressItem.FromAddress,
		ToAddress:   ethColdAddress2,
		Value:       value,
		Status:      models.CurrencyCollectionCold,
	})
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] SaveEthExtractTransactionToMysql error, txId:%s, value:%v, err:%v", txId, value, err)
	}

	mylog.DataLogger.Info().Msgf("[COLD] sendToMainWallet success, txId: %v, token: %v, value: %v", txId, addressItem.Token, value)

	return true, nil
}

func GetBalance(address string, coin string) (*big.Int, error) {

	var hexAmount string
	var err error

	if coin == EthName { //查看eth balance
		hexAmount, err = EthGetBalance(address, "latest")
		if err != nil {
			mylog.DataLogger.Error().Msgf("[COLD] GetBalance, block err: %v", err)
			return nil, err
		}

		//mylog.DataLogger.Info().Msgf("amount: %v", hexAmount)
	} else { //查看 coin balance
		config, err := mysql.SharedStore().GetAddressConfigByCoin(coin)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[COLD] GetBalance, GetContractAddress err: %v", err)
			return nil, err
		}

		hexAmount, err = TokenGetBalance(address, config.ContractAddress)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[COLD] GetBalance, TokenGetBalance  block err: %v", err)
			return nil, err
		}
	}

	bigAmount, err := GetBigFromHex(hexAmount)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[COLD] GetBalance, GetBigFromHex err:%v", err)
		return nil, err
	}

	//mylog.DataLogger.Info().Msgf("[COLD] GetBalance address:%s, coin:%s, balance:%d", address, coin, bigAmount.Uint64())

	return bigAmount, nil
}
