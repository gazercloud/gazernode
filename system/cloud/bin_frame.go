package cloud

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

type BinFrameHeader struct {
	TargetNodeId   string `json:"n"`
	Function       string `json:"f"`
	IsRequest      bool   `json:"r"`
	TransactionId  string `json:"t"`
	CloudSessionId string `json:"s"`
	Error          string `json:"e"`
	RequestType    string `json:"rt"`
}

type BinFrame struct {
	Header BinFrameHeader
	Data   []byte
}

func (c *BinFrame) Marshal() ([]byte, error) {
	headerBytes, err := json.Marshal(c.Header)
	if err != nil {
		return nil, err
	}

	frameSize := 4 + 4 + len(headerBytes) + len(c.Data)

	result := make([]byte, frameSize)
	offset := 0

	// FrameSize
	binary.LittleEndian.PutUint32(result[offset:], uint32(frameSize))
	offset += 4
	// Lengths
	binary.LittleEndian.PutUint32(result[offset:], uint32(len(headerBytes)))
	offset += 4
	// Header
	copy(result[offset:], headerBytes)
	offset += len(headerBytes)
	// Body
	copy(result[offset:], c.Data)
	offset += len(c.Data)

	return result, nil
}

func UnmarshalBinFrame(frameBytes []byte) (*BinFrame, error) {
	var c BinFrame
	actualFrameSize := len(frameBytes)
	if actualFrameSize < 8 {
		return nil, errors.New("wrong frame actual size (<8):" + fmt.Sprint(len(frameBytes)))
	}
	offset := 0
	frameSize := binary.LittleEndian.Uint32(frameBytes[offset:])
	if int(frameSize) != actualFrameSize {
		return nil, errors.New("wrong frame actual size (!=) " + fmt.Sprint(frameSize) + " != " + fmt.Sprint(actualFrameSize))
	}
	offset += 4
	lenHeader := int(binary.LittleEndian.Uint32(frameBytes[offset:]))
	dataSize := actualFrameSize - 8 - lenHeader
	if dataSize < 0 {
		return nil, errors.New("wrong data size. FrameSize:" + fmt.Sprint(frameSize) + " HeaderSize:" + fmt.Sprint(lenHeader) + " DataSize:" + fmt.Sprint(dataSize))
	}
	offset += 4

	err := json.Unmarshal(frameBytes[8:8+lenHeader], &c.Header)
	offset += lenHeader
	if err == nil {
		c.Data = make([]byte, len(frameBytes)-offset)
		copy(c.Data, frameBytes[8+lenHeader:])
	}

	return &c, err
}
