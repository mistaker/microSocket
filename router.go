package microSocket

import (
	"log"
	"microSocket/util"
	"reflect"
)

type module interface {
	Default()
	BeforeRequest(data map[string]string) bool
	AfterRequest(data map[string]string)
}

//--------------------------------------------------全局路由------------------------------------------------------------

type RouterMap struct {
	pools util.SafeMap
}

func NewRouterMap() *RouterMap {
	return &RouterMap{}
}

func (this *RouterMap) Register(name string, modules module) {
	this.pools.Set(name, modules)
}

func (this *RouterMap) Hook(moduleName string, funcName string, values map[string]string) {
	item := this.pools.Get(moduleName)
	if item == nil {
		log.Println("not find module " + moduleName)
		return
	}

	var modules module
	if v, ok := item.(module); ok {
		modules = v
	} else {
		return
	}

	//调用模块的 beforequest 请求
	if modules.BeforeRequest(values) == false {
		return
	}

	//反射
	moduleType := reflect.TypeOf(modules)
	moduleValue := reflect.ValueOf(modules)

	//调用相应的接口
	if funcs, exit := moduleType.MethodByName(funcName); exit {
		moduleValue.Method(funcs.Index).Call([]reflect.Value{reflect.ValueOf(values)})
	} else {
		modules.Default()
	}

	//调用模块的 afterquest 请求
	modules.AfterRequest(values)
}

//判断当前模块是否存在
func (this *RouterMap) ModuleExit(moduleName string) bool {
	if this.pools.Get(moduleName) == nil {
		return false
	}
	return true
}
