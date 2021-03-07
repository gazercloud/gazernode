package history

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"io/ioutil"
	"time"
)

type FileCache struct {
	filePath   string
	items      []*common_interfaces.ItemValue
	lastReadDT time.Time
}

func NewFileCache(filePath string) *FileCache {
	var c FileCache
	c.filePath = filePath
	c.readFile()
	return &c
}

func (c *FileCache) Read(begin int64, end int64) []*common_interfaces.ItemValue {
	result := make([]*common_interfaces.ItemValue, 0)
	for _, item := range c.items {
		if item.DT >= begin && item.DT < end {
			result = append(result, item)
		}
	}
	c.lastReadDT = time.Now().UTC()
	return result
}

func (c *FileCache) Write(item *common_interfaces.ItemValue) {
	c.items = append(c.items, item)
}

func (c *FileCache) readFile() {
	c.items = make([]*common_interfaces.ItemValue, 0)
	bs, err := ioutil.ReadFile(c.filePath)
	logger.Println("read file ", c.filePath)
	if err == nil {
		currentLine := make([]byte, 0)
		for _, b := range bs {
			if b == 10 || b == 13 {
				if len(currentLine) > 0 {
					var item common_interfaces.ItemValue
					err = json.Unmarshal(currentLine, &item)
					if err == nil {
						c.items = append(c.items, &item)
					} else {
						logger.Println("read error", err)
					}
					currentLine = make([]byte, 0)
				}
			} else {
				currentLine = append(currentLine, b)
			}
		}
	}
}
