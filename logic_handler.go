package main

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

var LoginChan = make(chan LoginResp, 1)
var GetAllRoomListChan = make(chan GetAllRoomListResp, 1)
var CreateRoomChan = make(chan CreateRoomResp, 1)
var DismissRoomChan = make(chan DismissRoomResp, 1)
var ChangeRoomSettingsChan = make(chan ChangeRoomSettingsResp, 1)
var JoinRoomChan = make(chan JoinRoomResp, 1)
var GetRoomAllMemberChan = make(chan GetRoomAllMemberResp, 1)
var SendInfoChan = make(chan SendInfoResp, 1)
var ExitRoomChan = make(chan ExitRoomResp, 1)

func HandleLoginResp(pProto *ProtoPack) bool {
	var ack LoginResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	fmt.Println("收到验证身份ACK, auth: ", ack.GetAuth())
	LoginChan <- ack
	return true
}

func HandleGetAllRoomListResp(pProto *ProtoPack) bool {
	var ack GetAllRoomListResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	GetAllRoomListChan <- ack
	return true
}

func HandleGetRoomAllMembersResp(pProto *ProtoPack) bool {
	var ack GetRoomAllMemberResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	GetRoomAllMemberChan <- ack
	return true
}

func HandleCreateRoomResp(pProto *ProtoPack) bool {
	var ack CreateRoomResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	CreateRoomChan <- ack
	return true
}

func HandleDismissRoomResp(pProto *ProtoPack) bool {
	var ack DismissRoomResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	DismissRoomChan <- ack
	return true
}

func HandleChangeRoomSettingsResp(pProto *ProtoPack) bool {
	var ack ChangeRoomSettingsResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	ChangeRoomSettingsChan <- ack
	return true
}

func HandleChangeRoomSettingsNtf(pProto *ProtoPack) bool { return true }

func HandleJoinRoomResp(pProto *ProtoPack) bool {
	var ack JoinRoomResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	JoinRoomChan <- ack
	return true
}

func HandleChangeJoinSettingsResp(pProto *ProtoPack) bool { return true }

func HandleSendInfoResp(pProto *ProtoPack) bool {
	var ack SendInfoResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	SendInfoChan <- ack
	return true
}

func HandleRecvInfoNtf(pProto *ProtoPack) bool {
	var ntf RecvInfoNtf
	err := proto.Unmarshal(pProto.body, &ntf)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ntf.GetId())
		return false
	}
	fmt.Printf("\r来自[%s]的消息：%s\n", ntf.GetSenderName(), ntf.GetMsg())
	fmt.Print("房间", curRoomId, "> ")
	return true
}

func HandleExitRoomResp(pProto *ProtoPack) bool {
	var ack ExitRoomResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...", ack.GetId())
		return false
	}
	ExitRoomChan <- ack
	return true
}
