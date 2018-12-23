# microSocket
这是一款十分适合学习的go语言socket框架
socket 和websocket 一键切换 业务代码完全不用改


能够非常简单就实现一个服务端

```go
package main

import (
	"log"
	msf "microSocket"
	"net"
)

var ser = msf.NewMsf(&msf.CommSocket{})

//框架事件
//----------------------------------------------------------------------------------------------------------------------
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
	log.Println("这个是接受消息事件",msg)
	return true
}
//----------------------------------------------------------------------------------------------------------------------
//框架业务逻辑
type Test struct {
}

func (this Test) Default(fd uint32,data map[string]string) bool {
	log.Println("default")
	return true
}

func (this Test) BeforeRequest(fd uint32,data map[string]string) bool {
	log.Println("before")
	return true
}

func (this Test) AfterRequest(fd uint32,data map[string]string) bool{
	log.Println("after")
	return true
}

func (this Test) Hello(fd uint32,data map[string]string) bool {
	log.Println("收到消息了")
	log.Println(data)
	ser.SessionMaster.WriteByid(fd,[]byte("hehehehehehehehe"))
	return true
}


//----------------------------------------------------------------------------------------------------------------------
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)
	ser.EventPool.RegisterEvent(&event{})
	ser.EventPool.RegisterStructFun("test", &Test{})
	ser.Listening(":8565")
}



```
客户端连接成功后发送"module:test|action:Hello"就能 触发相应模块事件

我也对该框架做了源码分析  [传送](https://www.jianshu.com/p/49974703cf3e)