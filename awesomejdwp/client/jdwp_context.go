package client

import (
	"fmt"
	"sync"
)

type JdwpReplies struct {
	ReplyPacketReceiverMap        map[uint32]chan *ReplyPacket
	ReplyCommandPacketReceiverMap map[uint32]chan *CommandPacket
	muReply                       sync.Mutex
	muCommand                     sync.Mutex
}

func NewJdwpReplies() *JdwpReplies {
	return &JdwpReplies{
		ReplyPacketReceiverMap:        make(map[uint32]chan *ReplyPacket),
		ReplyCommandPacketReceiverMap: make(map[uint32]chan *CommandPacket),
	}
}

func (jc *JdwpReplies) AddReceiverChannel(messageId uint32) chan *ReplyPacket {
	jc.muReply.Lock()
	defer jc.muReply.Unlock()

	jc.ReplyPacketReceiverMap[messageId] = make(chan *ReplyPacket)
	return jc.ReplyPacketReceiverMap[messageId]
}

func (jc *JdwpReplies) GetReceiverChannel(messageId uint32) (chan *ReplyPacket, error) {
	jc.muReply.Lock()
	defer jc.muReply.Unlock()

	if jc.ReplyPacketReceiverMap == nil {
		return nil, fmt.Errorf(
			"ReplyPacketReceiver channel does not exist for id: %d",
			messageId,
		)
	}

	return jc.ReplyPacketReceiverMap[messageId], nil
}

func (jc *JdwpReplies) PopReceiverChannel(messageId uint32) {
	jc.muReply.Lock()
	defer jc.muReply.Unlock()

	if jc.ReplyPacketReceiverMap[messageId] != nil {
		close(jc.ReplyPacketReceiverMap[messageId])
		delete(jc.ReplyPacketReceiverMap, messageId)
	} else {
		fmt.Println("ReplyPacketChannel does not exist for messageId: ", messageId)
	}

}
