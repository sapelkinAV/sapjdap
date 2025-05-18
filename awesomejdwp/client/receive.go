package client

import (
	"encoding/binary"
	"fmt"
)

const REPLY_HEADER_LENGTH = 11
const REPLY_FLAG = 0x80

func (jc *JdwpClient) receiveLoop() {
	go func() {
		for {
			reply, _, err := jc.readMessage()
			if err != nil {
				fmt.Println("Error occured during reading message, ", err)
				continue
			}

			go func() {
				ch, err := jc.JdwpContext.GetReceiverChannel(reply.ID)
				if err != nil {
					ch <- reply
					close(ch)
					jc.JdwpContext.PopReceiverChannel(reply.ID)
				}
			}()
		}
	}()
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
