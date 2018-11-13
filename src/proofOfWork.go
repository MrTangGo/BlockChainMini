package main

import (
	"math/big"
	"fmt"
	"bytes"
	"crypto/sha256"
)

const difficulty = 20

//定义一个工作量证明结构
type ProofOfWork struct {
	block  Block    //区块数据
	target *big.Int //难度系数
}

//创建一个POW工作量证明函数
func NewProofOfWork(block Block) *ProofOfWork {
	//新建一个工作量证明结构体(难度系数先不给，自己写一个)
	pow := ProofOfWork{
		block: block,
	}


	//自定义的难度值，先写成固定值
	//targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//var targetInt big.Int
	//targetInt.SetString(targetStr, 16)
	//pow.target = &targetInt

	//v2.5
	targetInt:=big.NewInt(1)//new完之后就是指针了
	//targetInt.Lsh(targetInt,256)
	//targetInt.Rsh(targetInt,20)
	targetInt.Lsh(targetInt,256-difficulty)
	pow.target = targetInt

	return &pow
}

//
func (pow *ProofOfWork) Run() ([]byte, uint64) {

	fmt.Println("===================挖矿中===================")

	var nonce uint64
	var curtBlockHash [32]byte

	for {
		//通过prepareDate获取挖矿要用到的相数据，再挖矿
		info := pow.prepareDate(nonce)
		curtBlockHash = sha256.Sum256(info)

		//判断hash:1.hash->big.int 2..Cmp
		var currentHashInt big.Int
		currentHashInt.SetBytes(curtBlockHash[:])

		if currentHashInt.Cmp(pow.target) == -1 {
			fmt.Printf("找到了符合条件的Hash值:%x,随机数：%d\n", curtBlockHash, nonce)
			break
		} else {
			nonce++
		}
	}

	return curtBlockHash[:], nonce
}

func (pow *ProofOfWork) prepareDate(nonce uint64) []byte {
	block := pow.block
	//需要hash的数据
	curtBlockArray := [][]byte{
		Uint2Byte(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		Uint2Byte(block.TimeStamp),
		Uint2Byte(block.Difficulty),
		Uint2Byte(nonce),
	}
	info := bytes.Join(curtBlockArray, []byte{})
	return info
}

func (pow *ProofOfWork) IsValid() bool {
	info := pow.prepareDate(pow.block.Nonce)
	hash := sha256.Sum256(info)

	var tempInt big.Int
	tempInt.SetBytes(hash[:])

	return tempInt.Cmp(pow.target) == -1
}
