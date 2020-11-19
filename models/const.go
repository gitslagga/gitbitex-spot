package models

const (
	TopicOrder   = "g_order"
	TopicAccount = "g_account"
	TopicFill    = "g_fill"
	TopicBill    = "g_bill"
)

const (
	CurrencyYtl    = "YTL"
	CurrencyBite   = "BITE"
	CurrencyUsdt   = "USDT"
	CurrencyEnergy = "ENERGY" // Deprecated

	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	TransferAccountAsset = 1
	TransferAccountPool  = 2
	TransferAccountSpot  = 3
	TransferAccountShop  = 4
)

const (
	// 1-赠送矿机，2-购买矿机
	MachineGiveAwayId = 1
	MachineFree       = 1
	MachineBuy        = 2

	// 1-Ytl兑换BITE，2-Bite兑换Ytl
	MachineYtlConvertBite = 1
	MachineBiteConvertYtl = 2
)

const (
	// 0-普通用户，1-一级达人，2-二级达人，3-三级达人，4-四级达人，5-五级达人
	MachineLevelZero = iota
	MachineLevelOne
	MachineLevelTwo
	MachineLevelThree
	MachineLevelFour
	MachineLevelFive
)

const (
	TotalActiveTransfer = iota
	TotalIssueReward
	TotalYtlConvert
	TotalHoldCoinProfit
	TotalPromoteProfit

	YtlConvertInviteOne
	YtlConvertInviteTwo
	YtlConvertInviteThree
	YtlConvertInviteFour
	YtlConvertInviteFive
	YtlConvertFeeOne
	YtlConvertFeeTwo
	YtlConvertFeeThree
	YtlConvertFeeFour
	YtlConvertFeeFive

	RateYtlConvertUsdt
	RateEnergyConvertUsdt // Deprecated
	RateBiteConvertUsdt
	RateUsdtConvertCny

	UsdtMinDepositNumber
	UsdtMinWithdrawNumber
	UsdtWithdrawFee
	UsdtCollectionAddress
	UsdtCollectFeeAddress
)
