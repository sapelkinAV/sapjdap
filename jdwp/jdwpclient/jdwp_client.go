package jdwpclient

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"sapelkinav/javadap/jdwp/data/binary"
	"sapelkinav/javadap/jdwp/data/endian"
	"sapelkinav/javadap/utils"
	"sync"
	"time"
)

var log, _ = utils.GetComponentLogger("jdwp", "client")

var (
	handshake = []byte("JDWP-Handshake")

	defaultIDSizes = IDSizes{
		FieldIDSize:         8,
		MethodIDSize:        8,
		ObjectIDSize:        8,
		ReferenceTypeIDSize: 8,
		FrameIDSize:         8,
	}
)

type Connection struct {
	in           io.Reader
	r            binary.Reader
	w            binary.Writer
	flush        func() error
	idSizes      IDSizes
	nextPacketID packetID
	events       map[EventRequestID]chan<- Event
	replies      map[packetID]chan<- replyPacket
	sync.Mutex
}

// Open creates a Connection using conn for I/O.
func Open(ctx context.Context, conn io.ReadWriteCloser) (*Connection, error) {
	if err := exchangeHandshakes(conn); err != nil {
		return nil, err
	}

	buf := bufio.NewWriterSize(conn, 1024)
	r := endian.Reader(conn, endian.BigEndian)
	w := endian.Writer(buf, endian.BigEndian)
	c := &Connection{
		in:      conn,
		r:       r,
		w:       w,
		flush:   buf.Flush,
		idSizes: defaultIDSizes,
		events:  map[EventRequestID]chan<- Event{},
		replies: map[packetID]chan<- replyPacket{},
	}
	go func() { c.recv(ctx) }()
	var err error
	c.idSizes, err = c.GetIDSizes()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func exchangeHandshakes(conn io.ReadWriter) error {
	if _, err := conn.Write(handshake); err != nil {
		return err
	}
	ok, err := expect(conn, handshake)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Bad handshake")
	}
	return nil
}

// expect reads c.in, expecting the specfified sequence of bytes. If the read
// data doesn't match, then the function returns immediately with false.
func expect(conn io.Reader, expected []byte) (bool, error) {
	got := make([]byte, len(expected))
	for len(expected) > 0 {
		n, err := conn.Read(got)
		if err != nil {
			return false, err
		}
		for i := 0; i < n; i++ {
			if got[i] != expected[i] {
				return false, nil
			}
		}
		got, expected = got[n:], expected[n:]
	}
	return true, nil
}

// get sends the specified command and waits for a reply.
func (c *Connection) get(cmd cmd, req interface{}, out interface{}) error {
	p, err := c.req(cmd, req)
	if err != nil {
		return err
	}
	return p.wait(out)
}

// req sends the specified command and returns a pending.
func (c *Connection) req(cmd cmd, req interface{}) (*pending, error) {
	data := bytes.Buffer{}
	if req != nil {
		e := endian.Writer(&data, endian.BigEndian)
		if err := c.encode(e, reflect.ValueOf(req)); err != nil {
			return nil, err
		}
	}

	id, replyChan := c.newReplyHandler()

	p := cmdPacket{id: id, cmdSet: cmd.set, cmdID: cmd.id, data: data.Bytes()}

	c.Lock()
	defer c.Unlock()

	if err := p.write(c.w); err != nil {
		return nil, err
	}
	if err := c.flush(); err != nil {
		return nil, err
	}

	dbg("<%v> send: %v, %+v", id, cmd, req)

	return &pending{c, replyChan, id}, nil
}

type pending struct {
	c  *Connection
	p  <-chan replyPacket
	id packetID
}

// wait blocks until the penging response is received, filling out with the
// response data.
func (p *pending) wait(out interface{}) error {
	select {
	case reply := <-p.p:
		if reply.err != ErrNone {
			dbg("<%v> recv err: %+v", p.id, reply.err)
			return reply.err
		}
		if out == nil {
			return nil
		}
		r := bytes.NewReader(reply.data)
		d := endian.Reader(r, endian.BigEndian)
		if err := p.c.decode(d, reflect.ValueOf(out)); err != nil {
			return err
		}
		dbg("<%v> recv: %+v", p.id, out)
		if offset, _ := r.Seek(0, 1); offset != int64(len(reply.data)) {
			panic(fmt.Errorf("Only %d/%d bytes read from reply packet", offset, len(reply.data)))
		}
		return nil
	case <-time.After(time.Second * 120):
		return fmt.Errorf("timeout")
	}
}

func (c *Connection) newReplyHandler() (packetID, <-chan replyPacket) {
	reply := make(chan replyPacket, 1)
	c.Lock()
	id := c.nextPacketID
	c.nextPacketID++
	c.replies[id] = reply
	c.Unlock()
	return id, reply
}
