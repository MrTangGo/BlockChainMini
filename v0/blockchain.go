package main

//定义一个区块链结构，用数组来实现，链接的时候使用append即可
type BlockChain struct {
	blocks []*Block
}

//定义一个创建区块链的方法
func NewBlockChain() *BlockChain {
	//在创建区块链的时候，添加一个创世块genesisBlock
	genesisBlock := NewBlock(genesisInfo, []byte{})
	blockChain := BlockChain{blocks: []*Block{genesisBlock}}
	return &blockChain
}

func (bc *BlockChain) AddBlock(data string) {
	//根据数组的下标找到最后一个区块，获取前区块哈希值
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash

	//创建新的区块，并且添加到区块链
	block := NewBlock(data, prevHash)
	bc.blocks = append(bc.blocks, block)
}

