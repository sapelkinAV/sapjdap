package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
length (4 bytes)
id (4 bytes)
flags (1 byte)
command set (1 byte)
command (1 byte)
In total 11 bytes
*/
const JDWP_PACKET_HEADER_SIZE = 11

type CommandPacket struct {
	// header
	ID     uint32
	Length uint32
	Flags  uint8
	CmdSet uint8
	Cmd    uint8

	//value
	Data []byte
}

func NewCommandPacket(command *Command, request interface{}) (*CommandPacket, error) {
	data := new(bytes.Buffer)
	if request != nil {
		err := binary.Write(data, binary.BigEndian, request)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request: %w", err)
		}
	}

	packet := CommandPacket{
		ID:     NextCommandId(),
		Flags:  0x00,
		CmdSet: uint8(command.Set),
		Cmd:    uint8(command.Id),
		Data:   data.Bytes(),
	}
	packet.Length = packet.GetLength()
	return &packet, nil
}

func (pkt *CommandPacket) PrepareByteBuffer() *bytes.Buffer {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, pkt.GetLength())
	binary.Write(buf, binary.BigEndian, pkt.ID)
	buf.WriteByte(pkt.Flags)
	buf.WriteByte(pkt.CmdSet)
	buf.WriteByte(pkt.Cmd)
	buf.Write(pkt.Data)

	return buf
}

func (pkt *CommandPacket) GetLength() uint32 {

	return JDWP_PACKET_HEADER_SIZE + uint32(len(pkt.Data))
}
