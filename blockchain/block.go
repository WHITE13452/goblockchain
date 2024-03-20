package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/merkletree"
	"goblockchain/transaction"
	"goblockchain/utils"
	"time"
)

//对于Nonce： 每个节点去寻找一个随机值（也就是nonce），将这个随机值作为候选区块的头部信息属性之一，
//			要求候选区块对自身信息（注意这里是包含了nonce的）进行哈希后表示为数值要小于一个难度目标值（也就是Target），最先寻找到nonce的节点即为卷王，可以将自己的候选区块发布并添加到区块链尾部。
//区块
type Block struct {
	Timestamp 	int64 //时间戳
	Hash 		[]byte //该区块的哈希值
	PrevHash 	[]byte //上个区块的哈希
	Height		int64	//区块在区块链中的高度，便于系统中各设备节点更新维护区块链，辅助检查本地utxo集是否过时
	Target		[]byte //目标难度值
	Nonce		int64 //随机值
	//Data 		[]byte //数据
	Transaction	[]*transaction.Transaction//交易信息
	MerkleTree	*merkletree.MerkleTree
	
}

//设置该区块的哈希值（时间戳+preHash+数据）
func (b *Block) SetHash() {
	//将三个数据连接起来，Join的第二个参数是连接符，这里设置为byte的空
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce),b.BackTransactionSummary(), b.MerkleTree.RootNode.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

//创建区块
func CreateBlock(prevhash []byte, height int64, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, height, []byte{}, 0, txs, merkletree.CreateMerkleTree(txs)}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}
//创世区块
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte("white is cool!"), 0, []*transaction.Transaction{tx})
	genesis.SetHash()
	// genesisWords := "Hello! BlockChain!"
	//创世区块的前一块哈希值为空即可
	return genesis
}

//返回当前区块所有交易的id（也就是交易的hash）总和
func (b *Block) BackTransactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transaction {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

//用于序列化区块，Badger键值对只能存储字节串
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encode := gob.NewEncoder(&res)
	err := encode.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}