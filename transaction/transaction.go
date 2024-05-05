package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/constcoe"
	"goblockchain/utils"
)

type Transaction struct {
	ID		[]byte		//其自身的id值（其实就是hash）
	Inputs	[]TxInput	//用于记录本次转账前置交易信息的txoutput
	//记录txouput可以实现找零，只需要将本次找零记录其中，设置流入方向就是本次交易的sender，也就是给sender找零
	OutPuts []TxOutput	//记录本次转账的的接受者和总数
}

func (tx *Transaction)	TxHash()	[]byte {
	var encoded bytes.Buffer
	var hash [32]byte
	//gob用来序列化结构体
	encoder := gob.NewEncoder(&encoded)	
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID()  {
	tx.ID = tx.TxHash()
}

//创始的交易信息
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	tx := Transaction{[]byte("This is the Fucking Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut} }
	return  &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) ==1 && tx.Inputs[0].OutIdx ==-1
}

//交易过程  https://zhuanlan.zhihu.com/p/435733874
//只是单纯的复制一下交易的信息，拿到这笔交易的input和output
func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txIn := range tx.Inputs {
		inputs = append(inputs, TxInput{txIn.TxID, txIn.OutIdx, nil, nil})
	}

	for _, txOut := range tx.OutPuts {
		outputs = append(outputs, TxOutput{txOut.Value, txOut.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

//用来辅助交易函数签名
func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.TxHash()
}

//签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey)  {
	if tx.IsBase() {
		return
	}

	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		signature, err := utils.Sign(string(plainhash), privKey)
		utils.Handle(err)
		tx.Inputs[idx].Sig = []byte(signature)
	}
}

func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainHash := tx.PlainHash(idx, input.PubKey)
		if !utils.VerifySign(string(plainHash), string(input.Sig), string(input.PubKey)){
			return false	
		}
	}
	return true
}