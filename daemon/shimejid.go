package daemon

import (
	"fmt"
	"time"
)

// Shimejid 守护协程
type Shimejid struct {
}

// Init 启动守护协程
func Init() *Shimejid {
	d := new(Shimejid)
	go d.run()
	return d
}

// SendMessage 发送文字
func (d *Shimejid) SendMessage(s string) {
	fmt.Println(s)
}

func (d *Shimejid) run() {
	for {
		fmt.Println("Shimejid running")
		time.Sleep(100 * time.Millisecond)
	}
}
