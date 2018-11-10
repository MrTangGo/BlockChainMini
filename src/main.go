package main

import (
	"fmt"
	"time"
)

//
//1. 定义结构
//1. 前区块哈希
//2. 当前区块哈希
//3. 数据
//2. 创建区块
//3. 生成哈希
//4. 引入区块链
//5. 添加区块
//6. 重构代码

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Hello Itcast")
	bc.AddBlock("Hello 航头")
	bc.AddBlock("Hello 航头1")

	for i, block := range bc.blocks {
		fmt.Printf("============== 区块高度：%d ===========\n", i)
		fmt.Printf("Version:%x\n",block.Version)
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("MerkelRoot:%x\n",block.MerkelRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp: %s\n", timeFormat)

		fmt.Printf("Difficulty:%d\n",block.Difficulty)
		fmt.Printf("Nonce:%d\n",block.Nonce)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Data)
		//校验函数
		pow := NewProofOfWork(*block)
		fmt.Printf("IsVald:%t\n",pow.IsValid())
	}
}
