// 交易信息被分为两部分，分别为input与output
package transaction

import (
	"bytes"
	"encoding/gob"
	"goblockchain/utils"
)

//交易信息的output，去哪里
type TxOutput struct {
	Value		int		//转出价值
	// ToAddress	[]byte	//资产接受者地址
	HashPubKey	[]byte //用公钥hash做output地址
}
//交易信息的input，来自哪里
type TxInput struct {
	TxID		[]byte	//指明支持本次交易的前置交易信息
	OutIdx		int		//前置交易信息中的第几次output
	// FromAddress	[]byte	//资产转出者地址	
	PubKey		[]byte 	//用公钥做input地址
	Sig			[]byte	
}

//本地化utxo，需要包括该output的所有信息
//因为utxo存储的是余额，也就是区块链中没有被用的output
type UTXO struct {
	TxID	[]byte
	OutIdx	int
	TxOutput
}

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.PubKey, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.HashPubKey, utils.PublicKeyHash(address))
}

func (u *UTXO) SerializeUTXO() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(u)
	utils.Handle(err)
	return res.Bytes()
}

func DeSerializeUTXO(data []byte) *UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	utils.Handle(err)
	return &utxo
}

