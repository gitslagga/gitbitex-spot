package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"time"
)

var minDeposit = decimal.NewFromFloat(0.1)
var ethColdAddress1 = "0x82d2658D3fF713fbDA59f39aEA584975D7442407"
var ethColdAddress2 string
var blockHeight uint64

// ETH归集任务
func StartTransactionScan() {
	var err error
	blockHeight, err = models.SharedRedis().GetEthLatestHeight()
	if err != nil {
		panic(err)
	}

	if blockHeight == 0 {
		blockHeight, err = EthBlockNumber()
		if err != nil {
			panic(err)
		}
		err = models.SharedRedis().SetEthLatestHeight(blockHeight)
		if err != nil {
			panic(err)
		}
	}

	mylog.DataLogger.Info().Msgf("first scan height:%v", blockHeight)

	TransactionScan()

	t := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-t.C:
			TransactionScan()
		}
	}
}

func TransactionScan() {
	maxHeight, err := EthBlockNumber()
	if err != nil {
		mylog.DataLogger.Error().Msgf("[TransactionScan] EthBlockNumber, err:%v", err)
		return
	}

	maxHeight -= 16

	for curBlock := blockHeight + 1; curBlock <= maxHeight; curBlock++ {
		err := updateBlockDataStatus(curBlock - 15)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[TransactionScan] updateBlockDataStatus, err: %v", err)
			break
		}

		transactions, err := GetBlockByNumber(curBlock)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[TransactionScan] getBlockByNumber, err:%v", err)
			break
		}

		var index int
		for index = 0; index < len(transactions); index++ {

			var rowTransaction = transactions[index]

			value, err := FromWeiWithDecimals(rowTransaction.Value, EthDecimals)
			if err != nil {
				mylog.DataLogger.Error().Msgf("[TransactionScan] HexParseEthvalue, err:%v", err)
				continue
			}

			if value != "0" { //ETH
				err = ETHDataHandle(rowTransaction, curBlock, value)
				if err != nil {
					mylog.DataLogger.Error().Msgf("[TransactionScan] ETHDataHandle err: %v", err)
					continue
				}
			} else { //ERC20
				boolean := InputDataIsTransfer(rowTransaction.Input)
				if boolean == false {
					continue
				}

				err = ERC20DataHandle(rowTransaction, curBlock)
				if err != nil {
					mylog.DataLogger.Error().Msgf("[TransactionScan] ERC20DataHandle err: %v", err)
					continue
				}
			}
		}

		if index == len(transactions) {
			blockHeight = curBlock
			err = models.SharedRedis().SetEthLatestHeight(blockHeight)
			if err != nil {
				mylog.DataLogger.Error().Msgf("[TransactionScan] redis SetEthLatestHeight err:%v", err)
			}
		} else {
			break
		}
	}
}

func updateBlockDataStatus(curBlock uint64) error {
	addressDeposits, err := mysql.SharedStore().GetAddressDepositsByBNStatus(curBlock, models.CurrencyDepositUnConfirm)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[updateBlockDataStatus] GetAddressDepositsByBNStatus err:%v, blockNum:%v", err, curBlock)
		return err
	}

	for _, deposit := range addressDeposits {
		// modify account
		err = modifyAccount(deposit)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[updateBlockDataStatus] modifyAccount, err:%v", err)
			return err
		}
	}

	return nil
}

func modifyAccount(deposit *models.AddressDeposit) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(deposit.UserId, deposit.Coin)
	if err != nil {
		return err
	}

	coinAsset.Available = coinAsset.Available.Add(deposit.Value)
	err = db.UpdateAccountAsset(coinAsset)
	if err != nil {
		return err
	}

	deposit.Status = models.CurrencyDepositConfirmed
	err = db.UpdateAddressDeposit(deposit)
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func ETHDataHandle(transaction *Transaction, blockNum uint64, value string) error {
	//1 查询address是否存在于db
	address, err := mysql.SharedStore().GetAddressByAddress(transaction.ToAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ETHDataHandle] GetAddressByAddress, err:%v, txId:%v", err, transaction.Txid)
		return err
	}
	if address == nil {
		return nil
	}

	mylog.DataLogger.Info().Msgf("[ETHDataHandle] begin ETH, userID:%v, txId:%v, address:%v, value:%v", address.Id, transaction.Txid, transaction.ToAddress, value)

	//2 eth_getTransactionReceipt
	transactionReceipt, err := GetTransactionReceipt(transaction.Txid)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ETHDataHandle] eth_getTransactionReceipt, err: %v, txId:%v, value:%v", err, transaction.Txid, value)
		return err
	}
	if transactionReceipt.Result.Status == "0x0" {
		mylog.DataLogger.Info().Msgf("[ETHDataHandle] eth_getTransactionReceipt, status failed, txId:%v", transaction.Txid)
		return nil
	}

	//3 是否小于最小充值金额
	valueF, err := decimal.NewFromString(value)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ETHDataHandle] decimal NewFromString err:%v, txId:%v", err, transaction.Txid)
		return err
	}
	if valueF.LessThan(minDeposit) {
		mylog.DataLogger.Info().Msgf("[ETHDataHandle] eth value less than minDeposit, value:%v, minDeposit:%v, txId:%v", value, minDeposit, transaction.Txid)
		return nil
	}

	//汇集
	ethSendMainItem := ethAddressItem{
		FromAddress: transaction.ToAddress,
		Count:       0,
	}
	ethToMainChan <- &ethSendMainItem

	mylog.DataLogger.Info().Msgf("[ETHDataHandle] deposit ETH, userID:%v, txId:%v", address.Id, transaction.Txid)

	return nil
}

func ERC20DataHandle(transaction *Transaction, blockNum uint64) error {
	//1 查询合约地址是否存在与db
	config, err := mysql.SharedStore().GetAddressConfigByContract(transaction.ToAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20DataHandle] getContractNameFromMysql, err:%v, address:%v, txId:%v", err, transaction.ToAddress, transaction.Txid)
		return err
	}
	if config == nil {
		return nil
	}

	//2 eth_getTransactionReceipt
	transactionReceipt, err := GetTransactionReceipt(transaction.Txid)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20DataHandle] eth_getTransactionReceipt, err: %v, txId:%v", err, transaction.Txid)
		return err
	}
	if transactionReceipt.Result.Status == "0x0" || len(transactionReceipt.Result.Logs) <= 0 {
		//mylog.DataLogger.Info().Msgf("[ERC20DataHandle] eth_getTransactionReceipt, status failed, txId:%v", transaction.Txid)
		return nil
	}

	for _, transactionLog := range transactionReceipt.Result.Logs {
		err = ERC20transactionLogHandle(&transactionLog, transaction.Txid, blockNum, config)
		if err != nil {
			mylog.DataLogger.Error().Msgf("[ERC20DataHandle] erc20transactionLogHandle, err:%v, txId:%v", err, transaction.Txid)
			continue
		}
	}

	return nil
}

func ERC20transactionLogHandle(transactionLog *TransactionReceiptLog, txId string, blockNum uint64, config *models.AddressConfig) error {
	//3 parse transaction log
	isTransfer, toAddress, value, err := ParseTransactionLog(transactionLog, config.Decimals)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20transactionLogHandle] ParseTransactionLog, err:%v, txId:%v", err, txId)
		return err
	} else if isTransfer == false {
		mylog.DataLogger.Error().Msgf("[ERC20transactionLogHandle] is not transfer, txId:%v", txId)
		return nil
	}

	//4 查询address是否存在
	address, err := mysql.SharedStore().GetAddressByAddress(toAddress)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20transactionLogHandle] SelectAddressIsExistFromMysql, err:%v, txId:%v", err, txId)
		return err
	}
	if address == nil {
		return nil
	}

	mylog.DataLogger.Info().Msgf("[ERC20transactionLogHandle] begin %v, userID:%v, txId:%v, address:%v, value:%v", config.Coin, address.Id, txId, toAddress, value)

	valueF, err := decimal.NewFromString(value)
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20transactionLogHandle] ParseFloat, err:%v, value:%v, txId:%v", err, value, txId)
		return err
	}
	if valueF.LessThan(config.MinDeposit) {
		mylog.DataLogger.Info().Msgf("[ERC20transactionLogHandle] ERC20 %v, value less than value:%v, minDeposit:%v, txId:%v", config.Coin, value, config.MinDeposit, txId)
		return nil
	}

	//汇集
	toMainItem := tokenAddressItem{
		FromAddress:   toAddress,
		Token:         config.Coin,
		GasPrice:      EthGasPrice,
		GasLimit:      Erc20GasLimit,
		SendTimestamp: time.Now().Unix(),
		Count:         0,
	}
	tokenToMainChan <- &toMainItem

	//5, db
	err = mysql.SharedStore().AddAddressDeposit(&models.AddressDeposit{
		UserId:   address.Id,
		BlockNum: blockNum,
		TxId:     txId,
		Coin:     config.Coin,
		Address:  toAddress,
		Value:    valueF,
		Actual:   valueF,
		Status:   models.CurrencyDepositUnConfirm,
	})
	if err != nil {
		mylog.DataLogger.Error().Msgf("[ERC20transactionLogHandle] SaveEthTransactionToDB, err:%v, blockNum:%v, txId:%s, value:%v", err, blockNum, txId, value)
		return err
	}

	mylog.DataLogger.Info().Msgf("[ERC20transactionLogHandle] deposit %v, toAddress:%v, txId:%v, value:%v", config.Coin, toAddress, txId, value)

	return nil
}
