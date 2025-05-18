package client

import (
	"context"
	"fmt"
	"net"
)

type JdwpClient struct {
	addr          string
	ctx           context.Context
	conn          net.Conn
	pktID         uint32
	ObjectIdSizes ObjectIdSizes
	JdwpContext   *JdwpContext
}

func NewJdwpClient(addr string) *JdwpClient {
	jc := &JdwpClient{
		addr:        addr,
		ctx:         context.Background(),
		JdwpContext: NewJdwpContext(),
	}
	return jc
}

func (jc *JdwpClient) SendCommand(
	command *Command,
	request interface{},
) (chan *ReplyPacket, error) {

	pkt, _ := NewCommandPacket(
		command,
		request,
	)

	_, err := jc.conn.Write(pkt.PrepareByteBuffer().Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}
	recieverChan := jc.JdwpContext.AddReceiverChannel(pkt.ID)

	return recieverChan, nil
}

func (jc *JdwpClient) addReceiverChannel(messageId uint32) {

}
