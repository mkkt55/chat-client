package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
var curRoomId int32 = 0 // 记录当前路径，如果没有在房间里就为空，不然为房间id
var clientStatus = UnAuth

type void struct{}

var voidHolder void
var setAllRoomIds map[int32]void = make(map[int32]void) // 所有房间

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
		// authFile, err = os.Open("./auth")
		// if err == nil {
		// 	bytes := make([]byte, 100)
		// 	n, err := authFile.Read(bytes)
		// 	fmt.Println(bytes[:n])
		// 	if err == nil {
		// 		authStr = string(bytes[:n])
		// 	}
		// }
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
		fmt.Print(curRoomId, " > ")
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
		case uint32(ProtoId_get_all_room_list_resp_id):
			HandleGetAllRoomListResp(pProto)
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
		case uint32(ProtoId_recv_info_ntf_id):
			HandleRecvInfoNtf(pProto)
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
		if len(param1) == 0 || param1 == ".." {
			cd(0)
		} else {
			nId, err := strconv.Atoi(param1)
			if err != nil {
				fmt.Println("请输入房间Id")
				return
			}
			cd(int32(nId))
		}
		break
	case "ls":
		ls()
		break
	case "mkroom":
		mkroom(param1, param2)
		break
	case "rm":
		nId, err := strconv.Atoi(param1)
		if err != nil {
			fmt.Println("请输入房间Id")
			return
		}
		rm(int32(nId))
		break
	case "send":
		send(param1)
		break
	default:
		fmt.Printf("未知命令：\"%s\"\n", cmd)
	}
}

func cd(targetRoomId int32) {
	getAllRoomIds()
	if targetRoomId == 0 {
		if curRoomId != 0 {
			var req ExitRoomReq
			req.RoomId = &curRoomId
			SendProto(&req, req.GetId())
			ack, ok := <-ExitRoomChan
			if !ok {
				fmt.Println("无法退出房间")
				return
			} else if ack.GetError() != ErrorId_err_none {
				fmt.Println("退出房间失败：", ack.GetError())
				return
			} else {
				fmt.Println("退出房间：", curRoomId)
				curRoomId = 0
			}
		}
	} else {
		_, exist := setAllRoomIds[targetRoomId]
		if !exist {
			fmt.Println("请输入房间Id")
			return
		}
		var req JoinRoomReq
		req.RoomId = &targetRoomId
		fmt.Print("请输入您加入的昵称：")
		var joinName string
		fmt.Scanln(&joinName)
		var settings JoinSettings
		settings.JoinName = &joinName
		req.Settings = &settings
		SendProto(&req, req.GetId())
		ack, ok := <-JoinRoomChan
		if !ok {
			fmt.Println("无法加入房间")
		} else if ack.GetError() != ErrorId_err_none {
			fmt.Println("加入房间失败：", ack.GetError())
		} else {
			fmt.Println("加入房间", targetRoomId)
			curRoomId = targetRoomId
		}
	}
}

func ls() {
	getAllRoomIds()
	fmt.Println("Show all room ids:")
	for k := range setAllRoomIds {
		fmt.Println(k)
	}
}

func mkroom(name string, open string) {
	fmt.Println(name)
	var req CreateRoomReq
	var rs RoomSettings
	rs.RoomName = &name
	req.Settings = &rs
	if len(name) == 0 {
		fmt.Print("请输入您想要设置的房间名：")
		fmt.Scanln(&name)
	}
	if len(open) != 0 {
		open := false
		req.Settings.Open = &open
	}
	SendProto(&req, req.GetId())

	ack, ok := <-CreateRoomChan
	if !ok {
		//
	}
	if ack.GetError() != ErrorId_err_none {
		fmt.Println("创建房间失败：", ack.GetError())
		return
	}
	cd(ack.GetNewRoomId())
}

func rm(roomId int32) {
	var req DismissRoomReq
	req.RoomId = &roomId
	SendProto(&req, req.GetId())

	ack, ok := <-DismissRoomChan
	if !ok {
		//
	}
	if ack.GetError() != ErrorId_err_none {
		fmt.Println("解散房间失败：", ack.GetError())
		return
	}
}

func send(msg string) {
	var req SendInfoReq
	req.Info = &msg
	SendProto(&req, req.GetId())

	ack, ok := <-SendInfoChan
	if !ok {
		//
	}
	if ack.GetError() != ErrorId_err_none {
		fmt.Println("发送失败：", ack.GetError())
		return
	}
}

func getAllRoomIds() {
	var req GetAllRoomListReq
	SendProto(&req, req.GetId())

	ack, ok := <-GetAllRoomListChan
	if !ok {
		//
	}
	roomIds := ack.GetRoomIds()
	setAllRoomIds = make(map[int32]void)
	for i := 0; i < len(roomIds); i++ {
		setAllRoomIds[roomIds[i]] = voidHolder
	}
}

func tryEnterChattingGroup(cmdArr []string) {
	//
}
