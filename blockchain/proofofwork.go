package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"math"
	"math/big"
)

//返回目标难度值
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	//向左移位，移的越多目标难度值越大，哈希取值落在的空间就更多就越容易找到符合条件的nonce
	target.Lsh(target, uint(256 - constcoe.Difficulty))
	return target.Bytes()
}

//根据不同的nonce会有不同的区块哈希值
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(int64(nonce)),
		b.Target,
		b.BackTransactionSummary(),
	}, []byte{})
	return data
}

//对于任意一个区块找到合适的nonce
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64 =0
	// nonce = 0
	intTarget.SetBytes(b.Target)
	
	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}
//检验是否合适
func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)

	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(&intTarget) == -1
}