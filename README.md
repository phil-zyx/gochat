# 一个简易的聊天服务

## 需求

1. 一个聊天系统，服务器与客户端
2. 单元测试
3. 功能详细：

```
1.通信基于TCP/IP协议。
2.玩家管理: 1.登陆第一次不存在玩家，则自动创建。 2.玩家名字唯一,重复则提示重试。
3.房间管理: 1.登陆之后，加入聊天室开始聊天，并下发最新的50条聊天记录。 2.可以在不同聊天房间进行切换，切换后行为重复登陆后相同处理。
4.脏词过滤:
1.用 【*】号替换脏词，最长匹配。例如 “hellboy” -> “****boy”
*附:脏词库 https://github.com/CloudcadeSF/google-profanity-words/blob/main/data/list.txt
5.GM指令:
1. /stats [username]
打印出玩家的登陆时间，在线时长，房间号。不限制输出格式。 2. /popular n (n为房间Id)
打印出最近10分钟内发送频率最高的词。
```

## 代码分析

## Server

- network：负责TCP server 的创建及消息监听
- player：玩家管理，创建获取玩家以及玩家的一些个人信息
- chatmgr：聊天室管理
- redisop：redis 存历史消息，一些存取函数
- wordfliter：脏词过滤

## Client

- for 循环监听消息及指令输入

## Proto

- 自定义C/S协议，建议使用 protobuf

## 运行

- ```
  cd project
  ```

- ```
  go run server/main.go
  ```

- ```
  go run client/main.go
  ```

## 一些说明

- server 的数据没有落地实现，mysql 或 mongodb 任意实现都有成熟的库（没实现是因为懒）。
- 采用 redis 来存放聊天记录，当做缓存，方便存取，用到了`"github.com/go-redis/redis"`
- GM 指令没有给服务器单独设定输入接口，所以直接通过 client 消息调用实现需求，将数据打印在日志里面。
- 拉取聊天室最近五十条消息，通过设定消息索引，chatroom 管理消息索引。
- 高频词汇：通过给每条消息的时间戳，从最近依次拉取时间范围内的消息，进行分词，获得高频词。