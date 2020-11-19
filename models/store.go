package models

import "github.com/shopspring/decimal"

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
	AddAddress(address *Address) error
	UpdateAddress(address *Address) error

	GetAccountTransferByUserId(userId int64) ([]*AccountTransfer, error)
	AddAccountTransfer(accountTransfer *AccountTransfer) error

	GetAccountAsset(userId int64, currency string) (*AccountAsset, error)
	GetAccountsAssetByUserId(userId int64) ([]*AccountAsset, error)
	GetAccountAssetForUpdate(userId int64, currency string) (*AccountAsset, error)
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
	GetMachineAddressByUserId(userId int64) ([]*MachineAddress, error)
	GetMachineAddressUsedCount(userId int64, machineId int64) (int, error)
	GetMachineAddressUsedList() ([]*MachineAddress, error)
	AddMachineAddress(machineAddress *MachineAddress) error
	UpdateMachineAddress(machineAddress *MachineAddress) error
	GetMachineLogByUserId(userId int64) ([]*MachineLog, error)
	GetLastMachineLog(machineAddressId int64) (*MachineLog, error)
	AddMachineLog(machineLog *MachineLog) error

	GetMachineConvertByUserId(userId int64) ([]*MachineConvert, error)
	GetMachineConvertSumNumber() (decimal.Decimal, error)
	GetMachineConvertSumFee() (decimal.Decimal, error)
	AddMachineConvert(machineConvert *MachineConvert) error

	GetMachineLevel() ([]*MachineLevel, error)
	GetMachineLevelById(machineLevelId int64) (*MachineLevel, error)
}
