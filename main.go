package main

import (
	"goblockchain/cli"
	"os"
)

func main() {

	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
	// txPool := make([]*transaction.Transaction, 0)
	// var tempTx *transaction.Transaction
	// var ok bool
	// var property int
	// chain := blockchain.CreateBlockChain()
	// property, _ = chain.FindUTXOs([]byte("white"))
	// fmt.Println("Balance of white: ", property)

	// tempTx, ok = chain.CreateTransaction([]byte("white"), []byte("Krad"), 100)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }
	// chain.Mine(txPool)
	
	// txPool = make([]*transaction.Transaction, 0)
	// property, _ = chain.FindUTXOs([]byte("white"))
	// fmt.Println("Balance of white: ", property)

	// //做一系列交易
	// tempTx, ok = chain.CreateTransaction([]byte("Krad"), []byte("Exia"), 200) // this transaction is invalid
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// tempTx, ok = chain.CreateTransaction([]byte("Krad"), []byte("Exia"), 50)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// tempTx, ok = chain.CreateTransaction([]byte("white"), []byte("Exia"), 100)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// chain.Mine(txPool)
	// txPool = make([]*transaction.Transaction, 0)

	// property, _ = chain.FindUTXOs([]byte("white"))
	// fmt.Println("Balance of white: ", property)
	// property, _ = chain.FindUTXOs([]byte("Krad"))
	// fmt.Println("Balance of Krad: ", property)
	// property, _ = chain.FindUTXOs([]byte("Exia"))
	// fmt.Println("Balance of Exia: ", property)
	// // time.Sleep(time.Second)
	// // blockChain.AddBlock("After genesis, I have something to say.")
	// // time.Sleep(time.Second)
	// // blockChain.AddBlock("WHITE is cool!")
	// // time.Sleep(time.Second)
	// // blockChain.AddBlock("I can't wait to follow his redbook!")
	// // time.Sleep(time.Second)

	// for _, block := range chain.Blocks {
	// 	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	// 	fmt.Printf("hash: %x\n", block.Hash)
	// 	fmt.Printf("Previous hash: %x\n", block.PrevHash)
	// 	fmt.Printf("nonce: %d\n", block.Nonce)
	// 	fmt.Println("Proof of Work validation:", block.ValidatePoW())
	// }


	// //I want to show the bug at this version.
	// fmt.Println("This is the BUG:")
	// tempTx, ok = chain.CreateTransaction([]byte("Krad"), []byte("Exia"), 30)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// tempTx, ok = chain.CreateTransaction([]byte("Krad"), []byte("white"), 30)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// chain.Mine(txPool)
	// txPool = make([]*transaction.Transaction, 0)

	// for _, block := range chain.Blocks {
	// 	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	// 	fmt.Printf("hash: %x\n", block.Hash)
	// 	fmt.Printf("Previous hash: %x\n", block.PrevHash)
	// 	fmt.Printf("nonce: %d\n", block.Nonce)
	// 	fmt.Println("Proof of Work validation:", block.ValidatePoW())
	// }

	// property, _ = chain.FindUTXOs([]byte("white"))
	// fmt.Println("Balance of white: ", property)
	// property, _ = chain.FindUTXOs([]byte("Krad"))
	// fmt.Println("Balance of Krad: ", property)
	// property, _ = chain.FindUTXOs([]byte("Exia"))
	// fmt.Println("Balance of Exia: ", property)
}
