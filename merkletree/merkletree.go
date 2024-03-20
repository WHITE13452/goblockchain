package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"goblockchain/transaction"
	"goblockchain/utils"

)


type MerkleTreeNode struct {
	leftNode 	*MerkleTreeNode
	reightNode 	*MerkleTreeNode
	Data		[]byte //hash
}

type MerkleTree struct {
	RootNode	*MerkleTreeNode
}

func CreateMerkleNode(left, right *MerkleTreeNode, data []byte) *MerkleTreeNode {
	tempNode := MerkleTreeNode{}

	//叶节点	
	if left == nil && right == nil {
		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.leftNode = left
	tempNode.reightNode = right

	return &tempNode
}

func CreateMerkleTree(txs []*transaction.Transaction) *MerkleTree {
	txslen := len(txs)
	//如果叶子结点为奇数个，那么复制最后一个节点
	if txslen%2 != 0 {
		txs = append(txs, txs[txslen-1])	
	}

	var nodePool []*MerkleTreeNode

	//交易都作为叶节点存储
	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleTreeNode //用来存上层节点
		poolLen := len(nodePool)
		//当发现某一层枝节点为奇数个时，我们将最后一个枝节点不做处理直接加入到上层枝节点最前面
		if poolLen % 2 != 0 {
			//刚开始temp为空，所以append是在最前面
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}

	return &merkleTree
}

//dfs用于返回叶子节点的mt验证路径
//route表示方向，0 = 左边，1 = 右边
//验证路径由两部分组成：1、指定用于拼凑的哈希值；2、该哈希值是拼在前边还是后边
func (mn *MerkleTreeNode) Find(data []byte, route []int, hashroute [][]byte) (bool, []int, [][]byte) {
	findFlag := false

	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashroute
	} else {
		if mn.leftNode != nil {
			route_t := append(route, 0)
			hashroute_t := append(hashroute, mn.reightNode.Data)
			findFlag, route_t, hashroute_t = mn.leftNode.Find(data, route_t, hashroute_t)
			if findFlag {
				return findFlag, route_t, hashroute_t
			} else {
				if mn.reightNode != nil {
					route_t = append(route, 1)
					hashroute_t = append(hashroute, mn.leftNode.Data)
					findFlag, route_t, hashroute_t = mn.reightNode.Find(data, route_t, hashroute_t)
					if findFlag {
						return findFlag, route_t, hashroute_t
					} else {
						return findFlag, route, hashroute
					}
				}
			}
		} else {
			return findFlag, route, hashroute
		}
	}
	return findFlag, route, hashroute

}

//把上边丑陋的dfs封装一下
//返回验证路径和是否找到该交易信息的信号
func (mn *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte, bool) {
	ok ,route, hashroute := mn.RootNode.Find(txid, []int{}, [][]byte{})
	return route, hashroute, ok
}

//SPV函数
//按照MT验证路径验证交易信息是否有效，如果成功则返回True，否则返回False
func SimplyPaymentValidation(txid, mtRootHash []byte, route []int, hashRoute [][]byte) bool {
	routeLen := len(route)
	tempHash := txid

	for i := routeLen - 1; i >= 0; i-- {
		if route[i] == 0 { //左边
			catenateHash := append(tempHash, hashRoute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if route[i] == 1 { //右边
			catenateHash := append(hashRoute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			utils.Handle(errors.New("error in validate route"))
		}
	}
	return bytes.Equal(tempHash, mtRootHash)
}


