package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	UnAuth   = 1
	WaitAuth = 2
	OutRoom  = 3
	InRoom   = 4
)

var page = 0
var path = ""
var status = UnAuth

func Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	for {
		//
		fmt.Printf("%s >", getPageName())
		cmd, _ := reader.ReadString('\n')
		handleCmd(cmd)
	}

}

func Auth() bool {
	time.Sleep(time.Second * 3)
	var pack LoginReq
	str := "111"
	pack.Auth = &str
	SendProto(&pack, pack.GetId())
	str = "222"
	pack.Auth = &str
	SendProto(&pack, pack.GetId())

	time.Sleep(time.Second * 3)
	return true
}

func getPageName() string {
	switch page {
	case 0:
		return "会话"
	case 1:
		return "联系人"
	case 2:
		return "设置"
	}
	return "?"
}

func getPath() string {
	return getPageName() + path
}

func handleCmd(cmd string) {
	var arr = strings.Split(cmd, " ")
	if len(arr) == 0 {
		fmt.Print(cmd)
		return
	}
	switch arr[0] {
	case "cd":
		cd(arr)
		break
	case "ls":
		ls(arr)
		break
	case "findu":
		findu(arr)
		break
	case "findg":
		findg(arr)
		break
	case "addu":
		findg(arr)
		break
	case "addg":
		findg(arr)
		break
	default:
		fmt.Printf("未知命令：%s", arr[0])
	}
}

func cd(cmdArr []string) {
	if len(cmdArr) < 2 {
		fmt.Print(`
		usage: cd [target_num]
		`)
	}
}

func ls(cmdArr []string) {
	if len(cmdArr) < 1 {
		fmt.Print(`
		usage: ls [target_num]
		`)
	}
}

func findu(cmdArr []string) {
	if len(cmdArr) < 2 {
		fmt.Print(`
		usage: findu [user_num]
		`)
	}
}

func findg(cmdArr []string) {
	if len(cmdArr) < 2 {
		fmt.Print(`
		usage: findg [group_num]
		`)
	}
}

func addu(cmdArr []string) {
	if len(cmdArr) < 2 {
		fmt.Print(`
		usage: addu [user_num]
		`)
	}
}

func addg(cmdArr []string) {
	if len(cmdArr) < 2 {
		fmt.Print(`
		usage: addg [group_num]
		`)
	}
}

func tryEnterChattingGroup(cmdArr []string) {
	//
}
