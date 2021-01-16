package daemon

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Shimejid 守护协程
type Shimejid struct {
}

type Mascot struct {
	XMLName     xml.Name     `xml:"Mascot"`
	ActionLists []ActionList `xml:"ActionList"`
}

type ActionList struct {
	XMLName xml.Name `xml:"ActionList"`
	Actions []Action `xml:"Action"`
}

type Action struct {
	XMLName    xml.Name    `xml:"Action"`
	Name       string      `xml:"Name,attr"`
	Type       string      `xml:"Type,attr"`
	BorderType string      `xml:"BorderType,attr"`
	Class      string      `xml:"Class,attr"`
	Loop       bool        `xml:"Loop,attr"`
	Condition  string      `xml:"Condition,attr"`
	Animations []Animation `xml:"Animation"`
	Actions    []Action    `xml:"Action"`
	ActionRefs []ActionRef `xml:"ActionReference"`
}

type ActionRef struct {
	XMLName xml.Name `xml:"ActionReference"`
	Name    string   `xml:"Name,attr"`
}

type Animation struct {
	XMLName xml.Name `xml:"Animation"`
	Poses   []Pose   `xml:"Pose"`
}

type Pose struct {
	XMLName     xml.Name `xml:"Pose"`
	Image       string   `xml:"Image,attr"`
	ImageAnchor string   `xml:"ImageAnchor,attr"`
	Velocity    string   `xml:"Velocity,attr"`
	Duration    int      `xml:"Duration,attr"`
	Sound       string   `xml:"Sound,attr"`
	Volume      int      `xml:"Volume,attr"`
}

// Start 启动守护协程
func Start() *Shimejid {
	d := new(Shimejid)
	d.loadActions()
	go d.run()
	return d
}

// SendMessage 发送文字
func (d *Shimejid) SendMessage(s string) {
	fmt.Println(s)
}

func (d *Shimejid) run() {
	for {
		// TODO
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *Shimejid) loadActions() {
	// 加载配置文件
	file, err := os.Open("conf/actions.xml")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	m := Mascot{}
	err = xml.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	for _, list := range m.ActionLists {
		for _, action := range list.Actions {
			fmt.Println(fmt.Sprintf("[Action]        Name=%s, Type=%s, Class=%s, BorderType=%s", action.Name, action.Type, action.Class, action.BorderType))
			for _, anime := range action.Animations {
				fmt.Println(fmt.Sprintf("- [Animation]"))
				for _, pose := range anime.Poses {
					fmt.Println(fmt.Sprintf("  - [Pose]      Image=%s, ImageAnchor=%s, Velocity=%s, Duration=%d, Sound=%s, Volume=%d", pose.Image, pose.ImageAnchor, pose.Velocity, pose.Duration, pose.Sound, pose.Volume))
				}
			}
			for _, ref := range action.ActionRefs {
				fmt.Println(fmt.Sprintf("  - [ActionRef] Name=%s", ref.Name))
			}
		}
	}
}
