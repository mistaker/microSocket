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
