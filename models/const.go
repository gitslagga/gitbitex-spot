package models

const (
	TopicOrder   = "g_order"
	TopicAccount = "g_account"
	TopicFill    = "g_fill"
	TopicBill    = "g_bill"
)

const (
	AccountMaxAddress = 100

	AccountCurrencyYtl  = "YTL"
	AccountCurrencyBite = "BITE"
	AccountCurrencyUsdt = "USDT"

	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	AccountAssetTransfer = 1
	AccountPoolTransfer  = 2
	AccountSpotTransfer  = 3
	AccountShopTransfer  = 4

	// 扫一扫支付币种, 1-未支付，2-已支付
	AccountScanCurrency    = "BITE"
	AccountScanUnPayment   = 1
	AccountScanPayment     = 2
	AccountScanReleaseRate = 0.05

	// 兑换收益和矿池收益的5%划分给商城账户
	AccountConvertShopRate = 0.05
	AccountHoldingShopRate = 0.05
	AccountPromoteShopRate = 0.05

	// 推广收益算力计算
	AccountPromoteStandard = 10000
	AccountPromoteMaxCal   = 90000
	AccountPromoteMinCal   = 10

	// 拼团
	AccountGroupCurrencyUsdt = "USDT"
	AccountGroupCurrencyBite = "BITE"
	AccountGroupPersonNumber = 7
	AccountGroupPlayUsdt     = 20
	AccountGroupPlayBite     = 20
	AccountGroupProcess      = 1
	AccountGroupFinish       = 2
	AccountGroupLogProcess   = 1
	AccountGroupLogSuccess   = 2
	AccountGroupLogFailed    = 3
	AccountGroupRefundUsdt   = 21
	AccountGroupRefundBite   = 21
	AccountGroupIncreaseDay  = 30
	AccountGroupDirectRate   = 0.02
	AccountGroupIndirectRate = 0.02
	AccountGroupNodeRate     = 0.02
	AccountGroupStakeUsdt    = 500
	AccountGroupStakeBite    = 500
	AccountGroupReleaseMode  = 0
	AccountGroupDelegateMode = 1
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

	// 最低持币量，最佳持币量
	ConfigMinHolding
	ConfigBestHolding
)

const (
	// 充币，提币
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

const (
	// 认购管理
	IssueFirstEveryMonth  = 0
	IssueSecondEveryMonth = 1
	IssueThreeEveryMonth  = 2
	IssueFourEveryMonth   = 3

	// 兑换，扫码，拼团 资金池
	AddressReleaseConvert = 1
	AddressReleaseScan    = 2
	AddressReleaseGroup   = 3
)
