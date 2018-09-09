package microSocket

import (
	"fmt"
	"log"
	"microSocket/util"
	"net"
	"time"
)

type MsfEventer interface {
	OnHandel(fd uint32, conn net.Conn) bool
	OnClose(fd uint32)
	OnMessage(fd uint32, msg map[string]string) bool
}

type Msf struct {
	EventPool     *RouterMap
	SessionMaster *SessionM
	MsfEvent      MsfEventer
}

func NewMsf(msfEvent MsfEventer) *Msf {
	return &Msf{
		SessionMaster: NewSessonM(),
		EventPool:     NewRouterMap(),
		MsfEvent:      msfEvent,
	}
}

func (this *Msf) Listening(address string) {
	tcpListen, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}

	fd := uint32(0)
	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//调用握手事件
		if this.MsfEvent.OnHandel(fd, conn) == false {
			continue
		}

		sess := NewSession(fd, conn)
		this.SessionMaster.SetSession(fd, sess)
		fd++

		go this.connHandle(conn, sess)
	}
}

func (this *Msf) connHandle(conn net.Conn, sess *session) {
	defer func() {
		this.SessionMaster.DelSessionById(sess.id)
		//调用断开链接事件
		this.MsfEvent.OnClose(sess.id)
		conn.Close()
	}()
	var errs error
	tempBuff := make([]byte, 0)
	readBuff := make([]byte, 14)
	data := make([]byte, 20)
	//设置最迟期限（心跳包）
	conn.SetReadDeadline(time.Now().Add(time.Duration(8) * time.Second))
	for {
		n, err := conn.Read(readBuff)
		//设置最迟期限（心跳包）
		conn.SetReadDeadline(time.Now().Add(time.Duration(8) * time.Second))
		if err != nil {
			return
		}
		tempBuff = append(tempBuff, readBuff[:n]...)
		tempBuff, data, errs = Depack(tempBuff)
		if errs != nil {
			log.Println(errs)
			return
		}

		if len(data) == 0 {
			continue
		}
		//把请求的到数据转化为map
		requestData := util.String2Map(string(data))
		if requestData["module"] == "" || requestData["action"] == "" || this.EventPool.ModuleExit(requestData["module"]) == false {
			log.Println("not find module ", requestData)
			continue
		}
		requestData["fd"] = fmt.Sprintf("%d", sess.id)

		//调用接收消息事件
		if this.MsfEvent.OnMessage(sess.id, requestData) == false {
			return
		}
		//路由
		this.EventPool.Hook(requestData["module"], requestData["action"], requestData)

	}
}
