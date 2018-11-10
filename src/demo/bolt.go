package main

import (
	"BlockChainMini/src/bolt"
	"github.com/astaxie/beego/logs"
	"fmt"
)

func main()  {

	//1.创建数据库
	db, err := bolt.Open("test.db", 0600, nil)
	if err!=nil{
		logs.Error(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		//2.通过桶的名字，找到桶
		bucket:=tx.Bucket([]byte("bucketName1"))
		//如果没有
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("bucketName1"))
			if err!=nil{
				logs.Error(err)
			}
		}

		//3.写数据
		err = bucket.Put([]byte("1"),[]byte("Hello,"))
		if err!=nil{
			logs.Error(err)
		}
		err = bucket.Put([]byte("2"),[]byte("World"))
		if err!=nil{
			logs.Error(err)
		}

		//5.读数据
		data1 := bucket.Get([]byte("1"))
		data2 := bucket.Get([]byte("2"))

		fmt.Printf("data1:%s\n",data1)
		fmt.Printf("data2:%s\n",data2)
		return nil
	})

}
