//一些全局常量
package constcoe

const (
	Difficulty = 12 //正常情况下是根据网络状况来确定难度值的，保证各节点在同一时间同一难度下进行竞争
	InitCoin =1000
	//相当于一个缓冲池，用于存储节点收集到的交易信息
	TransactionPoolFile = "./tmp/transaction_pool.data" 
	BCPath              = "./tmp/blocks"	
	BCFile              = "./tmp/blocks/MANIFEST" 	
	ChecksumLength      = 4 
	NetworkVersion      = byte(0x00) 
	Wallets             = "./tmp/wallets/" 
	WalletsRefList      = "./tmp/ref_list/" 
	UTXOSet				= "./tmp/utxoset/"
)