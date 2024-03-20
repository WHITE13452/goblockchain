//一个钱包对应一个用户，但是一个用户可以有多个钱包
//所有需要walletmanager来管理所有钱包，同一个用户的所有钱包都保存在一个机器上
//walletmanager中refList就是一个钱包的list，用map来存储

package wallet

import (
	"bytes"
	"crypto/elliptic"
	// "crypto/rand"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"os"
)

//钱包主要是为了保存密钥对
type Wallet struct {
	PrivateKey []byte
	PublicKey []byte
}

// //椭圆曲线生成密钥对
// func NewKeyPair() (ecdsa.PrivateKey, []byte) {
// 	curve := elliptic.P256()

// 	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	utils.Handle(err)
// 	//横纵坐标拼接保存
// 	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
// 	return *privateKey, publicKey
// }

func NewWallet() *Wallet {
	strprivateKey, strpublicKey, err := utils.GenKeyPair()
	utils.Handle(err)
	wallet := Wallet{[]byte(strprivateKey), []byte(strpublicKey)}
	return &wallet
}

func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

func (w *Wallet) Save()  {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)
	err = os.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}
	var w Wallet
	gob.Register(elliptic.P256())
	fileContent, err := os.ReadFile(filename)
	utils.Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}
