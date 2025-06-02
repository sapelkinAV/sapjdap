package client

type ReplyPacket struct {
	Length    uint32 // 4 bytes
	ID        uint32 // 4 bytes
	Flags     uint8  // 1 byte
	ErrorCode uint16 // 2 bytes
	Data      []byte
}
