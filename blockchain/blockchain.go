package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"goblockchain/constcoe"
	"goblockchain/transaction"
	"goblockchain/utils"
	"log"
	"runtime"

	"github.com/dgraph-io/badger/v4"
)

//区块链
type BlockChain struct {
	LastHash 		[]byte	//指当前区块链的最后一个区块	
	Database	*badger.DB
}

//用于遍历区块链
type BlockChainIterator struct {
	CurrentHash []byte
	Database	*badger.DB
}
//迭代器初始化
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	return &iterator
}

//迭代每次返回一个block，然后迭代器指向前一个区块的哈希值
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block
	
	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash
	return block
}

//判断迭代器是否终止,比较迭代器的currentHash和ogprehash就能确定迭代器是否迭代到头
func (chain *BlockChain) BackOgPreHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})
		utils.Handle(err)

		return err
	})	
	utils.Handle(err)

	return ogprevhash
}


//区块链创建区块的方法,使得区块链可以根据其它信息创建区块进行储存
func (bc *BlockChain) AddBlock(newBlock *Block) {
	// newBlock := CreateBlock(bc.Blocks[len(bc.Blocks) - 1].Hash, txs)
	// bc.Blocks = append(bc.Blocks, newBlock)
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	//先检查新区块的preHash是否等于lh
	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		fmt.Println("This block is out of age")
		runtime.Goexit()
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

// //区块链初始化函数，返回一个包含创世区块的区块链
// func CreateBlockChain() *BlockChain {
// 	blockChain := BlockChain{}
// 	blockChain.Blocks = append(blockChain.Blocks, GenesisBlock())
// 	return &blockChain
// }

//初始化一个区块链并且创建一个数据库保存
func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte		

	//检查是否有存储区块链的数据存在
	if utils.FileExists(constcoe.BCFile) {
		log.Fatal("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCPath)	//使用默认配置启动数据库
	// opts.Logger = nil	//让数据库的操作信息不输出到标准输出中

	db, err := badger.Open(opts)
	utils.Handle(err)
	
	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("genesis is created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)	//lh = lasthash，当前区块链最后一个区块的hash值
		utils.Handle(err)
		err = txn.Set([]byte("ogprevhash"), genesis.PrevHash)
		utils.Handle(err)
		lastHash = genesis.Hash
		return err
	})

	utils.Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain

}

//读取数据库加载区块链
func ContinueBlockChain() *BlockChain {
	if utils.FileExists(constcoe.BCFile) == false {
		log.Fatal("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) //读取最后一个区块值
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash=val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

//根据目标地址寻找可用的交易信息
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	//key是交易信息的id，value是output在该交易信息中的序号
	spentTxs := make(map[string][]int)	//记录遍历区块链时那些已经被使用的交易信息output
	
	iter := bc.Iterator()

	all:
		for{
			block := iter.Next()

			for _, tx := range block.Transaction {
				txID :=	hex.EncodeToString(tx.ID)

			IterOutputs:
				//遍历交易信息中的output，如果output在spentTxs中就跳过，说明该output已经被消费
				for outIdx, out := range tx.OutPuts {
					if spentTxs[txID] != nil {
						for _, spentOut := range spentTxs[txID] {
							if spentOut == outIdx {
								continue IterOutputs
							}
						}
					}

					if out.ToAddressRight(address) {
						unSpentTxs = append(unSpentTxs, *tx)
					}
				}
				if !tx.IsBase() {
					for _, in := range tx.Inputs {
						if in.FromAddressRight(address)	{
							inTxID := hex.EncodeToString(in.TxID)
							spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
						}
					}
				}
			}
			if bytes.Equal(block.PrevHash, bc.BackOgPreHash()) {
				break all
			}
		}
	return unSpentTxs
}

// func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
// 	unspentOuts := make(map[string]int)
// 	unspentTxs := bc.FindUnspentTransactions(address)
// 	accumulated := 0

// 	Work:
// 		for _, tx := range unspentTxs {
// 			txID := hex.EncodeToString(tx.ID)
// 			for outIdx, out := range tx.OutPuts {
// 				if out.ToAddressRight(address) {
// 					accumulated += out.Value
// 					unspentOuts[txID] = outIdx
// 					continue Work
// 				}
// 			}
// 		}
// 		return accumulated, unspentOuts
// }

//不用每次都找到全部的utxo，只用找到资产总值大于本次交易转账额的一部分utxo就行
func (bc *BlockChain) FindeSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTxs {
			txID := hex.EncodeToString(tx.ID)
			for outIdx, out := range tx.OutPuts {
				if out.ToAddressRight(address) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = outIdx
					if accumulated >= amount {
						break Work
					}
					continue Work
				}
			}
		}
		return accumulated, unspentOuts
}

func (bc *BlockChain) CreateTransaction(from_PubKey, to_HashPubKey []byte, amount int, privKey ecdsa.PrivateKey) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindeSpendableOutputs(from_PubKey, amount)

	if acc < amount {
		log.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}

	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{txID, outidx, from_PubKey, nil}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{amount, to_HashPubKey})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{acc - amount, utils.PublicKeyHash(from_PubKey)})
	}
	tx := transaction.Transaction{nil, inputs, outputs}

	tx.SetID()

	tx.Sign(privKey)

	return	&tx, true
}

func (chain *BlockChain) GetCurrentBlock() *Block {
	var block *Block
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(chain.LastHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	return block
}

func (bc *BlockChain) BackHeight() int64 {
	return bc.GetCurrentBlock().Height
}

//根据钱包地址返回该区块链中所有该地址的UTXO
func (bc *BlockChain) BackUTXOs(address []byte) []transaction.UTXO {
	var UTXOs []transaction.UTXO
	unspentTxs := bc.FindUnspentTransactions(address)

	Work:
		for _, tx := range unspentTxs {
			for outIdx, out := range tx.OutPuts {
				if out.ToAddressRight(address) {
					UTXOs = append(UTXOs, transaction.UTXO{tx.ID, outIdx, out})
					continue Work
				}
			}
		}
	return UTXOs
}

