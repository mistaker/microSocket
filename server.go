package microSocket

import (
	"log"
	"net"
)

type SocketTypes interface{
	ConnHandle(msf *Msf,sess *Session)
	Pack(data []byte)[]byte
}

type MsfEventer interface {
	OnHandel(fd uint32, conn net.Conn) bool
	OnClose(fd uint32)
	OnMessage(fd uint32, msg map[string]string) bool
}

type Msf struct {
	EventPool     *RouterMap
	SessionMaster *SessionM
	MsfEvent      MsfEventer
	SocketType        SocketTypes
}

func NewMsf(msfEvent MsfEventer,socketType SocketTypes) *Msf {
	msf := &Msf{
		EventPool:     NewRouterMap(),
		MsfEvent:      msfEvent,
		SocketType   :socketType,
	}
	msf.SessionMaster = NewSessonM(msf)
	return msf
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

		go this.SocketType.ConnHandle(this,sess)
	}
}



