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
)

const (
	// 1-赠送矿机，2-购买矿机
	MachineGiveAwayId = 1
	MachineFree       = 1
	MachineBuy        = 2
)

const (
	// 1-资产账户，2-矿池账户，3-币币账户，4-商城账户
	TransferAccountAsset = 1
	TransferAccountPool  = 2
	TransferAccountSpot  = 3
	TransferAccountShop  = 4
)

const (
	// 1-Ytl兑换BITE，2-Bite兑换Ytl
	ConvertYtlToBite = 1
	ConvertBiteToYtl = 2
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
