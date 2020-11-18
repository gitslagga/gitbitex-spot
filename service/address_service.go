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

	address.Username = username
	address.Password = password

	addressExists, err := GetAddressByAddress(address.Address)
	if err != nil {
		return nil, err
	}

	if addressExists != nil {
		address.Id = addressExists.Id
		return address, UpdateAddress(address)
	}

	address.ConvertFee = decimal.NewFromFloat(0.5)
	return address, AddAddress(address)
}

func AddressLogin(mnemonic, privateKey, password string) (address *models.Address, err error) {
	if mnemonic != "" {
		address, err = createAddressByMnemonic(mnemonic)
	} else {
		address, err = createAddressByPrivateKey(privateKey)
	}
	if err != nil {
		return nil, err
	}

	address.Password = password

	addressExists, err := GetAddressByAddress(address.Address)
	if err != nil {
		return nil, err
	}
	if addressExists != nil {
		address.Id = addressExists.Id
		address.Username = addressExists.Username
		return address, UpdateAddress(address)
	}

	address.Username = "Account1"
	address.ConvertFee = decimal.NewFromFloat(0.5)
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
	err = db.AddAccount(&models.Account{UserId: userId, Currency: models.CURRENCY_BITE})
	if err != nil {
		return err
	}
	err = db.AddAccount(&models.Account{UserId: userId, Currency: models.CURRENCY_USDT})
	if err != nil {
		return err
	}

	//资产账户(YTL,BITE,USDT)
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.CURRENCY_YTL})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.CURRENCY_BITE})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: models.CURRENCY_USDT})
	if err != nil {
		return err
	}

	//矿池账户(BITE)
	err = db.AddAccountPool(&models.AccountPool{UserId: userId, Currency: models.CURRENCY_BITE})
	if err != nil {
		return err
	}

	//购物账户(BITE,USDT)
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: models.CURRENCY_BITE})
	if err != nil {
		return err
	}
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: models.CURRENCY_USDT})
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

	configNum, err := strconv.ParseFloat(configs[0].Value, 64)
	if err != nil {
		return err
	}
	if number < configNum {
		return errors.New("激活数量不足|Insufficient number of activations")
	}

	//目标账户已经激活
	targetAddress, err := mysql.SharedStore().GetAddressByAddress(addressValue)
	if err != nil {
		return err
	}
	if targetAddress == nil {
		return errors.New("地址不存在|Address is not exists")
	}
	if targetAddress.ParentIds != "" {
		return errors.New("地址已经激活|Address is always activation")
	}

	return activationAddress(address, decimal.NewFromFloat(number), targetAddress, configs)
}

func activationAddress(address *models.Address, number decimal.Decimal,
	targetAddress *models.Address, configs []*models.Config) error {
	//进行激活
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CURRENCY_BITE)
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

	targetAddressAsset, err := db.GetAccountAssetForUpdate(targetAddress.Id, models.CURRENCY_BITE)
	if err != nil {
		return err
	}

	targetAddressAsset.Available = targetAddressAsset.Available.Add(number)
	err = db.UpdateAccountAsset(targetAddressAsset)
	if err != nil {
		return err
	}

	address.InviteNum++
	var inviteNum int
	var convertFee decimal.Decimal
	for i := 5; i < 10; i++ {
		inviteNum, err = strconv.Atoi(configs[i].Value)
		if err != nil {
			return err
		}
		convertFee, err = decimal.NewFromString(configs[i+5].Value)
		if err != nil {
			return err
		}

		if address.InviteNum >= inviteNum {
			address.ConvertFee = convertFee
		}
	}

	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	if address.ParentIds == "" {
		targetAddress.ParentIds = fmt.Sprintf("%d", address.Id)
	} else {
		targetAddress.ParentIds = fmt.Sprintf("%s-%d", address.ParentIds, address.Id)
	}
	err = db.UpdateAddress(targetAddress)
	if err != nil {
		return err
	}

	//激活下级，赠送上级一级矿机
	machine, err := db.GetMachineById(1)
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
		IsBuy:       0,
	})

	return db.CommitTx()
}
