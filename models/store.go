package models

import (
	"github.com/shopspring/decimal"
)

type Store interface {
	BeginTx() (Store, error)
	Rollback() error
	CommitTx() error

	GetConfigs() ([]*Config, error)
	GetConfigById(id int64) (*Config, error)
	UpdateConfig(config *Config) error

	GetUserByEmail(email string) (*User, error)
	AddUser(user *User) error
	UpdateUser(user *User) error

	GetAccount(userId int64, currency string) (*Account, error)
	GetAccountsByUserId(userId int64) ([]*Account, error)
	GetAccountForUpdate(userId int64, currency string) (*Account, error)
	AddAccount(account *Account) error
	UpdateAccount(account *Account) error

	GetUnsettledBillsByUserId(userId int64, currency string) ([]*Bill, error)
	GetUnsettledBills() ([]*Bill, error)
	AddBills(bills []*Bill) error
	UpdateBill(bill *Bill) error

	GetProductById(id string) (*Product, error)
	GetProducts() ([]*Product, error)

	GetOrderById(orderId int64) (*Order, error)
	GetOrderByClientOid(userId int64, clientOid string) (*Order, error)
	GetOrderByIdForUpdate(orderId int64) (*Order, error)
	GetOrdersByUserId(userId int64, statuses []OrderStatus, side *Side, productId string,
		beforeId, afterId int64, limit int) ([]*Order, error)
	AddOrder(order *Order) error
	UpdateOrder(order *Order) error
	UpdateOrderStatus(orderId int64, oldStatus, newStatus OrderStatus) (bool, error)

	GetLastFillByProductId(productId string) (*Fill, error)
	GetUnsettledFillsByOrderId(orderId int64) ([]*Fill, error)
	GetUnsettledFills(count int32) ([]*Fill, error)
	UpdateFill(fill *Fill) error
	AddFills(fills []*Fill) error

	GetLastTradeByProductId(productId string) (*Trade, error)
	GetTradesByProductId(productId string, count int) ([]*Trade, error)
	AddTrades(trades []*Trade) error

	GetTicksByProductId(productId string, granularity int64, limit int) ([]*Tick, error)
	GetLastTickByProductId(productId string, granularity int64) (*Tick, error)
	AddTicks(ticks []*Tick) error

	//development new
	GetAddressByAddress(addr string) (*Address, error)
	GetAddressById(id int64) (*Address, error)
	CountAddressByMachineLevelId(machineLevelId int64) (int, error)
	GetAddressByMachineLevelId(machineLevelId int64) ([]*Address, error)
	CountAddressByGroupUsdt() (int, error)
	GetAddressByGroupUsdt() ([]*Address, error)
	CountAddressByGroupBite() (int, error)
	GetAddressByGroupBite() ([]*Address, error)
	GetAddressByParentId(parentId int64) ([]*Address, error)
	AddAddress(address *Address) error
	UpdateAddress(address *Address) error
	GetAddressHoldingByUserId(userId, beforeId, afterId, limit int64) ([]*AddressHolding, error)
	GetLastAddressHolding() (*AddressHolding, error)
	AddAddressHolding(holding *AddressHolding) error
	GetTotalPowerList() ([]*TotalPower, error)
	GetAddressPromoteByUserId(userId, beforeId, afterId, limit int64) ([]*AddressPromote, error)
	GetLastAddressPromote() (*AddressPromote, error)
	AddAddressPromote(promote *AddressPromote) error

	GetAddressListByAddress(address string) (*AddressList, error)
	GetAddressListById(id int64) (*AddressList, error)
	GetAddressListByUserId(userId int64) ([]*AddressList, error)
	CountAddressListByUserId(userId int64) (int, error)
	AddAddressList(addressList *AddressList) error
	UpdateAddressList(addressList *AddressList) error
	DeleteAddressList(addressListId int64) error

	GetAccountAsset(userId int64, currency string) (*AccountAsset, error)
	GetAccountsAssetByUserId(userId int64) ([]*AccountAsset, error)
	GetAccountAssetForUpdate(userId int64, currency string) (*AccountAsset, error)
	GetIssueAccountAsset() ([]*AccountAsset, error)
	SumIssueAccountAsset() (decimal.Decimal, error)
	GetHoldingAccountAsset(minHolding decimal.Decimal) ([]*AccountAsset, error)
	AddAccountAsset(account *AccountAsset) error
	UpdateAccountAsset(account *AccountAsset) error

	GetAccountPool(userId int64, currency string) (*AccountPool, error)
	GetAccountsPoolByUserId(userId int64) ([]*AccountPool, error)
	GetAccountPoolForUpdate(userId int64, currency string) (*AccountPool, error)
	AddAccountPool(account *AccountPool) error
	UpdateAccountPool(account *AccountPool) error

	GetAccountShop(userId int64, currency string) (*AccountShop, error)
	GetAccountsShopByUserId(userId int64) ([]*AccountShop, error)
	GetAccountShopForUpdate(userId int64, currency string) (*AccountShop, error)
	AddAccountShop(account *AccountShop) error
	UpdateAccountShop(account *AccountShop) error

	GetBuyMachine() ([]*Machine, error)
	GetMachineById(machineId int64) (*Machine, error)
	GetMachineAddressByUserId(userId, before, after, limit int64) ([]*MachineAddress, error)
	CountMachineAddressUsed(userId int64, machineId int64) (int, error)
	GetMachineAddressUsedList() ([]*MachineAddress, error)
	AddMachineAddress(machineAddress *MachineAddress) error
	UpdateMachineAddress(machineAddress *MachineAddress) error
	GetMachineLogByUserId(userId, before, after, limit int64) ([]*MachineLog, error)
	GetLastMachineLog(machineAddressId int64) (*MachineLog, error)
	AddMachineLog(machineLog *MachineLog) error

	GetMachineConvertByUserId(userId, before, after, limit int64) ([]*MachineConvert, error)
	GetMachineConvertSumNum() (decimal.Decimal, error)
	GetMachineConvertSumFee() (decimal.Decimal, error)
	AddMachineConvert(machineConvert *MachineConvert) error

	GetMachineLevel() ([]*MachineLevel, error)
	GetMachineLevelById(machineLevelId int64) (*MachineLevel, error)
	GetMachineConfigs() ([]*MachineConfig, error)
	GetMachineConfigById(id int64) (*MachineConfig, error)
	UpdateMachineConfig(config *MachineConfig) error

	GetAccountTransferByUserId(userId, before, after, limit int64) ([]*AccountTransfer, error)
	AddAccountTransfer(accountTransfer *AccountTransfer) error
	GetAccountScanByUserId(userId, before, after, limit int64) ([]*AccountScan, error)
	GetAccountScanSumNumber(userId int64) (decimal.Decimal, error)
	GetAccountScanSumFee() (decimal.Decimal, error)
	AddAccountScan(accountScan *AccountScan) error

	GetValidAddressConfig() ([]*AddressConfig, error)
	GetAddressConfigByCoin(coin string) (*AddressConfig, error)
	GetAddressConfigByContract(contract string) (*AddressConfig, error)
	UpdateAddressConfig(config *AddressConfig) error
	AddAddressCollect(collect *AddressCollect) error
	GetAddressDepositsByUserId(userId, before, after, limit int64) ([]*AddressDeposit, error)
	GetAddressDepositsByBNStatus(blockNum uint64, status int) ([]*AddressDeposit, error)
	AddAddressDeposit(deposit *AddressDeposit) error
	UpdateAddressDeposit(deposit *AddressDeposit) error
	GetAddressWithdrawsByUserId(userId, before, after, limit int64) ([]*AddressWithdraw, error)
	GetAddressWithdrawsByOrderSN(orderSN string) (*AddressWithdraw, error)
	AddAddressWithdraw(withdraw *AddressWithdraw) error
	UpdateAddressWithdraw(withdraw *AddressWithdraw) error

	GetIssueByUserId(userId, beforeId, afterId, limit int64) ([]*Issue, error)
	GetIssueUsedList() ([]*Issue, error)
	AddIssue(issue *Issue) error
	UpdateIssue(issue *Issue) error
	GetIssueConfigs() ([]*IssueConfig, error)
	GetIssueLogByUserId(userId, beforeId, afterId, limit int64) ([]*IssueLog, error)
	GetLastIssueLog(issueId int64) (*IssueLog, error)
	AddIssueLog(issueLog *IssueLog) error

	GetGroupById(groupId int64) (*Group, error)
	GetGroupByUserIdCoin(userId int64, coin string) (*Group, error)
	GetGroupByCoin(coin string, beforeId, afterId, limit int64) ([]*Group, error)
	AddGroup(group *Group) error
	UpdateGroup(group *Group) error
	GetGroupLogByUserId(userId, beforeId, afterId, limit int64) ([]*GroupLog, error)
	GetGroupLogPublicity(beforeId, afterId, limit int64) ([]*GroupLog, error)
	GetGroupLogByGroupId(groupId int64) ([]*GroupLog, error)
	GetGroupLogByGroupIdUserId(groupId, userId int64) (*GroupLog, error)
	GetGroupLogSumNum(coin string) (decimal.Decimal, error)
	AddGroupLog(groupLog *GroupLog) error
	UpdateGroupLog(groupLog *GroupLog) error

	GetAddressReleaseByUserId(userId, beforeId, afterId, limit int64) ([]*AddressRelease, error)
	GetLastAddressRelease(releaseType int) (*AddressRelease, error)
	AddAddressRelease(release *AddressRelease) error
}
