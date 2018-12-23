package microSocket

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"io"
	"log"
	"microSocket/util"
	"strings"
)

type WebSocket struct {
}

//ws接收消息
func (this *WebSocket) ConnHandle(msf *Msf, sess *Session) {
	defer func() {
		msf.SessionMaster.DelSessionById(sess.Id)
		//调用断开链接事件
		msf.EventPool.OnClose(sess.Id)
	}()

	if this.Handshake(sess) == false {
		return
	}

	var (
		buf     []byte
		err     error
		fin     byte
		opcode  byte
		mask    byte
		mKey    []byte
		length  uint64
		l       uint16
		payload byte
		tembuf  []byte
	)

	for {
		buf = make([]byte, 2)
		_, err = io.ReadFull(sess.Con, buf)
		if err != nil {
			break
		}
		fin = buf[0] >> 7
		opcode = buf[0] & 0xf
		if opcode == 8 {
			break
		}
		mask = buf[1] >> 7
		payload = buf[1] & 0x7f
		switch {
		case payload < 126:
			length = uint64(payload)
		case payload == 126:
			buf = make([]byte, 2)
			io.ReadFull(sess.Con, buf)
			binary.Read(bytes.NewReader(buf), binary.BigEndian, &l)
			length = uint64(l)
		case payload == 127:
			buf = make([]byte, 8)
			io.ReadFull(sess.Con, buf)
			binary.Read(bytes.NewReader(buf), binary.BigEndian, &length)
		}
		if mask == 1 {
			mKey = make([]byte, 4)
			io.ReadFull(sess.Con, mKey)
		}
		buf = make([]byte, length)
		io.ReadFull(sess.Con, buf)
		if mask == 1 {
			for i, v := range buf {
				buf[i] = v ^ mKey[i%4]
			}
		}
		//更新最近接收到消息的时间
		sess.updateTime()
		if len(buf) == 0 {
			continue
		}
		tembuf = append(tembuf,buf...)
		if fin == 0 {
			continue
		}
		//把请求的到数据转化为map
		requestData := util.String2Map(string(tembuf))
		tembuf = make([]byte,0)
		//路由调用
		if msf.Hook(sess.Id,requestData) == false {
			return
		}
	}
}

//websocket 打包事件
func (this *WebSocket) Pack(data []byte) []byte {
	length := len(data)
	frame := []byte{129}
	switch {
	case length < 126:
		frame = append(frame, byte(length))
	case length <= 0xffff:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(length))
		frame = append(frame, byte(126))
		frame = append(frame, buf...)
	case uint64(length) <= 0xffffffffffffffff:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(length))
		frame = append(frame, byte(127))
		frame = append(frame, buf...)
	default:
		return []byte{}
	}
	frame = append(frame, data...)
	return frame
}

//握手包
func (this *WebSocket) Handshake(sess *Session) bool {
	reader := bufio.NewReader(sess.Con)
	key := ""
	str := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
			return false
		}
		if len(line) == 0 {
			break
		}
		str = string(line)
		if strings.HasPrefix(str, "Sec-WebSocket-Key") {
			key = str[19:43]
		}
	}
	sha := sha1.New()
	io.WriteString(sha, key+"258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	key = base64.StdEncoding.EncodeToString(sha.Sum(nil))

	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"Sec-WebSocket-Accept: " + key + "\r\n" +
		"Upgrade: websocket\r\n\r\n"
	sess.Con.Write([]byte(header))
	return true
}
