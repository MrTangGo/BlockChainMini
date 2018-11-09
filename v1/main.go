package main

import (
	"fmt"
)


func main() {
	bc := NewBlockChain()
	bc.AddBlock("Hello World")
	bc.AddBlock("Good")
	bc.AddBlock("asdf")
	bc.AddBlock("Gos234od")
	bc.AddBlock("sdfsd3445")
	bc.AddBlock("fdsasd224t5")

	for i, block := range bc.blocks {
		fmt.Printf("============== 区块高度：%d ===========\n", i)
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Data)
	}
}
