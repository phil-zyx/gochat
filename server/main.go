package main

import (
	"github.com/gochat/server/conf"
	"github.com/gochat/server/network"
	"github.com/gochat/server/redisop"
	"github.com/gochat/server/wordfilter"
)

// Init 全局初始化
func init() {
	wordfilter.Init("server/conf/sensitive.txt")
}

func main() {
	tcpServer := network.NewChatServer(conf.Addr)
	tcpServer.RedisCli = redisop.InitClient()
	tcpServer.Start()
}