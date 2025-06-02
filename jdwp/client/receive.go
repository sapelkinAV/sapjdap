package client

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"sapelkinav/javadap/jdwp/data/endian"
	"sapelkinav/javadap/jdwp/event/task"
)

const REPLY_HEADER_LENGTH = 11
const REPLY_FLAG = 0x80

// recv decodes all the incoming reply or command packets, forwarding them on
// to the corresponding chans. recv is blocking and should be run on a new
// go routine.
// recv returns when ctx is stopped or there's an IO error.
func (c *Connection) recv(ctx context.Context) {
	for !task.Stopped(ctx) {
		packet, err := c.readPacket()
		switch err {
		case nil:
		case io.EOF:
			return
		default:
			if !task.Stopped(ctx) {
				log.Warn().Err(err).Msg("Failed to read packet")
			}
			return
		}

		switch packet := packet.(type) {
		case replyPacket:
			c.Lock()
			out, ok := c.replies[packet.id]
			delete(c.replies, packet.id)
			c.Unlock()
			if !ok {
				log.Warn().Err(err).Uint32("Unexpected reply for packet", uint32(packet.id))
				continue
			}
			out <- packet

		case cmdPacket:
			switch {
			case packet.cmdSet == cmdSetEvent && packet.cmdID == cmdCompositeEvent:
				d := endian.Reader(bytes.NewReader(packet.data), endian.BigEndian)
				l := events{}
				if err := c.decode(d, reflect.ValueOf(&l)); err != nil {
					log.Warn().Err(err).Msg("Couldn't decode composite event data. Error: ")
					continue
				}

				for _, ev := range l.Events {
					dbg("<%v> event: %T %+v", ev.request(), ev, ev)

					c.Lock()
					handler, ok := c.events[ev.request()]
					c.Unlock()

					if ok {
						handler <- ev
					} else {
						dbg("No event handler registered for %+v", ev)
					}
				}

			default:
				dbg("received unknown packet %+v", packet)
				// Unknown packet. Ignore.
			}
		}
	}
}
