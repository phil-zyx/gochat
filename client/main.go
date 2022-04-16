package main

import (
	"encoding/json"
	"fmt"
	"github.com/gochat/client/client"
	"github.com/gochat/proto"
	"log"
)

func main() {
	var (
		key      int
		loop     = true
		userName string
		userID   int64
		cli      *client.Client
	)

	for loop {
		log.Println("----------------Welcome to the chat room--------------")
		log.Println("Select the options:")
		log.Print("1、Sign in")
		log.Print("2、Sign up")
		log.Print("3、Exit the system")

		// get user input
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			log.Println("sign In Please")
			log.Print("input player name:")
			fmt.Scanf("%s\n", &userName)
			log.Print("input player id(when id == 0,means create player):")
			fmt.Scanf("%s\n", &userID)
			cli = client.NewClient()
			cli.Dial("127.0.0.1:8888")
			loginMsg := proto.LoginReq{PlayerName: userName, PlayerID: userID}
			data, _ := json.Marshal(loginMsg)
			msg := proto.Message{Type: proto.MsgTypeLogin, Data: string(data)}
			msgData, _ := json.Marshal(msg)
			cli.SendMsg(msgData)
			go cli.ReadCmdArgs()
			for {
				client.ShowHandleAfterLogin(cli)
			}
		case 2:
			log.Println("Logout...")
			if cli != nil {
				cli.Close()
			}
		case 3:
			log.Println("Exit...")
			loop = false
		default:
			log.Println("Select is invalid!")
		}
	}
}
