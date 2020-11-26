package models

const (
	TopicOrder   = "g_order"
	TopicAccount = "g_account"
	TopicFill    = "g_fill"
	TopicBill    = "g_bill"
)

const (
	AccountCurrencyYtl  = "YTL"
	AccountCurrencyBite = "BITE"
	AccountCurrencyUsdt = "USDT"

	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	AccountAssetTransfer = 1
	AccountPoolTransfer  = 2
	AccountSpotTransfer  = 3
	AccountShopTransfer  = 4

	// 扫一扫支付币种, 1-未支付，2-已支付
	AccountScanCurrency  = "BITE"
	AccountScanUnPayment = 1
	AccountScanPayment   = 2
)

const (
	// 1-赠送矿机，2-购买矿机
	MachineGiveAwayId = 1
	MachineFree       = 1
	MachineBuy        = 2

	// 1-Ytl兑换BITE，2-Bite兑换Ytl
	MachineYtlConvertBite = 1
	MachineBiteConvertYtl = 2

	// 达人等级有效推荐数（活跃的大于1）
	MachineValidInvite = 1
	// YTL兑换BITE有效推荐数（购买矿机）
	MachineEffectInvite = 2

	// 0-普通用户，1-一级达人，2-二级达人，3-三级达人，4-四级达人，5-五级达人
	MachineLevelZero  = 0
	MachineLevelOne   = 1
	MachineLevelTwo   = 2
	MachineLevelThree = 3
	MachineLevelFour  = 4
	MachineLevelFive  = 5
)

const (
	ConfigActiveTransfer = iota
	ConfigIssueReward
	ConfigYtlConvert
	ConfigHoldCoinProfit
	ConfigPromoteProfit

	ConfigYtlConvertUsdt
	ConfigBiteConvertUsdt
	ConfigUsdtConvertCny
	ConfigYtlConvertBiteFee

	// 扫一扫支付条件
	ConfigScanStartHour
	ConfigScanEndHour
	ConfigScanMinPayment
	ConfigScanMaxPayment
	ConfigScanDayPayment
	ConfigScanFeePayment
)

const (
	CurrencyCollectionMaster = 1
	CurrencyCollectionCold   = 2

	CurrencyDepositUnConfirm = 1
	CurrencyDepositConfirmed = 2

	CurrencyWithdrawReview   = 1
	CurrencyWithdrawSuccess  = 2
	CurrencyWithdrawPassed   = 3
	CurrencyWithdrawUnPass   = 4
	CurrencyWithdrawCanceled = 5
)
