package main

import (
	"crypto/sha256"
	"time"
	"bytes"
	"encoding/binary"
	"log"
)

//data , prevHash, Hash

const genesisInfo = "2009年1月3日，财政大臣正处于实施第二轮银行紧急援助的边缘"

type Block struct {
	Version       uint64 //版本号
	PrevBlockHash []byte //前区块哈希值
	MerkelRoot []byte //这是一个哈希值，后面v5用到
	TimeStamp uint64 //时间戳，从1970.1.1到现在的秒数
	Difficulty uint64 //通过这个数字，算出一个哈希值：0x00010000000xxx
	Nonce uint64 // 这是我们要找的随机数，挖矿就找证书
	Hash []byte //当前区块哈希值, 正常的区块不存在，我们为了方便放进来
	Data []byte //数据本身，区块体，先用字符串表示，v4版本的时候会引用真正的交易结构
}

func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:       00,
		PrevBlockHash: prevHash,
		MerkelRoot:    []byte{}, //先填写为空
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulty:    difficulty,
		Nonce:         0,        //目前不挖矿，随便写一个值
		Hash:          []byte{}, //见SetHash函数
		Data:          []byte(data),
	}

	//通过工作量证明的方法得到hash与随机数
	pow:=NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash=hash
	block.Nonce=nonce


	return &block
}

func (block *Block) SetHash() {
	//将区块的各个字段拼接成一个[]byte{}
	//var info []byte
	//info = append(info, Uint2Byte(block.Version)...)
	//info = append(info, block.PrevBlockHash...)
	//info = append(info, block.MerkelRoot...)
	//info = append(info, Uint2Byte(block.TimeStamp)...)
	//info = append(info, Uint2Byte(block.Difficulty)...)
	//info = append(info, Uint2Byte(block.Nonce)...)
	//info = append(info, block.Data...)

	//使用Join代替append
	bytesArray := [][]byte{
		Uint2Byte(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		Uint2Byte(block.TimeStamp),
		Uint2Byte(block.Difficulty),
		Uint2Byte(block.Nonce),
		block.Data,
	}

	info := bytes.Join(bytesArray, []byte{})

	//对区块的数据进行哈希运算，返回[32]byte
	hash := sha256.Sum256(info)
	block.Hash = hash[:]
}

//将uint转换成[]byte
func Uint2Byte(num uint64) []byte {
	var buffer bytes.Buffer

	//这是一个序列化的过程, 将num转换成buffer字节流
	err := binary.Write(&buffer, binary.BigEndian, &num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
