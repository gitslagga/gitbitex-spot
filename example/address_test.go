package example

import (
	"crypto/ecdsa"
	"crypto/md5"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"log"
	"testing"
)

func TestCreateMnemonic(t *testing.T) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(160)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
}

func TestGetAddressByMnemonic(t *testing.T) {
	mnemonic := "paddle either arctic cereal puzzle shiver fiscal become turn bridge tunnel brisk"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := wallet.PublicKeyHex(account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("privateKey:", privateKey)
	fmt.Println("publicKey", publicKey)
	fmt.Println("address:", account.Address.Hex())
}

func TestGetAddressByPrivateKey(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("826fa3446c54d5ced4845e164172f6bef8cb419ccc4ca686811dcd0df429d5bb")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyValue := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyValue := hexutil.Encode(publicKeyBytes)[2:]

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	fmt.Println("privateKey:", privateKeyValue)
	fmt.Println("publicKey", publicKeyValue)
	fmt.Println("address:", address)
}

func TestGetBiteAddress(t *testing.T) {
	hash := md5.Sum([]byte(`0x82d2658D3fF713fbDA59f39aEA584975D7442407`))
	biteAddress := fmt.Sprintf("BITE%x", hash)

	fmt.Println("hash:", hash)
	fmt.Println("biteAddress:", biteAddress)
}
