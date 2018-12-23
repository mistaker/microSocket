package microSocket

import (
	"net"
	"reflect"
)

const (
	DEFAULTACTION = "Default"
	BEFORACTION   = "BeforeRequest"
	AFTERACTION   = "AfterRequest"
)

type module interface {
	Default(fd uint32,data map[string]string) bool
	BeforeRequest(fd uint32,data map[string]string) bool
	AfterRequest(fd uint32,data map[string]string) bool
}

type eventer interface {
	OnHandel(fd uint32, conn net.Conn) bool
	OnClose(fd uint32)
	OnMessage(fd uint32, msg map[string]string) bool
}

type RoutersMap struct {
	pools    map[string] func(uint32,map[string]string) bool
	strPools map[string] map[string] func(uint32,map[string]string) bool
	structs  map[string] module
	events   eventer
}

func NewRoutersMap() *RoutersMap{
	return &RoutersMap{
		pools :make(map[string]func(uint32,map[string]string) bool),
		strPools:make(map[string]map[string] func(uint32,map[string]string) bool),
		structs : make(map[string]module),
	}
}

//注册事件
func (this *RoutersMap) RegisterEvent (events eventer) {
	this.events = events
}

//注册单个逻辑
func (this *RoutersMap) RegisterFun (methodName string ,funcs func(uint32 ,map[string]string) bool) bool {
	if _, exit := this.pools[methodName]; !exit {
		this.pools[methodName] = funcs
		return true
	}
	return false
}

//结构体 注册
func (this *RoutersMap) RegisterStructFun (moduleName string,mod module) bool {
	if _, exit := this.strPools[moduleName]; exit {
		return false
	}
	this.strPools[moduleName] = make(map[string] func(uint32,map[string]string) bool)
	this.structs[moduleName] = mod

	temType  := reflect.TypeOf(mod)
	temValue := reflect.ValueOf(mod)
	for i := 0 ; i < temType.NumMethod(); i++ {
		tem := temValue.Method(i).Interface()
		if temFunc ,ok := tem.(func(uint32, map[string]string) bool); ok {
			this.strPools[moduleName][temType.Method(i).Name] = temFunc
		}
	}
	return true
}

func (this *RoutersMap) HookAction (funcionName string,fd uint32, data map[string]string) bool{
	if action ,exit := this.pools[funcionName]; exit {
		return action(fd, data)
	} else {
		return false
	}
}

func (this *RoutersMap) HookModule(mouleName string, method string,fd uint32, data map[string]string) bool {
	if _, exit := this.strPools[mouleName]; !exit {
		return false
	}

	if this.strPools[mouleName][BEFORACTION](fd, data) == false {
		return false
	}
	if action, exit := this.strPools[mouleName][method]; exit {
		if action(fd, data) == false {
			return false
		}
	} else {
		if this.strPools[mouleName][DEFAULTACTION](fd, data) == false {
			return false
		}
	}
	if this.strPools[mouleName][AFTERACTION](fd, data) == false {
		return false
	}
	return true
}

func (this *RoutersMap) OnClose(fd uint32) {
	if this.events != nil {
		this.events.OnClose(fd)
	}
}

func (this *RoutersMap) OnHandel(fd uint32, conn net.Conn) bool {
	if this.events != nil {
		return this.events.OnHandel(fd, conn)
	}
	return true
}

func (this *RoutersMap) OnMessage(fd uint32, msg map[string]string) bool {
	if this.events != nil {
		return this.events.OnMessage(fd ,  msg)
	}
	return true
}

