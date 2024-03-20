//UTXO本地化后的一些增删改查

package utxoset

import (
	"bytes"
	"fmt"
	"goblockchain/transaction"
	"goblockchain/utils"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v4"
)


var (
	info  		= "INFO"
	infoname 	= info + "NAME"
	infoheight 	= info + "HEIGHT"
	utxokey  	= "UTXO"
	utxokeyorder = ":ORDER:"
)

/**
一个UTXO集应该包括一个数据库用于维护其中的UTXO，
一个Name用于辨别这是什么的UTXO集，
以及一个Height用于说明该UTXO集维护的UTXO的有效性，便于判断UTXO集是否需要更新
**/
type UTXOSet struct {
	Name	[]byte
	DB		*badger.DB
	Height	int64
}

func GetUtxoSetFile(dir string) string {
	fileAddress := dir + "/" + "MANIFEST"
	return fileAddress
}

func ToUtxoKey(txID []byte, order int) []byte {
	utxokey := bytes.Join([][]byte{[]byte(utxokey), txID, []byte(utxokeyorder), utils.ToHexInt(int64(order))}, []byte{})
	return utxokey
}

func CreateUTXOSet(Name []byte, dir string, utxos []transaction.UTXO, height int64) *UTXOSet {
	if utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("UTXOSet has already existed, now rebuild it.")
		err := os.RemoveAll(dir)
		utils.Handle(err)
	}
	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	utxoSet := UTXOSet{Name, db, height}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoname), Name)
		if err != nil {
			return err
		}
		
		err = txn.Set([]byte(infoheight), utils.ToHexInt(height))
		if err != nil {
			return err
		}

		for _, utxo := range utxos {
			utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
			err = txn.Set(utxoKey, utxo.SerializeUTXO())
			return err
		}
		return nil
	})
	utils.Handle(err)
	return &utxoSet
}

func LoadUTXOSet(dir string) *UTXOSet {
	if !utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("No UTXOSet found, please create one first")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	var name []byte
	var height int64

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(infoname))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			name = val
			return nil
		})
		if err != nil {
			return err
		}
		item, err = txn.Get([]byte(infoheight))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			height = utils.ByteToInt64(val)
			return nil
		})
		return err
	})
	utils.Handle(err)

	utxoSet := UTXOSet{name, db, height}
	return &utxoSet
}

func (us *UTXOSet) AddUtxo(utxo *transaction.UTXO)  {
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
		err := txn.Set(utxoKey, utxo.SerializeUTXO())
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

func (us *UTXOSet) AddUTXO(txID []byte, outIdx int, output transaction.TxOutput)  {
	utxo := transaction.UTXO{txID, outIdx, output}
	us.AddUtxo(&utxo)
}

func (us *UTXOSet) DelUTXO(txID []byte, order int)  {
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(txID, order)
		err := txn.Delete(utxoKey)
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

func (us *UTXOSet) UpdateHeight(height int64)  {
	us.Height = height
	err := us.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoheight), utils.ToHexInt(height))
		return err
	})
	utils.Handle(err)
}

//检查数据库中的某一个键值对信息是否是height或name一类的描述信息还是一般的utxo信息
func IsInfo(infoKey []byte) bool {
	if bytes.HasPrefix(infoKey, []byte(info)) {
		return true
	} else {
		return false
	}
}