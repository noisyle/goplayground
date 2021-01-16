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

// SendCmd 发送指令
func SendCmd(cmd *Cmd) {
	fmt.Println(cmd.Op)
}

// Init 初始化
func Init() {
	loadActions()
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
		fmt.Println("[ActionList]")
		for _, action := range list.Actions {
			_printAction(&action, 0)
		}
	}
}

func _printAction(action *action, level int) {
	prefix := "                    "[:level*2]
	fmt.Println(fmt.Sprintf("%s- [Action]        Name=%s, Type=%s, Class=%s, BorderType=%s", prefix, action.Name, action.Type, action.Class, action.BorderType))
	for _, anime := range action.Animations {
		fmt.Println(fmt.Sprintf("%s  - [Animation]", prefix))
		for _, pose := range anime.Poses {
			fmt.Println(fmt.Sprintf("%s    - [Pose]      Image=%s, ImageAnchor=%s, Velocity=%s, Duration=%s, Sound=%s, Volume=%d", prefix, pose.Image, pose.ImageAnchor, pose.Velocity, pose.Duration, pose.Sound, pose.Volume))
		}
	}
	for _, ref := range action.ActionRefs {
		fmt.Println(fmt.Sprintf("%s  - [ActionRef]   Name=%s, Duration=%s", prefix, ref.Name, ref.Duration))
	}
	for _, act := range action.Actions {
		_printAction(&act, level+1)
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
	XMLName   xml.Name `xml:"ActionReference"`
	Name      string   `xml:"Name,attr"`
	Duration  string   `xml:"Duration,attr"`
	InitialVX string   `xml:"InitialVX,attr"`
	InitialVY string   `xml:"InitialVY,attr"`
	X         string   `xml:"X,attr"`
	Y         string   `xml:"Y,attr"`
	TargetX   string   `xml:"TargetX,attr"`
	TargetY   string   `xml:"TargetY,attr"`
	LookRight string   `xml:"LookRight,attr"`
	Condition string   `xml:"Condition,attr"`
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
	Duration    string   `xml:"Duration,attr"`
	Sound       string   `xml:"Sound,attr"`
	Volume      int      `xml:"Volume,attr"`
}
