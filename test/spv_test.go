//测试spv和merkletree

package test

import (
	"crypto/sha256"
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/merkletree"
	"goblockchain/transaction"
	"strconv"
	"strings"
	"testing"
)

func GenerateTransaction(outCash int, inAccount string, toAccount string, prevTxID string, outIdx int) *transaction.Transaction {
	prevTxIDHash := sha256.Sum256([]byte(prevTxID))
    inAccountHash := sha256.Sum256([]byte(inAccount))
    toAccountHash := sha256.Sum256([]byte(toAccount))
    txIn := transaction.TxInput{prevTxIDHash[:], outIdx, inAccountHash[:], nil}
    txOut := transaction.TxOutput{outCash, toAccountHash[:]}
    tx := transaction.Transaction{[]byte("This is the Base Transaction!"),
        []transaction.TxInput{txIn}, []transaction.TxOutput{txOut}} //Whether set ID is not nessary
    tx.SetID() //Here the ID is reset to the hash of the whole transaction. Signature is skipped
    return &tx
}
var transactionTests = []struct {
    outCash   int
    inAccount string
    toAccount string
    prevTxID  string
    outIdx    int
}{
    {
        outCash:   10,
        inAccount: "LLL",
        toAccount: "CCC",
        prevTxID:  "prev1",
        outIdx:    1,
    },
    {
        outCash:   20,
        inAccount: "EEE",
        toAccount: "OOO",
        prevTxID:  "prev2",
        outIdx:    1,
    },
    {
        outCash:   30,
        inAccount: "OOO",
        toAccount: "EEE",
        prevTxID:  "prev3",
        outIdx:    0,
    },
    {
        outCash:   100,
        inAccount: "CCC",
        toAccount: "LLL",
        prevTxID:  "prev4",
        outIdx:    1,
    },
    {
        outCash:   50,
        inAccount: "AAA",
        toAccount: "OOO",
        prevTxID:  "prev5",
        outIdx:    1,
    },
    {
        outCash:   110,
        inAccount: "OOO",
        toAccount: "AAA",
        prevTxID:  "prev6",
        outIdx:    0,
    },
    {
        outCash:   200,
        inAccount: "LLL",
        toAccount: "CCC",
        prevTxID:  "prev7",
        outIdx:    1,
    },
    {
        outCash:   500,
        inAccount: "EEE",
        toAccount: "OOO",
        prevTxID:  "prev8",
        outIdx:    1,
    },
}

func GenerateBlock(txs []*transaction.Transaction, prevBlock string) *blockchain.Block {
	prevBlockHash := sha256.Sum256([]byte(prevBlock))
	testBlock := blockchain.CreateBlock(prevBlockHash[:], -1, txs)
	return testBlock
}

var spvTests = []struct {
    txContained []int
    prevBlock   string
    findTX      []int
    wants       []bool
}{
    {
        txContained: []int{0},
        prevBlock:   "prev1",
        findTX:      []int{0, 1},
        wants:       []bool{true, false},
    },
    {
        txContained: []int{0, 1, 2, 3, 4, 5, 6, 7},
        prevBlock:   "prev2",
        findTX:      []int{3, 7, 5},
        wants:       []bool{true, true, true},
    },
    {
        txContained: []int{0, 1, 2, 3},
        prevBlock:   "prev3",
        findTX:      []int{0, 1, 5},
        wants:       []bool{true, true, false},
    },
    {
        txContained: []int{0, 3, 5, 6, 7},
        prevBlock:   "prev4",
        findTX:      []int{0, 1, 6, 7},
        wants:       []bool{true, false, true, true},
    },
    {
        txContained: []int{0, 1, 2, 4, 5, 6, 7},
        prevBlock:   "prev5",
        findTX:      []int{0, 1, 3},
        wants:       []bool{true, true, false},
    },
}

func TestSPV(t *testing.T)  {
	primeTxs := []*transaction.Transaction{}
	for _, tx := range transactionTests {
		tx := GenerateTransaction(tx.outCash, tx.inAccount, tx.toAccount, tx.prevTxID, tx.outIdx)
		primeTxs = append(primeTxs, tx)
	}

	fmt.Println("TestSPV Begin...")
    fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for idx, test := range spvTests {
		fmt.Println("Current test No: ", idx)
        fmt.Println("Merkle Tree is like:")
		mtGraphPaint(test.txContained)
		txs := []*transaction.Transaction{}
		for _, txidx := range test.txContained {
			txs = append(txs, primeTxs[txidx])
		}
		
		testBlock := GenerateBlock(txs, test.prevBlock)
        fmt.Println("------------------------------------------------------------------")
        for num, findidx := range test.findTX {
            fmt.Println("Find transaction:", findidx)
            fmt.Printf("Transaction ID: %x\n", primeTxs[findidx].ID)
            route, hashroute, ok := testBlock.MerkleTree.BackValidationRoute(primeTxs[findidx].ID)
            if ok {
                fmt.Println("Validate route has been found: ", route)
                fmt.Println("Route is like:")
                routePaint(route)
            } else {
                fmt.Println("Has not found the referred transaction")
            }
            spvRes := merkletree.SimplyPaymentValidation(primeTxs[findidx].ID, testBlock.MerkleTree.RootNode.Data, route, hashroute)
            fmt.Println("SPV result: ", spvRes, ", Want result: ", test.wants[num])
            if spvRes != test.wants[num] {
                t.Errorf("test %d find %d: SPV is not right", idx, findidx)
            }
            fmt.Println("------------------------------------------------------------------")
        }
        fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}
}

func mtGraphPaint(txContained []int)  {
	mtLayer := [][]string{}
	bottomLayer := []string{}

	for i := 0; i < len(txContained); i++ {
		bottomLayer = append(bottomLayer, strconv.Itoa(txContained[i]))
	}
	if len(bottomLayer) % 2 == 1 {
		bottomLayer = append(bottomLayer, bottomLayer[len(bottomLayer) - 1])
	}

	mtLayer = append(mtLayer, bottomLayer)

	for len(mtLayer[len(mtLayer) - 1]) != 1 {
		tempLayer := []string{}
		if len(mtLayer[len(mtLayer) - 1]) % 2 == 1 {

            
            //    {
            //        txContained: []int{0, 3, 5, 6, 7},
            //       prevBlock:   "prev4",
            //        findTX:      []int{0, 1, 6, 7},
            //        wants:       []bool{true, false, true, true},
            //    },          
            
            // fmt.Println("testestestestestsetsetestset")
            // fmt.Println(len(mtLayer[len(mtLayer) - 1]))
            // fmt.Println(mtLayer[len(mtLayer) - 1])
            // fmt.Println(mtLayer[len(mtLayer) - 1][len(mtLayer[len(mtLayer)-1]) - 1])
            // fmt.Println("testestestestestsetsetestset")
            tempLayer = append(tempLayer, mtLayer[len(mtLayer) - 1][len(mtLayer[len(mtLayer)-1]) - 1])
			mtLayer[len(mtLayer)-1][len(mtLayer[len(mtLayer)-1]) - 1] = "->"			
		}
		for i := 0; i < len(mtLayer[len(mtLayer)-1])/2; i++ {
            tempLayer = append(tempLayer, mtLayer[len(mtLayer)-1][2*i]+mtLayer[len(mtLayer)-1][2*i+1])
            
            // fmt.Println("tempLayertempLayertempLayertempLayertempLayer")
            // fmt.Println(tempLayer)
            // fmt.Println("tempLayertempLayertempLayertempLayertempLayer")
        }

        mtLayer = append(mtLayer, tempLayer)
        // fmt.Println("mtLayermtLayermtLayermtLayermtLayermtLayermtLayer")
        // fmt.Println(mtLayer)
        // fmt.Println("mtLayermtLayermtLayermtLayermtLayermtLayermtLayer")
    }

	layers := len(mtLayer)
    fmt.Println(strings.Repeat(" ", layers-1), mtLayer[layers-1][0])
	foreSpace := 0
    for i := layers - 2; i >= 0; i-- {
        var str1, str2 string
        str1 += strings.Repeat(" ", foreSpace)
        str2 += strings.Repeat(" ", foreSpace)

        for j := 0; j < len(mtLayer[i]); j++ {
            str1 += strings.Repeat(" ", i+1)
            if j%2 == 0 {
                if mtLayer[i][j] == "->" {
                    foreSpace += (i+1)*2 + 1
                    str1 = strings.Repeat(" ", foreSpace) + str1
                    str2 = strings.Repeat(" ", foreSpace) + str2
                } else {
                    str1 += "/"
                }

            } else {
                str1 += "\\"
            }
            str1 += strings.Repeat(" ", len(mtLayer[i][j])-1)
            str2 += strings.Repeat(" ", i+1)
            str2 += mtLayer[i][j]
        }
        fmt.Println(str1)
        fmt.Println(str2)
    }
}
func routePaint(route []int) {
    probe := len(route)
    fmt.Println(strings.Repeat(" ", probe) + "*")
    for i := 0; i < len(route); i++ {
        var str1, str2 string
        str1 += strings.Repeat(" ", probe)
        if route[i] == 0 {
            str1 += "/"
            probe -= 1
        } else {
            str1 += "\\"
            probe += 1
        }
        str2 += strings.Repeat(" ", probe) + "*"
        fmt.Println(str1)
        fmt.Println(str2)
    }
}


