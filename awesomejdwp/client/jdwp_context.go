package client

import (
	"context"
	"fmt"
	"sync"
)

type JdwpContext struct {
	ctx                           context.Context
	ReplyPacketReceiverMap        map[uint32]chan *ReplyPacket
	ReplyCommandPacketReceiverMap map[uint32]chan *CommandPacket
	muReply                       sync.Mutex
	muCommand                     sync.Mutex
}

func NewJdwpContext() *JdwpContext {
	return &JdwpContext{
		ctx:                           context.Background(),
		ReplyPacketReceiverMap:        make(map[uint32]chan *ReplyPacket),
		ReplyCommandPacketReceiverMap: make(map[uint32]chan *CommandPacket),
	}
}

func (jc *JdwpContext) AddReceiverChannel(messageId uint32) chan *ReplyPacket {
	jc.muReply.Lock()
	defer jc.muReply.Unlock()

	jc.ReplyPacketReceiverMap[messageId] = make(chan *ReplyPacket)
	return jc.ReplyPacketReceiverMap[messageId]
}

func (jc *JdwpContext) GetReceiverChannel(messageId uint32) (chan *ReplyPacket, error) {
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

func (jc *JdwpContext) PopReceiverChannel(messageId uint32) {
	jc.muReply.Lock()
	defer jc.muReply.Unlock()

	if jc.ReplyPacketReceiverMap[messageId] != nil {
		close(jc.ReplyPacketReceiverMap[messageId])
		delete(jc.ReplyPacketReceiverMap, messageId)
	} else {
		fmt.Println("ReplyPacketChannel does not exist for messageId: ", messageId)
	}

}
