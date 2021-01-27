package main

import (
	"fmt"
	"log"
	"os"
)

const (
	UnAuth   = 1
	WaitAuth = 2
	OutRoom  = 3
	InRoom   = 4
)

var fileName = "chat-client.log"
var logFile *os.File
var logger *log.Logger
var clientPath = ""
var clientStatus = UnAuth

func Init() bool {
	var err error
	logFile, err = os.Create(fileName)
	if err != nil {
		log.Fatal("获取日志文件失败")
	}
	logger = log.New(logFile, "??? ", log.Ldate|log.Ltime|log.Llongfile)
	if logger == nil {
		log.Fatal("日志记录功能初始化失败")
	}
	if !InitConnection() {
		fmt.Println("连接服务器失败...")
		return false
	}
	go dealFromNet()
	if !auth() {
		fmt.Println("验证身份失败...")
		ReleaseConnection()
		return false
	}
	return true
}

func auth() bool {
	var pack LoginReq
	_, err := os.Stat("./auth")
	var authFile *os.File
	authStr := ""
	if err == nil {
		authFile, err = os.Open("./auth")
		if err == nil {
			bytes := make([]byte, 100)
			n, err := authFile.Read(bytes)
			fmt.Println(bytes[:n])
			if err == nil {
				authStr = string(bytes[:n])
			}
		}
	}

	pack.Auth = &authStr
	fmt.Println(authStr)
	SendProto(&pack, pack.GetId())

	ack, ok := <-LoginChan
	if !ok {
		if authFile != nil {
			authFile.Close()
		}
		return false
	}
	logger.Println("err: ", ack.GetError())
	logger.Println("auth: ", ack.GetAuth())
	fmt.Println(ack.GetAuth())
	if len(ack.GetAuth()) != 0 {
		authFile, err = os.Create("./auth")
		authFile.Write([]byte(ack.GetAuth()))
	}
	if authFile != nil {
		authFile.Close()
	}
	return true
}

func Run() {
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	for {
		fmt.Print(clientPath, " > ")
		var cmd, param1, param2, param3, param4 string
		_, _ = fmt.Scanln(&cmd, &param1, &param2, &param3, &param4)
		logger.Print("Read cmd from console: ", cmd)
		if len(cmd) == 0 {
			continue
		}
		handleCmd(cmd, param1, param2, param3, param4)
	}
}

func dealFromNet() {
	for {
		pProto, err := ReadProto()
		if err != nil {
			continue
		}
		logger.Println("Receive proto, id: ", pProto.protoId)
		switch pProto.protoId {
		case uint32(ProtoId_login_resp_id):
			HandleLoginResp(pProto)
			break
		case uint32(ProtoId_create_room_resp_id):
			HandleCreateRoomResp(pProto)
			break
		case uint32(ProtoId_dismiss_room_resp_id):
			HandleDismissRoomResp(pProto)
			break
		case uint32(ProtoId_change_room_settings_resp_id):
			HandleChangeRoomSettingsResp(pProto)
			break
		case uint32(ProtoId_change_room_settings_ntf_id):
			HandleChangeRoomSettingsNtf(pProto)
			break
		case uint32(ProtoId_join_room_resp_id):
			HandleJoinRoomResp(pProto)
			break
		case uint32(ProtoId_change_join_settings_resp_id):
			HandleChangeJoinSettingsResp(pProto)
			break
		case uint32(ProtoId_send_info_resp_id):
			HandleSendInfoResp(pProto)
			break
		case uint32(ProtoId_exit_room_resp_id):
			HandleExitRoomResp(pProto)
			break
		default:
			logger.Println("未知网络消息：", pProto.protoId)
		}
	}
}

func handleCmd(cmd string, param1 string, param2 string, param3 string, param4 string) {
	switch cmd {
	case "cd":
		cd(param1)
		break
	case "ls":
		ls()
		break
	default:
		fmt.Printf("未知命令：\"%s\"\n", cmd)
	}
}

func cd(path string) {
	if len(path) == 0 || path == ".." {
		clientPath = ""
		fmt.Println("退出房间")
	} else {
		clientPath = path
		fmt.Println("进入房间", path)
	}
}

func ls() {
	fmt.Printf("In ls cmd\n")
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
