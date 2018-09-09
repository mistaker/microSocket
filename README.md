# microSocket
这是一款十分适合学习的go语言socket框架

能够非常简单就实现一个服务端

```go
package main

import (
	"fmt"
	"log"
	msf "microSocket"
	"net"
	"strconv"
)

var ser = msf.NewMsf(&event{})

//框架事件
type event struct {
}

//客户端握手成功事件
func (this event) OnHandel(fd uint32, conn net.Conn) bool {
	log.Println(fd, "链接成功类")
	return true
}

//断开连接事件
func (this event) OnClose(fd uint32) {
	log.Println(fd, "链接断开类")
}

//接收到消息事件
func (this event) OnMessage(fd uint32, msg map[string]string) bool {
	return true
}

//---------------------------------------------------------------------
//框架业务逻辑
type Test struct {
}

func (this Test) Default() {
	fmt.Println("is default")
}

func (this Test) BeforeRequest(data map[string]string) bool {
	log.Println("before")
	return true
}

func (this Test) AfterRequest(data map[string]string) {
	log.Println("after")
}

func (this Test) Hello(data map[string]string) {
	fd, _ := strconv.Atoi(data["fd"])
	log.Println("收到消息了")
	ser.SessionMaster.WriteByid(uint32(fd), "Hello")
}

//---------------------------------------------------------------------

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)
	ser.EventPool.Register("test", &Test{})
	ser.Listening(":8565")
}

```
其他代码客户端连接成功后发送以下字节数组就能 有反应
[72 101 97 100 101 114 0 0 0 32 110 97 109 101 58 106 100 124 109 111 100 117 108 101 58 116 101 115 116 124 97 99 116 105 111 110 58 72 101 108 108 111
我也对该框架做了源码分析  [传送](https://www.jianshu.com/p/49974703cf3e)