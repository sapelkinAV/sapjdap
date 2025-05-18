package client

import (
	"fmt"
	"net"
	"time"
)

func (jc *JdwpClient) Connect() error {
	var err error
	if jc.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", jc.addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	jc.conn = conn
	// add logic of handshake and initialization of vm ids
	jc.handshake()

	return err
}

func (jc *JdwpClient) handshake() error {
	var err error
	var handshake = []byte("JDWP-Handshake")
	_, err = jc.conn.Write(handshake)
	if err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	resp := make([]byte, len(handshake))
	_, err = jc.conn.Read(resp)
	if err != nil {
		return fmt.Errorf("failed to read handshake response: %w", err)
	}

	if string(resp) != string(handshake) {
		return fmt.Errorf("invalid handshake response: got %q", resp)
	}

	fmt.Println("JDWP handshake successful!")
	return err
}
