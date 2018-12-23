package microSocket

import (
	"log"
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

func (this *Session) write(msg string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	_ ,errs := this.Con.Write([]byte(msg))
	return errs
}

func (this *Session)close(){
	this.Con.Close()
}

func (this *Session)updateTime(){
	this.times = time.Now().Unix()
}
//---------------------------------------------------SESSION管理类------------------------------------------------------

type SessionM struct {
	isWebSocket bool
	ser     *Msf
	sessions sync.Map
}

func NewSessonM(msf *Msf) *SessionM {
	if msf == nil {
		return nil
	}

	return &SessionM{
		ser : msf,
	}
}

func (this *SessionM) GetSessionById(id uint32) *Session {
	tem ,exit := this.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*Session) ; ok {
			return sess
		}
	}
	return nil
}

func (this *SessionM) SetSession(fd uint32, conn net.Conn) {
	sess := NewSession(fd, conn)
	this.sessions.Store(fd,sess)
}

//关闭连接并删除
func (this *SessionM) DelSessionById(id uint32) {
	tem ,exit := this.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*Session) ; ok {
			sess.close()
		}
	}
	this.sessions.Delete(id)
}

//向所有客户端发送消息
func (this *SessionM) WriteToAll(msg []byte) {
	msg = this.ser.SocketType.Pack(msg)
	this.sessions.Range(func(key,val interface{})bool{
		if val, ok := val.(*Session); ok {
			if err := val.write(string(msg)); err != nil {
				this.DelSessionById(key.(uint32))
				log.Println(err)
			}
		}
		return true
	})
}

//向单个客户端发送信息
func (this *SessionM) WriteByid(id uint32, msg []byte) bool {
	//把消息打包
	msg = this.ser.SocketType.Pack(msg)

	tem ,exit := this.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*Session) ; ok {
			if err := sess.write(string(msg)); err == nil {
				return true
			}
		}
	}
	this.DelSessionById(id)
	return false
}

//心跳检测   每秒遍历一次 查看所有sess 上次接收消息时间  如果超过 num 就删除该 sess
func (this *SessionM)HeartBeat(num int64){
	for {
		time.Sleep(time.Second)
		this.sessions.Range(func (key,val interface{})bool {
			tem , ok := val.(*Session)
			if !ok {
				return true
			}

			if time.Now().Unix() - tem.times > num {
				this.DelSessionById(key.(uint32))
			}
			return true
		})

	}
}
