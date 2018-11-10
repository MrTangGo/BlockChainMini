package main

import (
	"BlockChainMini/src/bolt"
	"github.com/astaxie/beego/logs"
	"fmt"
	"os"
	"time"
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

//定义一个创建区块链的方法
func NewBlockChain() *BlockChain {
	var lastHash []byte
	//1.创建数据库
	db, err := bolt.Open(blockChainName, 0600, nil)
	if err!=nil{
		logs.Error(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		//2.通过桶的名字，找到桶
		bucket:=tx.Bucket([]byte(blockBucket))
		//如果没有
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err!=nil{
				logs.Error(err)
			}
			//3.写数据
			//创建区块的时候添加一个创世区块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			err = bucket.Put([]byte(genesisBlock.Hash),  genesisBlock.Serialize() /*将区块序列化成字节流*/)
			if err!=nil{
				logs.Error(err)
			}
			//更新最后一个区块的哈希
			err = bucket.Put([]byte(lastHashKey),genesisBlock.Hash)
			if err!=nil{
				logs.Error(err)
			}

			//更新内存中的最后一个区块的hash值
			lastHash = genesisBlock.Hash
		}else{
			//如果bucket存在，直接Get回来就可以
			lastHash = bucket.Get([]byte(lastHashKey))
		}
		return nil
	})

	return &BlockChain{db,lastHash}
}

//添加区块
func (bc *BlockChain) AddBlock(data string) {
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
		newBlock := NewBlock(data, prevBlockHash)

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
		fmt.Printf("Data :%s\n", block.Data)
		pow := NewProofOfWork(*block)
		fmt.Printf("IsValid : %v\n\n", pow.IsValid())

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("打印结束!")
			break
		}
	}
}

