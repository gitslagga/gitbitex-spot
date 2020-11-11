package service

import (
	"crypto/ecdsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
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
	return address, AddAddress(address)
}

func AddressAsset(userId int64) error {
	accounts, err := mysql.SharedStore().GetAccountsByUserId(userId)
	if len(accounts) > 0 || err != nil {
		return nil
	}

	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	//币币账户(USDT,BITC)
	err = db.AddAccount(&models.Account{UserId: userId, Currency: "USDT"})
	if err != nil {
		return err
	}
	err = db.AddAccount(&models.Account{UserId: userId, Currency: "BITC"})
	if err != nil {
		return err
	}

	//资产账户(YTL,BITC,ENERGY,USDT)
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: "YTL"})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: "BITC"})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: "ENERGY"})
	if err != nil {
		return err
	}
	err = db.AddAccountAsset(&models.AccountAsset{UserId: userId, Currency: "USDT"})
	if err != nil {
		return err
	}

	//矿池账户(BITC)
	err = db.AddAccountPool(&models.AccountPool{UserId: userId, Currency: "BITC"})
	if err != nil {
		return err
	}

	//购物账户(BITC,USDT)
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: "BITC"})
	if err != nil {
		return err
	}
	err = db.AddAccountShop(&models.AccountShop{UserId: userId, Currency: "USDT"})
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
		return nil, errors.New("cannot convert claim to MapClaims")
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	addressValue, found := claim["address"]
	if !found {
		return nil, errors.New("bad token")
	}
	addr := addressValue.(string)

	privateKeyVal, found := claim["private_key"]
	if !found {
		return nil, errors.New("bad token")
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

func AddAddress(address *models.Address) error {
	return mysql.SharedStore().AddAddress(address)
}

func UpdateAddress(address *models.Address) error {
	return mysql.SharedStore().UpdateAddress(address)
}

func ActivationAddress(address *models.Address, number float64, targetAddress string) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	//TODO 如果激活所需数量大于等于配置数量，并且该账户未激活，进行激活

	return db.CommitTx()
}
