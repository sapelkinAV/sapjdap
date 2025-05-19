package client

import (
	"context"
	"fmt"
	"net"
	"sapelkinav/javadap/utils"
)

var log, _ = utils.GetComponentLogger("jdwp", "client")

type JdwpClient struct {
	addr          string
	ctx           context.Context
	conn          net.Conn
	pktID         uint32
	objectIdSizes *ObjectIdSizes
	jdwpReplies   *JdwpReplies
}

func NewJdwpClient(addr string, ctx context.Context) *JdwpClient {
	jc := &JdwpClient{
		addr:        addr,
		ctx:         ctx,
		jdwpReplies: NewJdwpReplies(),
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
		return nil, utils.LogError(log, err, "failed to send command")
	}
	receiverChan := jc.jdwpReplies.AddReceiverChannel(pkt.ID)

	return receiverChan, nil
}

func (jc *JdwpClient) HelloWorld() {

	fmt.Println("Hello World")

}
