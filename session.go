package microSocket

import (
	"net"
	"sync"
	"time"
)

//---------------------------------------------一个session代表一个连接--------------------------------------------------
type Session struct {
	Id    uint32
	Con   net.Conn
	times int64
	lock  sync.Mutex
}

func NewSession(id uint32, con net.Conn) *Session {
	return &Session{
		Id :    id,
		Con:   con,
		times: time.Now().Unix(),
	}
}

func (this *Session) Write(msg string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	_ ,errs := this.Con.Write([]byte(msg))
	return errs
}

func (this *Session)Close(){
	this.Con.Close()
}

func (this *Session)UpdateTime(){
	this.times = time.Now().Unix()
}
//---------------------------------------------------SESSION管理类------------------------------------------------------

type SessionM struct {
	sessions map[uint32]*Session
	num      uint32
	lock     sync.RWMutex
	isWebSocket bool
	ser     *Msf
}

func NewSessonM(msf *Msf) *SessionM {
	return &SessionM{
		sessions: make(map[uint32]*Session),
		num:      0,
		ser : msf,
	}
}

func (this *SessionM) GetSessionById(id uint32) *Session {
	if v, exit := this.sessions[id]; exit {
		return v
	}
	return nil
}

func (this *SessionM) SetSession(fd uint32, conn net.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	sess := NewSession(fd, conn)
	this.sessions[fd] = sess
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

//向所有客户端发送消息
func (this *SessionM) WriteToAll(msg []byte) {
	for i,_ := range this.sessions {
		this.WriteByid(i,msg)
	}
}

//向单个客户端发送信息
func (this *SessionM) WriteByid(id uint32, msg []byte) bool {
	//把消息打包
	msg = this.ser.SocketType.Pack(msg)

	if v, exit := this.sessions[id]; exit {
		if err := v.Write(string(msg)); err != nil {
			this.DelSessionById(id)
			return false
		} else {
			return true
		}
	}
	return false
}

//心跳检测   每秒遍历一次 查看所有sess 上次接收消息时间  如果超过 num 就删除该 sess
func (this *SessionM)HeartBeat(num int64){
	for {
		time.Sleep(time.Second)
		for i,v:= range this.sessions{
			if time.Now().Unix() - v.times > num {
				this.DelSessionById(i)
			}
		}
	}
}
