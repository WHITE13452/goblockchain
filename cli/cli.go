package cli

import (
	"bytes"
	"flag"
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/utils"
	"goblockchain/wallet"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

//cli.go
func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to White's tiny blockchain system, usage is as follows:")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a wallet.")
	fmt.Println("And then you can use the wallet address to create a blockchain and declare the owner.")
	fmt.Println("Make transactions to expand the blockchain.")
	fmt.Println("In addition, don't forget to run mine function after transatcions are collected.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -refname REFNAME                       ----> Creates and save a wallet. The refname is optional.")
	fmt.Println("walletinfo -refname NAME -address Address           ----> Print the information of a wallet. At least one of the refname and address is required.")
	fmt.Println("walletsupdate                                       ----> Registrate and update all the wallets (especially when you have added an existed .wlt file).")
	fmt.Println("walletslist                                         ----> List all the wallets found (make sure you have run walletsupdate first).")
	fmt.Println("createblockchain -refname NAME -address ADDRESS     ----> Creates a blockchain with the owner you input (address or refname).")
	fmt.Println("Please make sure the UTXO set init before querying the balance of a wallet.")
	fmt.Println("initutxoset                                         ----> Init all the UTXO sets of known wallets.")
	fmt.Println("balance -refname NAME -address ADDRESS              ----> Back the balance of a wallet using the address (or refname) you input.")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain.")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block.")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> Make a transaction and put it into candidate block using refname.")
	fmt.Println("mine                                                ----> Mine and add a block to the chain.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	// fmt.Println("Welcome to White's tiny blockchain system, usage is as follows:")
	// fmt.Println("--------------------------------------------------------------------------------------------------------------")
	// fmt.Println("All you need is to first create a blockchain and declare the owner.")
	// fmt.Println("And then you can make transactions.")
	// fmt.Println("--------------------------------------------------------------------------------------------------------------")
	// fmt.Println("createblockchain -address ADDRESS                   ----> Creates a blockchain with the owner you input")
	// fmt.Println("balance -address ADDRESS                            ----> Back the balance of the address you input")
	// fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain")
	// fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block")
	// fmt.Println("mine                                                ----> Mine and add a block to the chain")
	// fmt.Println("--------------------------------------------------------------------------------------------------------------")
}

func (cli *CommandLine) createWallet(refname string)  {
	newWallet := wallet.NewWallet()
	newWallet.Save()
	refList := wallet.LoadRefList()
	refList.BindRefName(string(newWallet.Address()), refname)
	refList.Save()
	fmt.Println("Succeed in creating wallet.")
}

func (cli *CommandLine) walletInfoRefName(refname string)  {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.walletInfo(address)
}

func (cli *CommandLine) walletInfo(address string)  {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	fmt.Printf("Wallet address:%x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
	fmt.Printf("Reference Name:%s\n", (*refList)[address])
}

func (cli *CommandLine) walletsUpdate()  {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("Succeed in updating wallets.")
}

func  (cli *CommandLine) walletsList()  {
	refList := wallet.LoadRefList()
	for address := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s\n", (*refList)[address])
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}

func (cli *CommandLine) sendRefName(fromRefname, toRefname string, amount int){
	refList := wallet.LoadRefList()
	FromAddress, err := refList.FindRef(fromRefname)
	utils.Handle(err)
	ToAddress, err := refList.FindRef(toRefname)
	utils.Handle(err)
	cli.send(FromAddress, ToAddress, amount)
} 

func (cli *CommandLine) createBlockChainRefName(refname string)  {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.createBlockChain(address)
}

func (cli *CommandLine) createBlockChain(address string)  {
	newChain := blockchain.InitBlockChain(utils.Address2PubHash([]byte(address)))
	newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}

func (cli *CommandLine) balanceRefName(refname string)  {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.balance(address)
}

// // 根据address统计余额
// func (cli *CommandLine) balance(address string)  {
// 	chain := blockchain.ContinueBlockChain()
// 	defer chain.Database.Close()

// 	wlt := wallet.LoadWallet(address)

// 	balance, _ := chain.FindUTXOs(wlt.PublicKey)
// 	fmt.Printf("Address:%s, Balance:%d \n", address, balance)
// }

func (cli *CommandLine) balance(address string)  {
	wlt := wallet.LoadWallet(address)
	balance := wlt.GetBalance()
	fmt.Printf("Address:%s, Balance:%d \n", address, balance)
}

func (cli *CommandLine) getBlockChainInfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iterator := chain.Iterator()
	ogprehash := chain.BackOgPreHash()
	for {
		block := iterator.Next()
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Transactions:%v\n", block.Transaction)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(block.PrevHash, ogprehash) {
			break
		}
	}

}
func (cli *CommandLine) send(from, to string, amount int)  {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	fromWallet := wallet.LoadWallet(from)
	privKey,err := utils.BuildPrivateKey(string(fromWallet.PrivateKey))
	utils.Handle(err)
	tx, ok := chain.CreateTransaction(fromWallet.PublicKey, utils.Address2PubHash([]byte(to)), amount, *privKey)
	if !ok {
		fmt.Println("Failed to create Transaction")
		return
	}

	tp := blockchain.CreateTransactionPool()
	tp.AddTransaction(tx)
	tp.SaveFile()
	fmt.Println("Success!")
}

func (cli *CommandLine) mine()  {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("Finish Mining")
	newblock := chain.GetCurrentBlock()
	refList := wallet.LoadRefList()
	for address := range *refList {
		wlt := wallet.LoadWallet(address)
		wlt.ScanBlock(newblock)
	}
	fmt.Println("Finish Updating UTXO sets")
}

func (cli *CommandLine) validateArgs()  {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) initUtxosSet()  {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	refList := wallet.LoadRefList()
	for addr, _ := range *refList {
		wlt := wallet.LoadWallet(addr)
		utxoSet := wlt.CreateUTXOSet(chain)
		utxoSet.DB.Close()
	}
	fmt.Println("Succeed in initializing UTXO sets.")
}

func (cli *CommandLine) Run()  {
	cli.validateArgs()
	
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)   
	walletInfoCmd := flag.NewFlagSet("walletinfo", flag.ExitOnError)       
	walletsUpdateCmd := flag.NewFlagSet("walletsupdate ", flag.ExitOnError)
	walletsListCmd := flag.NewFlagSet("walletslist", flag.ExitOnError)     
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	//dont need it in a real blockchain
	sendByRefNameCmd := flag.NewFlagSet("sendbyrefname", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)
	initUtxoCmd := flag.NewFlagSet("initutxoset", flag.ExitOnError)

	createWalletRefName := createWalletCmd.String("refname", "", "The refname of the wallet, and this is optimal") 
	walletInfoRefName := walletInfoCmd.String("refname", "", "The refname of the wallet")                          
	walletInfoAddress := walletInfoCmd.String("address", "", "The address of the wallet")                          
	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	createBlockChainByRefNameOwner := createBlockChainCmd.String("refname", "", "The name refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")
	balanceRefName := balanceCmd.String("refname", "", "Who needs to get balance amount")
	sendByRefNameFrom := sendByRefNameCmd.String("from", "", "Source refname")           
	sendByRefNameTo := sendByRefNameCmd.String("to", "", "Destination refname")          
	sendByRefNameAmount := sendByRefNameCmd.Int("amount", 0, "Amount to send")           
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet": 
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletinfo": 
		err := walletInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletsupdate": 
		err := walletsUpdateCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletslist": 
		err := walletsListCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "blockchaininfo":
		err := getBlockChainInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)
	
	case "sendbyrefname": 
		err := sendByRefNameCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "initutxoset":
		err := initUtxoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(*createWalletRefName)
	}
	if walletInfoCmd.Parsed() {
		if *walletInfoAddress == "" {
			if *walletInfoRefName == "" {
				walletInfoCmd.Usage()
				runtime.Goexit()	
			} else {
				cli.walletInfoRefName(*walletInfoRefName)
			}
		} else {
			cli.walletInfo(*walletInfoAddress)
		}
	}

	if walletsUpdateCmd.Parsed() {
		cli.walletsUpdate()
	}

	if walletsListCmd.Parsed() {
		cli.walletsList()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			if *createBlockChainByRefNameOwner == "" {
				createBlockChainCmd.Usage()
				runtime.Goexit()
			} else {
				cli.createBlockChainRefName(*createBlockChainByRefNameOwner)
			}
		} else {
			cli.createBlockChain(*createBlockChainOwner)
		}
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			if *balanceRefName == "" {
				balanceCmd.Usage()
				runtime.Goexit()
			} else {
				cli.balanceRefName(*balanceRefName)
			}
		} else {
			cli.balance(*balanceAddress)
		}
	}

	if sendByRefNameCmd.Parsed() {
		if *sendByRefNameFrom == "" || *sendByRefNameTo == "" || *sendByRefNameAmount <= 0 {
			sendByRefNameCmd.Usage()
			runtime.Goexit()
		}
		cli.sendRefName(*sendByRefNameFrom, *sendByRefNameTo, *sendByRefNameAmount)
	}


	if sendCmd.Parsed() {
		if *sendFromAddress == "" || *sendToAddress == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}

	if getBlockChainInfoCmd.Parsed() {
		cli.getBlockChainInfo()
	}

	if initUtxoCmd.Parsed() {
		cli.initUtxosSet()
	}

	if mineCmd.Parsed() {
		cli.mine()
	}

	
}