package bin_client

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type BinFrame struct {
	Channel  string
	Password string
	Data     []byte
}

func (c *BinFrame) Marshal() ([]byte, error) {
	lenChannel := len(c.Channel)
	lenPassword := len(c.Password)
	lenData := len(c.Data)

	if lenChannel > 255 {
		return nil, errors.New("wrong channelId length")
	}

	if lenPassword > 255 {
		return nil, errors.New("wrong password length")
	}

	if lenData > 100*1024 {
		return nil, errors.New("wrong data length")
	}

	// FrameSize(4) + lenChannel(4) + lenPassword(4)+ lenData(4) + lenChannel + lenPassword + lenData
	// min: 4 + 4 + 4 + 4 = 16 bytes
	frameSize := 16 + lenChannel + lenPassword + lenData

	result := make([]byte, frameSize)
	offset := 0

	// FrameSize
	binary.LittleEndian.PutUint32(result[offset:], uint32(frameSize))
	offset += 4

	// Lengths
	binary.LittleEndian.PutUint32(result[offset:], uint32(lenChannel))
	offset += 4
	binary.LittleEndian.PutUint32(result[offset:], uint32(lenPassword))
	offset += 4
	binary.LittleEndian.PutUint32(result[offset:], uint32(lenData))
	offset += 4

	// Channel
	copy(result[offset:], c.Channel)
	offset += lenChannel
	// Password
	copy(result[offset:], c.Password)
	offset += lenPassword
	// Data
	copy(result[offset:], c.Data)
	offset += lenData

	return result, nil
}

func UnmarshalBinFrame(frameBytes []byte) (*BinFrame, error) {
	var c BinFrame
	actualFrameSize := len(frameBytes)
	if actualFrameSize < 16 {
		return nil, errors.New("wrong frame actual size")
	}
	offset := 0
	frameSize := binary.LittleEndian.Uint32(frameBytes[offset:])
	if int(frameSize) != actualFrameSize {
		return nil, errors.New("wrong frame actual size (!=) " + fmt.Sprint(frameSize) + " != " + fmt.Sprint(actualFrameSize))
	}
	lenChannel := binary.LittleEndian.Uint32(frameBytes[offset+4:])
	lenPassword := binary.LittleEndian.Uint32(frameBytes[offset+8:])
	lenData := binary.LittleEndian.Uint32(frameBytes[offset+12:])
	offset += 16

	if 16+int(lenChannel)+int(lenPassword)+int(lenData) > int(frameSize) {
		return nil, errors.New("not enough data")
	}

	c.Channel = string(frameBytes[offset : offset+int(lenChannel)])
	offset += int(lenChannel)
	c.Password = string(frameBytes[offset : offset+int(lenPassword)])
	offset += int(lenPassword)
	c.Data = make([]byte, lenData)
	copy(c.Data, frameBytes[offset:offset+int(lenData)])
	offset += int(lenPassword)

	return &c, nil
}
