package client

import (
	"encoding/binary"
	"fmt"
	"sapelkinav/javadap/utils"
)

const REPLY_HEADER_LENGTH = 11
const REPLY_FLAG = 0x80

func (jc *JdwpClient) initializeReceiveLoop() {
	go func() {
		for {
			select {
			case <-jc.ctx.Done():
				return
			default:
				reply, _, err := jc.readMessage()
				if err != nil {
					utils.LogError(log, err, "Error occured during reading message, ")
					continue
				}
				go prepareReplyChannel(jc, reply)
			}
		}
	}()
}

func prepareReplyChannel(jc *JdwpClient, reply *ReplyPacket) {
	ch, err := jc.jdwpReplies.GetReceiverChannel(reply.ID)
	if err != nil {
		close(ch)
		utils.LogError(log, err, "Error during preparing channel to recieve message")
		jc.jdwpReplies.PopReceiverChannel(reply.ID)
	}
	ch <- reply
	jc.jdwpReplies.PopReceiverChannel(reply.ID)
}

func (jc *JdwpClient) readMessage() (*ReplyPacket, *CommandPacket, error) {
	head := make([]byte, REPLY_HEADER_LENGTH)
	_, err := jc.conn.Read(head)
	if err != nil {
		return nil, nil, err
	}

	flag := head[8]
	if flag == REPLY_FLAG {

		reply := &ReplyPacket{}
		reply.Length = binary.BigEndian.Uint32(head[0:4])
		reply.ID = binary.BigEndian.Uint32(head[4:8])
		reply.Flags = flag
		reply.ErrorCode = binary.BigEndian.Uint16(head[9:11])

		dataLen := int(reply.Length) - REPLY_HEADER_LENGTH
		if dataLen > 0 {
			reply.Data = make([]byte, dataLen)
			_, err := jc.conn.Read(reply.Data)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read reply data: %w", err)
			}
		}
		return reply, nil, nil

	} else {
		return nil, nil, fmt.Errorf("parsing messages that are not replies not implemented yet")
	}

}
