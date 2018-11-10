package main

import (
	"fmt"
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
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Data)
	}
}
