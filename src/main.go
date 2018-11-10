package main



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
	bc.PrintChain()
}
