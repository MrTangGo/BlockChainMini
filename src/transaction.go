package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type TXInput struct {
	//1. 引用的交易id
	TxId []byte
	//2. 引用的output的索引
	Index int64
	//3. 解锁脚本
	Sig string //用地址来代替
}

type TXOutput struct {
	//1. 金额
	Value float64 //一定要大写
	//2. 锁定脚本
	PubKeyHash string //也用地址代替，只要比对地址是否相同，就认为可以解锁
}


type Transaction struct {
	TXID []byte
	//多个输入
	TXInputs []TXInput
	//多个输出
	TXOutputs []TXOutput
}

func (tx *Transaction) SetTxID() {
	//使用gob编码，生成交易的哈希

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

func NewCoinbaseTx(miner, data string) *Transaction {
	//挖矿交易的特点， 没有输入， 只有输出
	input := TXInput{nil, -1, data}
	output := TXOutput{12.5, miner}

	tx := Transaction{nil, []TXInput{input}, []TXOutput{output}}
	tx.SetTxID()

	return &tx
}






