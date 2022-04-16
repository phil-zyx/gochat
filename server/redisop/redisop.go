package redisop

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gochat/proto"
	"github.com/gochat/server/conf"
	"github.com/gochat/server/util"
	"strconv"
	"strings"
)

var messageKeyPre string = "M:"
func InitClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Redis Ping() no respond!")
		panic(err)
	}
	return client
}

// GetRecentMessage .
func GetRecentMessage(client *redis.Client, roomName string, msgIDMin, msgIDMax int64) []proto.ChatMessage {
	var ret []proto.ChatMessage
	for id := msgIDMin; id <= msgIDMax; id++ {
		key := messageKeyPre + roomName + strconv.FormatInt(id, 10)
		key = strings.TrimSpace(key)

		fmt.Println("GetMessage key =", key)
		messages, err := client.HGetAll(key).Result()
		if err == redis.Nil {
			continue
		}
		detail := messages["detail"]
		msg := proto.ChatMessage{}
		err = json.Unmarshal([]byte(detail), &msg)
		if err == nil {
			ret = append(ret, msg)
		}
	}
	return ret
}

// GetLast10MinMostPopularWords .
func GetLast10MinMostPopularWords(client *redis.Client, roomName string, msgID int64, min int32) string {
	set := make(map[string]int32)
	var max int32
	popular := ""
	if msgID <= 0 {
		return popular
	}

	for i := msgID; i > 0; i-- {
		key := messageKeyPre + roomName + strconv.FormatInt(i, 10)
		key = strings.TrimSpace(key)
		messages, err := client.HGetAll(key).Result()
		if err == redis.Nil {
			continue
		}
		detail := messages["detail"]
		msg := proto.ChatMessage{}
		err = json.Unmarshal([]byte(detail), &msg)
		if err != nil {
			continue
		}
		if msg.CreateTs+int64(min*util.Minute) < util.NowTs() {
			break
		}
		for _, word := range strings.Split(msg.Content, " ") {
			set[word]++
			if set[word] > max {
				popular = word
				max = set[word]
			}
		}
	}
	return popular
}

// AddMessage 存入消息
func AddMessage(client *redis.Client, msg proto.ChatMessage) bool {
	key := messageKeyPre + msg.RoomName + strconv.FormatInt(msg.ID, 10)
	key = strings.TrimSpace(key)

	_, err := client.HSet(key, "id", msg.ID).Result()
	fmt.Println("HSet:", key)
	if err != nil {
		fmt.Println("redisop AddMessage HSet msg Error:", err.Error())
		return false
	}
	detail, err1 := json.Marshal(msg)
	if err1 != nil {
		return false
	}
	client.HSet(key, "detail", detail)
	return true

}
