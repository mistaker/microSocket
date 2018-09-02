# microSocket
这是一款十分适合学习的go语言socket框架

能够非常简单就实现一个服务端

```go
package main

import (
	"fmt"
	"log"
	msf "microSocket"
	"strconv"
)

var ser = msf.NewMsf()

//框架逻辑
//---------------------------------------------------------------------
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
我也对该框架做了源码分析  [传送](https://www.jianshu.com/p/49974703cf3e)