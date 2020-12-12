package service

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/utils"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tyler-smith/go-bip39"
	"strconv"
	"time"
)

func CreateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func createAddressByMnemonic(mnemonic string) (*models.Address, error) {
	seed, err := hdwallet.NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return nil, err
	}

	publicKey, err := wallet.PublicKeyHex(account)
	if err != nil {
		return nil, err
	}
	return &models.Address{
		Address:    account.Address.Hex(),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Mnemonic:   mnemonic,
	}, nil
}

func createAddressByPrivateKey(privateKey string) (*models.Address, error) {
	privateK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicK := privateK.Public()
	publicKeyECDSA, ok := publicK.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateK)
	privateKey = hexutil.Encode(privateKeyBytes)[2:]

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKey := hexutil.Encode(publicKeyBytes)[2:]

	return &models.Address{
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA).Hex(),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func AddressRegister(username, password, mnemonic string) (*models.Address, error) {
	address, err := createAddressByMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	//子地址不能注册或登录，会引发地址混乱冲突
	addressList, err := mysql.SharedStore().GetAddressListByAddress(address.Address)
	if err != nil {
		return nil, err
	}
	if addressList != nil {
		return nil, errors.New("子地址不能注册或登录|Sub address cannot be registered or logged in")
	}

	addressExists, err := GetAddressByAddress(address.Address)
	if err != nil {
		return nil, err
	}

	if addressExists != nil {
		addressExists.Username = username
		addressExists.Password = password
		return addressExists, UpdateAddress(addressExists)
	}

	address.Username = username
	address.Password = password
	address.AddressBite = utils.GetBiteAddress(address.Address)

	config, err := GetConfigById(models.ConfigYtlConvertBiteFee + 1)
	if err != nil {
		return nil, err
	}
	address.ConvertFee, err = decimal.NewFromString(config.Value)
	if err != nil {
		return nil, err
	}

	return address, AddAddress(address)
}

func AddressLogin(username, password, mnemonic, privateKey string) (address *models.Address, err error) {
	if mnemonic != "" {
		address, err = createAddressByMnemonic(mnemonic)
		if err != nil {
			return nil, err
		}
	} else {
		address, err = createAddressByPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		address, err := GetAddressByAddress(address.Address)
		if err != nil {
			return nil, err
		}
		if address == nil {
			return nil, errors.New("外部私钥不能导入|External private key cannot be imported")
		}
	}

	//子地址不能注册或登录，会引发地址混乱冲突
	addressList, err := mysql.SharedStore().GetAddressListByAddress(address.Address)
	if err != nil {
		return nil, err
	}
	if addressList != nil {
		return nil, errors.New("子地址不能注册或登录|Sub address cannot be registered or logged in")
	}

	addressExists, err := GetAddressByAddress(address.Address)
	if err != nil {
		return nil, err
	}
	if addressExists != nil {
		addressExists.Username = username
		addressExists.Password = password
		return addressExists, UpdateAddress(addressExists)
	}

	address.Username = username
	address.Password = password
	address.AddressBite = utils.GetBiteAddress(address.Address)

	config, err := GetConfigById(models.ConfigYtlConvertBiteFee + 1)
	if err != nil {
		return nil, err
	}
	address.ConvertFee, err = decimal.NewFromString(config.Value)
	if err != nil {
		return nil, err
	}

	return address, AddAddress(address)
}

func AddressAsset(userId int64) error {
	accounts, err := mysql.SharedStore().GetAccountsAssetByUserId(userId)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		return nil
	}

	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	//币币账户(USDT,BITE)
	err = db.AddAccount(&models.Account{UserId: userId, Currency: models.AccountCurrencyBite})
	if err != nil {
		return err
	}
	err = db.AddAccount(&models.Account{UserId: userId, Currency: models.AccountCurrencyUsdt})
	if err != nil {
		return err
	}

	//资产账户(YTL,BITE,USDT)
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.AccountCurrencyYtl})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.AccountCurrencyBite})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.AccountCurrencyUsdt})
	if err != nil {
		return err
	}

	//矿池账户(BITE)
	err = db.AddAccountPool(&models.AccountPool{UserId: userId, Currency: models.AccountCurrencyBite})
	if err != nil {
		return err
	}

	//购物账户(BITE,USDT)
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: models.AccountCurrencyBite})
	if err != nil {
		return err
	}
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: models.AccountCurrencyUsdt})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func CreateFrontendToken(address *models.Address) (string, error) {
	claim := jwt.MapClaims{
		"id":          address.Id,
		"address":     address.Address,
		"private_key": address.PrivateKey,
		"exp":         time.Now().Add(time.Second * time.Duration(60*60*24*7)).Unix(),
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(conf.GetConfig().JwtSecret))
}

func CheckFrontendToken(tokenStr string) (*models.Address, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetConfig().JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无法将声明转换为MapClaims|cannot convert claim to MapClaims")
	}
	if !token.Valid {
		return nil, errors.New("无效的token|token is invalid")
	}

	addressValue, found := claim["address"]
	if !found {
		return nil, errors.New("损坏的令牌|bad token")
	}
	addr := addressValue.(string)

	privateKeyVal, found := claim["private_key"]
	if !found {
		return nil, errors.New("损坏的令牌|bad token")
	}
	privateKey := privateKeyVal.(string)

	address, err := GetAddressByAddress(addr)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("bad token")
	}
	if address.PrivateKey != privateKey {
		return nil, errors.New("bad token")
	}
	return address, nil
}

func GetAddressByAddress(address string) (*models.Address, error) {
	return mysql.SharedStore().GetAddressByAddress(address)
}

func GetAddressById(id int64) (*models.Address, error) {
	return mysql.SharedStore().GetAddressById(id)
}

func GetAddressByUsername(username string) (bool, error) {
	address, err := mysql.SharedStore().GetAddressByUsername(username)
	if err != nil {
		return false, err
	}

	addressList, err := mysql.SharedStore().GetAddressListByUsername(username)
	if err != nil {
		return false, err
	}

	return address == nil && addressList == nil, nil
}

func GetAddressByParentId(parentId int64) ([]*models.Address, error) {
	return mysql.SharedStore().GetAddressByParentId(parentId)
}

func CountAddressByMachineLevelId(machineLevelId int64) (int, error) {
	return mysql.SharedStore().CountAddressByMachineLevelId(machineLevelId)
}

func AddAddress(address *models.Address) error {
	return mysql.SharedStore().AddAddress(address)
}

func UpdateAddress(address *models.Address) error {
	return mysql.SharedStore().UpdateAddress(address)
}

func ActivationAddress(address *models.Address, number float64, addressValue string) error {
	//激活数量需要大于等于配置数量
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	configNum, err := strconv.ParseFloat(configs[models.ConfigActiveTransfer].Value, 64)
	if err != nil {
		return err
	}
	if number < configNum {
		return errors.New("激活数量不足|Insufficient number of activations")
	}

	//目标账户已经激活
	targetAddress, err := mysql.SharedStore().GetAddressByUAddressBite(addressValue)
	if err != nil {
		return err
	}
	if targetAddress == nil {
		return errors.New("地址不存在|Address is not exists")
	}
	if targetAddress.ParentId != 0 {
		return errors.New("地址已经激活|Address is always activation")
	}

	return activationAddress(address, decimal.NewFromFloat(number), targetAddress)
}

func activationAddress(address *models.Address, number decimal.Decimal, targetAddress *models.Address) error {
	//进行激活
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(address.Id, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	if addressAsset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	addressAsset.Available = addressAsset.Available.Sub(number)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	targetAddressAsset, err := db.GetAccountAssetForUpdate(targetAddress.Id, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	targetAddressAsset.Available = targetAddressAsset.Available.Add(number)
	err = db.UpdateAccountAsset(targetAddressAsset)
	if err != nil {
		return err
	}

	//确认上下级关系
	targetAddress.ParentId = address.Id
	if address.ParentIds == "" {
		targetAddress.ParentIds = fmt.Sprintf("%d", address.Id)
	} else {
		targetAddress.ParentIds = fmt.Sprintf("%s,%d", address.ParentIds, address.Id)
	}
	err = db.UpdateAddress(targetAddress)
	if err != nil {
		return err
	}

	//赠送上级一级矿机
	machine, err := db.GetMachineById(models.MachineGiveAwayId)
	if err != nil {
		return err
	}
	err = db.AddMachineAddress(&models.MachineAddress{
		MachineId:   machine.Id,
		UserId:      address.Id,
		Number:      machine.Number.Add(machine.Number.Mul(machine.Profit)).Div(decimal.NewFromInt(int64(machine.Release))),
		TotalNumber: machine.Number.Add(machine.Number.Mul(machine.Profit)),
		Day:         machine.Release,
		TotalDay:    machine.Release,
		IsBuy:       models.MachineFree,
	})

	return db.CommitTx()
}
