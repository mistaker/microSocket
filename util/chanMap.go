package util

var (
	ADD  interface{} = 1
	DEL  interface{} = 2
	GET  interface{} = 3
)


type safeMap struct {
	Msq     chan *[3] interface{}       //['type','id','value',channle]
	data    map[interface{}]interface{}
	chanl   chan interface{}
}

func NewSafeMap() *safeMap {
	tem := &safeMap{}
	tem.init()
	return tem
}

func (this *safeMap) init() {
	this.Msq   = make(chan *[3]interface{},10)
	this.data  = make(map[interface{}]interface{})
	this.chanl = make(chan interface{},0)
	go this.run()
}

func (this *safeMap) run() {
	for {
		select {
		case msg := <- this.Msq :
			switch msg[0] {
			case ADD :
				this.dataAdd(msg[1],msg[2])
			case DEL :
				this.dataDel(msg[1])
			case GET :
				this.dataGet(msg[1])
			}
		}
	}
}

func (this *safeMap) msqChan (typ,id,val interface{}) *[3]interface{}{
	return &[...]interface{}{typ,id,val}
}

//保存 或者更新元素
func (this *safeMap) dataAdd (id , value interface{}) {
	this.data[id] = value
}

//删除元素
func (this *safeMap) dataDel (id interface{}) {
	delete(this.data,id)
}

//获得元素
func (this *safeMap) dataGet (id interface{}) {
	if val ,exit := this.data[id] ;exit {
		this.chanl <- val
		return
	}
	this.chanl <- nil
}

//----------------------------------------------------对外接口--------------------------------
func (this *safeMap) Add (id ,value interface{}) {
	this.Msq <- this.msqChan(ADD,id,value)
}

func (this *safeMap) Del (id interface{}) {
	this.Msq <- this.msqChan(DEL,id ,nil)
}

func (this *safeMap) Get (id interface{}) interface{} {
	this.Msq <- this.msqChan(GET,id,nil)
	res := <- this.chanl
	return res
}

//获得 长度
func (this *safeMap) GetLength() uint32{
	return uint32(len(this.data))
}

/*
func main() {
	sa := NewSafeMap()

//	sa.Add(1,1)
	sa.Add(2,3)
	fmt.Println(2,sa.Get(2))
	sa.Del(2)
	fmt.Println(2,sa.Get(2))
}

*/