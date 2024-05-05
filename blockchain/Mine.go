package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"goblockchain/transaction"
	"goblockchain/utils"
	"log"
)

//一个完整的Mine过程在构造候选区块时应该先检查交易信息池中的所有交易信息的有效性，
//这包括验证是否引用了已花费的Output，是否重复引用了同一UTXO，
//Input与Output资产总额是否对应，交易信息的签名信息。
func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}

	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		PubKey := tx.Inputs[0].PubKey
		unspentOutputs := bc.FindUnspentTransactions(PubKey)
		inputAmount := 0
		outputAmount := 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.OutPuts {
			outputAmount += output.Value
		}
		if outputAmount != inputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}

	}
	return true
}

func isInputRight(txs []transaction.Transaction, in transaction.TxInput) (bool, int) {
	 for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.OutPuts[in.OutIdx].Value
		}
	 }
	 return false, 0
}


//假设现在的单个节点每次挖矿都能胜出并将自己的候选区块加入到区块链中。
func (bc *BlockChain) RunMine()  {
	transactionPool := CreateTransactionPool()

	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

	//完成pow
	candidateBlock := CreateBlock(bc.LastHash, bc.BackHeight() + 1, transactionPool.PubTx)
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}