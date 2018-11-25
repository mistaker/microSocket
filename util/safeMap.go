package util

import "sync"

type SafeMap struct {
	Data map[string]interface{}
	Lock sync.RWMutex
}

func (this *SafeMap) Get(k string) interface{} {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if v, exit := this.Data[k]; exit {
		return v
	}
	return nil
}

func (this *SafeMap) Set(k string, v interface{}) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.Data == nil {
		this.Data = make(map[string]interface{})
	}
	this.Data[k] = v
}


