package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"go/types"
	"log"
	"time"
)

func main {

}

//区块
type Block struct {
	Timestamp int64 //时间戳
	Hash []byte //该区块的哈希值
	PrevHash []byte //上个区块的哈希
	Data []byte //数据
}

//区块链
type BlockChain struct {
	Blocks []*Block
}

//设置该区块的哈希值（时间戳+preHash+数据）
func (b *Block) SetHash() {
	//将三个数据连接起来，Join的第二个参数是连接符，这里设置为byte的空
	information := bytes.Join([][]byte{ToHexInt(b.Timestamp), b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}
//将int64转换为字节穿
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}