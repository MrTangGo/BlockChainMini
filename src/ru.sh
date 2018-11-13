#!/bin/bash
rm ./*.db
rm blockchain
go build -o blockchain *.go
./blockchain createBlockChain 'Tom'