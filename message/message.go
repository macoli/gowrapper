package message

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blinkbean/dingtalk"
)

type Msg struct {
	Code    int64
	Title   string
	Content string
	Time    string
	Err     error
}

// Notice 组合通知
func (m Msg) Notice() {
	//m.Print()
	//m.DingSend()
	//m.MailSend()
}

// Print
func (m Msg) Print() {
	fmt.Println("=====消息通知=====")
	fmt.Printf("时间: %v\n", m.Time)
	fmt.Printf("标题: %v\n", m.Title)
	if m.Err != nil {
		fmt.Printf("错误详情: %v\n", m.Err)
	} else {
		fmt.Printf("内容: %v\n", m.Content)
	}
}

type Auth struct {
	SignSecret  string
	AccessToken string
}

// dingSend 信息通过钉钉发送
func (m Msg) DingSend() error {
	// 打开文件
	file, _ := os.Open("conf.json")
	// 关闭文件
	defer file.Close()
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)

	conf := Auth{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(&conf)

	robot := dingtalk.InitDingTalkWithSecret(conf.AccessToken, conf.SignSecret)

	text := ""
	if m.Err != nil {
		text = fmt.Sprintf(`- 时间: %s
	- 标题: %s
	- 错误详情: %v`, m.Time, m.Title, m.Err)
	} else {
		text = fmt.Sprintf(`- 时间: %s
	- 标题: %s
	- 内容: %s`, m.Time, m.Title, m.Content)
	}

	err = robot.SendMarkDownMessage(m.Title, text)
	if err != nil {
		return err
	}
	return nil
}

// dingSend 信息通过钉钉发送
func (m Msg) MailSend() error {
	return nil
}
