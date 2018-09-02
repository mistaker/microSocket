package microSocket

import (
	"log"
	"net"
	"sync"
	"time"
)

//-------------------------------------------一个session代表一个连接------------------------------------------
type session struct {
	id    uint32
	con   net.Conn
	times int64
	lock  sync.Mutex
}

func NewSession(id uint32, con net.Conn) *session {
	return &session{
		id:    id,
		con:   con,
		times: time.Now().Unix(),
	}
}

func (this *session) Write(msg string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	_ ,errs := this.con.Write([]byte(msg))
	return errs
}

func (this *session)Close(){
	this.con.Close()
}

//.......................................SESSION管理类.......................................

type SessionM struct {
	sessions map[uint32]*session
	num      uint32
	lock     sync.RWMutex
}

func NewSessonM() *SessionM {
	return &SessionM{
		sessions: make(map[uint32]*session),
		num:      0,
	}
}

func (this *SessionM) GetSessionById(id uint32) *session {
	if v, exit := this.sessions[id]; exit {
		return v
	}
	return nil
}

func (this *SessionM) SetSession(id uint32, sess *session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.sessions[id] = sess
}

//关闭连接并删除
func (this *SessionM) DelSessionById(id uint32) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v,exit := this.sessions[id];exit{
		v.Close()
	}
	delete(this.sessions, id)
}

func (this *SessionM) WriteByid(id uint32, msg string) bool {
	if v, exit := this.sessions[id]; exit {
		if err := v.Write(msg); err != nil {
			this.DelSessionById(id)
			return false
		} else {
			return true
		}
	}
	return false
}

func (this *SessionM) WriteToAll(msg string) {
	for i, v := range this.sessions {
		if err := v.Write(msg); err != nil {
			log.Println(err)
			this.DelSessionById(i)
		}
	}
}
