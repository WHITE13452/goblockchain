package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/constcoe"
	"log"
	"os"

	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/ripemd160"
)

//错误处理
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//将int64转换为字节串
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return	false
	}
	return true
}

//公钥hash
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

//检查生成函数
func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constcoe.ChecksumLength]
}

//base256转base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input))
	Handle(err)
	return decode	
}

//公钥哈希生成钱包地址	
func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionHash := append([]byte{constcoe.NetworkVersion}, pubKeyHash...)
	checkSum := Checksum(networkVersionHash)
	finalHash := append(networkVersionHash, checkSum...)
	address := Base58Encode(finalHash)
	return address	
}

func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-constcoe.ChecksumLength]
	return pubKeyHash
}

func ByteToInt64(num []byte) int64 {
	var num64 int64
	buff := bytes.NewBuffer(num)
	err := binary.Read(buff, binary.BigEndian, &num64)
	Handle(err)
	return num64
}
