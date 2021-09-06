package message

import "fmt"

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

// dingSend 信息通过钉钉发送
func (m Msg) DingSend() error {
	return nil
}

// dingSend 信息通过钉钉发送
func (m Msg) MailSend() error {
	return nil
}
