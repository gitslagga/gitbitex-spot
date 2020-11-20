package models

const (
	TopicOrder   = "g_order"
	TopicAccount = "g_account"
	TopicFill    = "g_fill"
	TopicBill    = "g_bill"
)

const (
	AccountCurrencyYtl    = "YTL"
	AccountCurrencyBite   = "BITE"
	AccountCurrencyUsdt   = "USDT"
	AccountCurrencyEnergy = "ENERGY" // Deprecated

	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	AccountAssetTransfer = 1
	AccountPoolTransfer  = 2
	AccountSpotTransfer  = 3
	AccountShopTransfer  = 4
)

const (
	// 1-赠送矿机，2-购买矿机
	MachineGiveAwayId = 1
	MachineFree       = 1
	MachineBuy        = 2

	// 1-Ytl兑换BITE，2-Bite兑换Ytl
	MachineYtlConvertBite = 1
	MachineBiteConvertYtl = 2

	// 有效推荐数（活跃的大于1）
	MachineValidInvite = 1
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
	ConfigActiveTransfer = iota
	ConfigIssueReward
	ConfigYtlConvert
	ConfigHoldCoinProfit
	ConfigPromoteProfit

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

	ConfigYtlConvertUsdt
	ConfigEnergyConvertUsdt // Deprecated
	ConfigBiteConvertUsdt
	ConfigUsdtConvertCny

	UsdtMinDepositNumber
	UsdtMinWithdrawNumber
	UsdtWithdrawFee
	UsdtCollectionAddress
	UsdtCollectFeeAddress
)
