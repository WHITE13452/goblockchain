package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"os"
	"path/filepath"
	"strings"

)

//实际应用中不需要别名去管理，这里是为了方便演示
//用来记录机器上记录的钱包,key为钱包地址，value为钱包别名
type RefList map[string]string

//保存RefList
func (r *RefList) Save()  {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)

	err = os.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

//用于扫描机器上保存的所有钱包文件（.wlt文件）（特别是检查是否存在从其他机器上拷贝的钱包）
func (r *RefList) Update()  {
	//扫描wallets文件夹下所有文件
	err := filepath.Walk(constcoe.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName) - 4:], ".wlt") ==0 {
			_, ok := (*r) [fileName[:len(fileName) - 4]]
			if !ok {
				(*r)[fileName[:len(fileName) - 4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var refList RefList
	if utils.FileExists(filename) {
		fileContent, err := os.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&refList)
		utils.Handle(err)
	} else {
		refList = make(RefList)
		refList.Update()
	}
	return &refList
}

//绑定别名
func (r *RefList) BindRefName(adddress, refname string)  {
	(*r)[adddress] = refname
}

//通过别名调取钱包地址
func (r *RefList) FindRef(refname string) (string, error) {
	tmp := ""
	for key, value := range *r {
		if value == refname {
			tmp = key
			break
		}
	}
	if tmp == "" {
		err := errors.New("the refname is not found")
		return tmp, err
	}
	return tmp, nil
}

