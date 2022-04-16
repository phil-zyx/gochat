package network

import (
	"github.com/go-redis/redis"
	"github.com/gochat/server/chatmgr"
	"github.com/gochat/server/player"
	"log"
	"net"
)

type ChatServer struct {
	Addr      string // IP地址
	ln        net.Listener
	ChatMgr   *chatmgr.ChatMgr
	PlayerMgr *player.Mgr
	RedisCli  *redis.Client
}

func NewChatServer(addr string) *ChatServer {
	return &ChatServer{
		Addr: addr,
		ChatMgr: chatmgr.NewMgr(),
		PlayerMgr: player.NewMgr(),
	}
}

func (cs *ChatServer) Start() {
	// 开启监听
	ln, err := net.Listen("tcp", cs.Addr)
	defer ln.Close()
	if err != nil {
		log.Fatalf("ChatServer Start Err:%v", err)
		return
	}
	log.Printf("tcp server listening on: %v", ln.Addr())
	// 保留监听器
	cs.ln = ln
	cs.run()
}

// run 启动服务器
func (cs *ChatServer) run() {
	for {
		// 接收连接
		conn, err := cs.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("accept error: %v", err)
				continue
			}
			return
		}
		// 消息处理
		go cs.AcceptMsgHandle(conn)
	}
}

// Close 关闭TCPServer
func (cs *ChatServer) Close() {
	cs.ln.Close()
}
