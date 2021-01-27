package main

import (
	"google.golang.org/protobuf/proto"
)

var LoginChan = make(chan LoginResp, 1)

func HandleLoginResp(pProto *ProtoPack) bool {
	var ack LoginResp
	err := proto.Unmarshal(pProto.body, &ack)
	if err != nil {
		logger.Println("Unmarshal proto fail...")
		return false
	}
	LoginChan <- ack
	return true
}

func HandleCreateRoomResp(pProto *ProtoPack) bool {
	return true
}

func HandleDismissRoomResp(pProto *ProtoPack) bool { return true }

func HandleChangeRoomSettingsResp(pProto *ProtoPack) bool { return true }

func HandleChangeRoomSettingsNtf(pProto *ProtoPack) bool { return true }

func HandleJoinRoomResp(pProto *ProtoPack) bool { return true }

func HandleChangeJoinSettingsResp(pProto *ProtoPack) bool { return true }

func HandleSendInfoResp(pProto *ProtoPack) bool { return true }

func HandleExitRoomResp(pProto *ProtoPack) bool { return true }
