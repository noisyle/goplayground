package daemon

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// Cmd 指令
type Cmd struct {
	Op string
}

// CmdChan 指令管道
var CmdChan = make(chan *Cmd, 100)

// Start 启动守护线程
func Start() {
	loadActions()
	go func() {
		for {
			select {
			case cmd := <-CmdChan:
				fmt.Println(cmd.Op)
				if cmd.Op == "exit" {
					return
				}
			}
		}
	}()
}

func loadActions() {
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
	m := mascot{}
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

type mascot struct {
	XMLName     xml.Name     `xml:"Mascot"`
	ActionLists []actionList `xml:"ActionList"`
}

type actionList struct {
	XMLName xml.Name `xml:"ActionList"`
	Actions []action `xml:"Action"`
}

type action struct {
	XMLName    xml.Name    `xml:"Action"`
	Name       string      `xml:"Name,attr"`
	Type       string      `xml:"Type,attr"`
	BorderType string      `xml:"BorderType,attr"`
	Class      string      `xml:"Class,attr"`
	Loop       bool        `xml:"Loop,attr"`
	Condition  string      `xml:"Condition,attr"`
	Animations []animation `xml:"Animation"`
	Actions    []action    `xml:"Action"`
	ActionRefs []actionRef `xml:"ActionReference"`
}

type actionRef struct {
	XMLName xml.Name `xml:"ActionReference"`
	Name    string   `xml:"Name,attr"`
}

type animation struct {
	XMLName xml.Name `xml:"Animation"`
	Poses   []pose   `xml:"Pose"`
}

type pose struct {
	XMLName     xml.Name `xml:"Pose"`
	Image       string   `xml:"Image,attr"`
	ImageAnchor string   `xml:"ImageAnchor,attr"`
	Velocity    string   `xml:"Velocity,attr"`
	Duration    int      `xml:"Duration,attr"`
	Sound       string   `xml:"Sound,attr"`
	Volume      int      `xml:"Volume,attr"`
}
