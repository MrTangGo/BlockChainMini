package main

import (
	"math/big"
	"fmt"
	"crypto/sha256"
	"bytes"
)

//定义一个工作量证明结构
type ProofOfWork struct{
	block Block	//区块数据
	target *big.Int	//难度系数
}

//添加创建POW的函数
func NewProofOfWork(block Block) *ProofOfWork {
	//新建一个工作量证明结构体(难度系数先不给，自己写一个)
	pow := ProofOfWork{
		block: block,
	}

	//自定义的难度值，先写成固定值
	targetString := "0010000000000000000000000000000000000000000000000000000000000000"
	bigIntTmp := big.Int{}
	bigIntTmp.SetString(targetString, 16)
	pow.target = &bigIntTmp

	return &pow
}

func (pow *ProofOfWork)Run()([]byte,uint64){

	fmt.Printf("pow run...\n")

	var Nonce uint64
	var hash [32]byte

	for;;{
		//计算出添加了随机值的hash
		hash = sha256.Sum256(pow.prepareData(Nonce))
		//把hash转化为一个big.int
		tTmp := big.Int{}
		tTmp.SetBytes(hash[:])

		//与难度值比较，符合就退出循环
		if tTmp.Cmp(pow.target) == -1 {
			fmt.Printf("found hash : %x, %d\n", hash, Nonce)
			break
		}else {
			Nonce++
		}
	}
	return hash[:],Nonce
}

func (pow *ProofOfWork)prepareData(Nonce uint64)[]byte  {
	block := pow.block

	tmp := [][]byte{
		Uint2Byte(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		Uint2Byte(block.TimeStamp),
		Uint2Byte(block.Difficulty),
		Uint2Byte(block.Nonce),
		block.Data,
	}

	data:=bytes.Join(tmp,[]byte{})
	return data
}

