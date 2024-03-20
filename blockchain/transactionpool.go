package blockchain

import (
	"bytes"
	"encoding/gob"
	"goblockchain/constcoe"
	"goblockchain/transaction"
	"goblockchain/utils"
	"os"
)

//交易信息池,储存节点收集到的交易信息
type TransactionPool struct {
	PubTx []*transaction.Transaction
}

func (tp *TransactionPool) AddTransaction(tx *transaction.Transaction)  {
	tp.PubTx = append(tp.PubTx, tx)
}

func (tp *TransactionPool) SaveFile()  {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(tp)
	utils.Handle(err)
	//0644是8进制的644 （110，100，100）表示不同用户对文件读写执行的权限
	err = os.WriteFile(constcoe.TransactionPoolFile, content.Bytes(), 0644)
	utils.Handle(err)
}

func (tp *TransactionPool) LoadFile() error {
	if !utils.FileExists(constcoe.TransactionPoolFile) {
		return nil		
	}

	var transactionPool TransactionPool

	fileContent, err := os.ReadFile(constcoe.TransactionPoolFile)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&transactionPool)

	if err != nil {
		return err	
	}

	tp.PubTx = transactionPool.PubTx
	return nil
}

func CreateTransactionPool() *TransactionPool {
	// transactionPool := TransactionPool{}
	var transactionPool TransactionPool
	err := transactionPool.LoadFile()
	utils.Handle(err)
	return &transactionPool
}

func RemoveTransactionPoolFile() error {
	err := os.Remove(constcoe.TransactionPoolFile)	
	return err
}