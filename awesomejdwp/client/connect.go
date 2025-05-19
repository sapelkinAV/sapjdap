package client

import (
	"fmt"
	"net"
	"sapelkinav/javadap/utils"
	"time"
)

func (jc *JdwpClient) Connect() error {
	log.Info().Str("addr", jc.addr).Msg("Connecting to JDWP")

	if jc.conn != nil {
		log.Info().Msg("Already connected to JDWP")
		return nil
	}

	// Try to connect with retries
	maxRetries := 5
	backoff := 500 * time.Minute

	var err error
	var conn net.Conn

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Debug().Int("attempt", attempt).Msg("Attempting to connect to JDWP")

		// Set connect timeout
		conn, err = net.DialTimeout("tcp", jc.addr, 3*time.Second)
		if err == nil {
			break
		}

		log.Warn().Err(err).Int("attempt", attempt).Msg("Connection attempt failed")

		if attempt < maxRetries {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
	}

	if err != nil {
		return utils.LogError(log, err, "Failed to connect to JDWP after all attempts")
	}

	// Set read/write deadlines to prevent hanging
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		conn.Close()
		return utils.LogError(log, err, "Failed to set connection deadline")
	}

	jc.conn = conn
	log.Info().Msg("Successfully connected to JDWP")

	// Perform handshake
	if err := jc.handshake(); err != nil {
		jc.Close()
		return err
	}

	// Reset deadline after successful handshake
	if err := jc.conn.SetDeadline(time.Time{}); err != nil {
		jc.Close()
		return utils.LogError(log, err, "Failed to reset connection deadline")
	}

	jc.initializeReceiveLoop()

	// Get VM ID sizes
	objectIdSizes, err := jc.getVirtualMachineIdSizes()
	if err != nil {
		jc.Close()
		return err
	}

	jc.objectIdSizes = objectIdSizes
	log.Info().Msg("JDWP connection fully established")

	return nil
}

func (jc *JdwpClient) Close() {
	if jc.conn != nil {
		log.Info().Msg("Closing JDWP connection")
		err := jc.conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing JDWP connection")
		}
		jc.conn = nil
	}
}
func (jc *JdwpClient) handshake() error {
	log.Debug().Msg("Starting JDWP handshake")

	handshakeBytes := []byte("JDWP-Handshake")

	// Write handshake
	n, err := jc.conn.Write(handshakeBytes)
	if err != nil {
		return utils.LogError(log, err, "Failed to send handshake")
	}

	if n != len(handshakeBytes) {
		err := fmt.Errorf("incomplete handshake send: sent %d of %d bytes", n, len(handshakeBytes))
		return utils.LogError(log, err, "Handshake error")
	}

	log.Debug().Msg("Handshake sent, waiting for response")

	// Read response
	resp := make([]byte, len(handshakeBytes))
	n, err = jc.conn.Read(resp)
	if err != nil {
		return utils.LogError(log, err, "Failed to read handshake response")
	}

	if n != len(handshakeBytes) {
		err := fmt.Errorf("incomplete handshake response: got %d of %d bytes", n, len(handshakeBytes))
		return utils.LogError(log, err, "Handshake error")
	}

	// Verify response
	if string(resp) != string(handshakeBytes) {
		err := fmt.Errorf("invalid handshake response: expected %q, got %q", handshakeBytes, resp)
		return utils.LogError(log, err, "Handshake validation error")
	}

	log.Info().Msg("JDWP handshake successful")
	return nil
}
