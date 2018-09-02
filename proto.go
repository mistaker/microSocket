package microSocket

import (
	"bytes"
	"encoding/binary"
	"errors"
)

/*
   一个简单的通讯协议，由 header + 信息长度 ＋ 信息内容组成
*/

const (
	CONSTHEADER       = "Header"
	CONSTHEADERLENGTH = 6
	CONSTMLENGTH      = 4
)

func Enpack(message []byte) []byte {
	return append(append([]byte(CONSTHEADER), IntToBytes(len(message))...), message...)
}

func Depack(buff []byte) ([]byte, []byte, error) {
	length := len(buff)

	//如果包长小于header 就直接返回 因为接收的数据不完整
	if length < CONSTHEADERLENGTH+CONSTMLENGTH {
		return buff, nil, nil
	}

	//如果header不是 指定的header 说明此数据已经被污染 直接返回错误
	if string(buff[:CONSTHEADERLENGTH]) != CONSTHEADER {
		return []byte{}, nil, errors.New("header is not safe")
	}

	msgLength := BytesToInt(buff[CONSTHEADERLENGTH : CONSTHEADERLENGTH+CONSTMLENGTH])
	if length < CONSTHEADERLENGTH+CONSTMLENGTH+msgLength {
		return buff, nil, nil
	}

	data := buff[CONSTHEADERLENGTH+CONSTMLENGTH : CONSTHEADERLENGTH+CONSTMLENGTH+msgLength]
	buffs := buff[CONSTHEADERLENGTH+CONSTMLENGTH+msgLength:]
	return buffs, data, nil
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
