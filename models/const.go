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
	GiveAwayMachineId = 1
	FreeMachine       = 0
	BuyMachine        = 1
)

const (
	ActiveTransferNumber = iota
	IssueRewardNumber
	YtlConvertNumber
	HoldCoinProfitNumber
	PromoteProfitNumber

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

	YtlConvertUsdtRate
	EnergyConvertUsdtRate // Deprecated
	BiteConvertUsdtRate
	UsdtConvertCnyRate

	UsdtMinDepositNumber
	UsdtMinWithdrawNumber
	UsdtWithdrawFee
	UsdtCollectionAddress
	UsdtCollectFeeAddress
)
