package history

import (
	"encoding/binary"
	"github.com/gazercloud/gazernode/common_interfaces"
	"io/ioutil"
	"time"
)

type FileCache struct {
	filePath   string
	items      []common_interfaces.ItemValue
	lastReadDT time.Time
}

func NewFileCache(filePath string) *FileCache {
	var c FileCache
	c.filePath = filePath
	c.readFileBin()
	return &c
}

func (c *FileCache) Read(begin int64, end int64) []*common_interfaces.ItemValue {
	result := make([]*common_interfaces.ItemValue, 0)
	for index, item := range c.items {
		if item.DT >= begin && item.DT < end {
			result = append(result, &c.items[index])
		}
	}
	c.lastReadDT = time.Now().UTC()
	return result
}

func (c *FileCache) Write(item *common_interfaces.ItemValue) {
	c.items = append(c.items, *item)
}

func (c *FileCache) readFileBin() {
	bs, err := ioutil.ReadFile(c.filePath)

	if err == nil {
		offset := 0
		count := 0
		for offset < len(bs) {
			if bs[offset] != 0xAA {
				offset++
				continue
			}
			if offset+5 > len(bs) {
				offset++
				continue
			}
			blockSize := int(bs[offset+1])
			if bs[offset+blockSize-1] != 0x55 {
				offset += 1
				continue
			}

			block := bs[offset : offset+blockSize]

			valueLen := int(block[1+1+8])
			if valueLen+1+1+8+1 > len(block) {
				offset += 1
				continue
			}

			uomLen := int(block[1+1+8+1+valueLen])
			if 1+1+8+1+valueLen+1+uomLen > len(block) {
				offset += 1
				continue
			}

			count++
			offset += blockSize
		}

		c.items = make([]common_interfaces.ItemValue, count)

		offset = 0
		index := 0
		for offset < len(bs) {
			if bs[offset] != 0xAA {
				offset++
				continue
			}

			blockSize := int(bs[offset+1])
			block := bs[offset : offset+blockSize]

			if block[blockSize-1] != 0x55 {
				offset += 1
				continue
			}

			item := &c.items[index]
			item.DT = int64(binary.LittleEndian.Uint64(block[1+1:]))
			valueLen := int(block[1+1+8])
			if valueLen+1+1+8+1 > len(block) {
				offset += 1
				continue
			}
			item.Value = string(block[1+1+8+1 : 1+1+8+1+valueLen])
			uomLen := int(block[1+1+8+1+valueLen])
			if 1+1+8+1+valueLen+1+uomLen > len(block) {
				offset += 1
				continue
			}
			item.UOM = string(block[1+1+8+1+valueLen+1 : 1+1+8+1+valueLen+1+uomLen])

			offset += blockSize
			index++
		}
	}
}
