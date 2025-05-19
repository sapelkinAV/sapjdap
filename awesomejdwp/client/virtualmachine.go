package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sapelkinav/javadap/utils"
)

/*
Id sizes in bytes
(this object should be initialized on start)
*/
type ObjectIdSizes struct {
	FieldIdSize     uint32
	MethodIdSize    uint32
	ObjectIdSize    uint32
	ReferenceIdSize uint32
	FrameIdSize     uint32
}

func (jc *JdwpClient) getVirtualMachineIdSizes() (*ObjectIdSizes, error) {
	reply, err := jc.SendCommand(&cmdVirtualMachineIDSizes, struct{}{})
	if err != nil {
		utils.LogError(log, err, "Could not retrieve virtual machine id sizes")
		return nil, err
	}
	replyPacket := <-reply

	// Check for error code in the reply
	if replyPacket.ErrorCode != 0 {
		err := fmt.Errorf("JDWP command returned error code: %d", replyPacket.ErrorCode)
		log.Error().Err(err).Uint16("errorCode", replyPacket.ErrorCode).Msg("Command failed")
		return nil, err
	}

	if len(replyPacket.Data) < 20 {
		err := fmt.Errorf("received insufficient data for ID sizes: got %d bytes, expected at least 20", len(replyPacket.Data))
		log.Error().Err(err).Msg("Invalid response data for ID sizes")
		return nil, err
	}

	objectIds := &ObjectIdSizes{}

	buf := bytes.NewReader(replyPacket.Data)

	// Read each field
	if err := binary.Read(buf, binary.BigEndian, &objectIds.FieldIdSize); err != nil {
		return nil, fmt.Errorf("failed to read field ID size: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &objectIds.MethodIdSize); err != nil {
		return nil, fmt.Errorf("failed to read method ID size: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &objectIds.ObjectIdSize); err != nil {
		return nil, fmt.Errorf("failed to read object ID size: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &objectIds.ReferenceIdSize); err != nil {
		return nil, fmt.Errorf("failed to read reference ID size: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &objectIds.FrameIdSize); err != nil {
		return nil, fmt.Errorf("failed to read frame ID size: %w", err)
	}

	// Log the parsed values
	log.Debug().
		Uint32("fieldIdSize", objectIds.FieldIdSize).
		Uint32("methodIdSize", objectIds.MethodIdSize).
		Uint32("objectIdSize", objectIds.ObjectIdSize).
		Uint32("referenceIdSize", objectIds.ReferenceIdSize).
		Uint32("frameIdSize", objectIds.FrameIdSize).
		Msg("Parsed object ID sizes")

	return objectIds, nil

}
