package main

import (
	"BlockChainMini/src/bolt"
	"github.com/astaxie/beego/logs"
	"fmt"
	"os"
	"time"
	"log"
)

const blockChainName = "blockChain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

//定义一个区块链结构，用数组来实现，链接的时候使用append即可
type BlockChain struct {
	//数据库的句柄
	Db *bolt.DB
	//最后一个区块的哈希
	lastHash []byte
}

//创建一个新的区块链
func CreateBlockChain(address string) *BlockChain {
	if isDbExist() {
		fmt.Printf("区块链已经存在!\n")
		os.Exit(1)
	}
	var lastHash []byte
	db, err := bolt.Open(blockChainName, 0600, nil)

	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//2. 找到我们的桶，通过桶的名字
		// Returns nil if the bucket does not exist.
		bucket := tx.Bucket([]byte(blockBucket))
		//如果没有找到，先创建
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}

			//3. 写数据
			//在创建区块链的时候，添加一个创世块genesisBlock
			coinbase := NewCoinbaseTx(address, genesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})

			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化成字节流*/)
			if err != nil {
				log.Panic(err)
			}

			//一定要记得更新"lastHashKey" 这个key对应的值，最后一个区块的哈希
			err = bucket.Put([]byte(lastHashKey), genesisBlock.Hash)

			//更新内存中最后区块哈希值
			lastHash = genesisBlock.Hash
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

//定义一个创建区块链的方法
func NewBlockChain() *BlockChain {
	var lastHash []byte
	//1.创建数据库
	db, err := bolt.Open(blockChainName, 0600, nil)
	if err!=nil{
		logs.Error(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		//2.通过桶的名字，找到桶
		bucket:=tx.Bucket([]byte(blockBucket))
		//如果没有
		if bucket == nil {
			fmt.Printf("获取区块时不应该为空")
		}

		lastHash = bucket.Get([]byte(lastHashKey))

		return nil
	})

	return &BlockChain{db,lastHash}
}

//添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//最后一个区块的哈希值,也就是新区块的前哈希值
	prevBlockHash := bc.lastHash

	// 更新数据库
	bc.Db.Update(func(tx *bolt.Tx) error {
		//1. 找到bucket
		bucket := tx.Bucket([]byte(blockBucket))
		//2. 判断有没有，
		if bucket == nil {
			fmt.Printf("添加区块时，bucket不应为空，请检查!")
			//没有， 直接报错退出
			os.Exit(1)
		}
		//有，写入数据
		newBlock := NewBlock(txs, prevBlockHash)

		//更新数据库
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		bucket.Put([]byte(lastHashKey), newBlock.Hash)

		//更新内存
		bc.lastHash = newBlock.Hash
		return nil
	})
}

type Iterator struct {
	Db          *bolt.DB //来自于区块链
	currentHash []byte   //随着移动改变
}

//创建一个迭代器, 最初指向最后一个区块
func (bc *BlockChain) NewIterator() *Iterator {
	return &Iterator{Db: bc.Db, currentHash: bc.lastHash}
}

func (it *Iterator) Next() *Block {
	var block *Block
	it.Db.View(func(tx *bolt.Tx) error {
		//找到bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Printf("遍历区块时，bucket不应为空，请检查!")
			os.Exit(1)
		}

		//读取数据：currentHash
		blockTmp := bucket.Get(it.currentHash)
		block = Deserialize(blockTmp)

		//currentHash左移
		it.currentHash = block.PrevBlockHash

		return nil
	})

	return block
}


func (bc *BlockChain) PrintChain() {
	it := bc.NewIterator()

	for ; ; {

		block := it.Next()

		fmt.Printf("===============================\n")
		fmt.Printf("Version :%d\n", block.Version)
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("MerkeRoot :%x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp: %s\n", timeFormat)
		fmt.Printf("Difficulty :%d\n", block.Difficulty)
		fmt.Printf("Nonce :%d\n", block.Nonce)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Transactions[0].TXInputs[0].Sig)
		pow := NewProofOfWork(*block)
		fmt.Printf("IsValid : %v\n\n", pow.IsValid())

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("打印结束!")
			break
		}
	}
}

//判断区块链文件是否存在
func isDbExist() bool {
	if _, err := os.Stat(blockChainName); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *BlockChain) FindMyUtxos(address string) []TXOutput{

	it := bc.NewIterator()

	var UTXOInfos []TXOutput

	spentUTXO :=make(map[string][]int64)

	//遍历区块
	for ;;{
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
			OUTPUT:
			//遍历output
			for i /*0, 1, 2, 3*/ , output := range tx.TXOutputs {
				fmt.Printf("当前索引为：%d\n", i)

				if address == output.PubKeyHash {

					//checkout if this address is used
					key := string(tx.TXID)
					//检查一下这个output是否已经被用过了
					if len(spentUTXO[key]) /*[]int64*/ != 0 {
						fmt.Printf("当前交易里面有%s消耗过的output\n", address)
						//spentUTXOs[0x222] = []int64{0}
						//spentUTXOs[0x333] = []int64{0} //中间值
						//spentUTXOs[0x333] = []int64{0, 1}
						for _, j /*0, 1*/ := range spentUTXO[key] {
							if int64(i) == j {
								fmt.Printf("i==j,这个output被消耗了，不统计\n")
								continue OUTPUT
							}
						}
					}
					UTXOInfos = append(UTXOInfos, output)
					fmt.Printf("找到了一个属于%s的output\n", address)
				}
			}

			//遍历input，找到这个address已经消耗过得output，标识出来, 在遍历output前检测，过滤
			for _,input:=range tx.TXInputs{
				if input.Sig == address{
					fmt.Printf("%s 已经消耗的output:%d\n",address,input.Index)
					key :=string(input.TxId)
					spentUTXO[key]= append(spentUTXO[key],input.Index)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXOInfos
}

func (bc *BlockChain) GetBalance(address string) {

	//找到所有属于address的utxo的数组
	utxos := bc.FindMyUtxos(address)

	//总额
	total := 0.0

	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("%s 的余额为: %f\n", address, total)
}

