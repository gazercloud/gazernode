package history

import (
	"encoding/binary"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/settings"
	"github.com/gazercloud/gazernode/utilities"
	"github.com/gazercloud/gazernode/utilities/logger"
	"os"
	"strconv"
	"sync"
	"time"
)

type Item struct {
	ss               *settings.Settings
	id               uint64
	historyDepthDays int
	data             []*common_interfaces.ItemValue
	mtx              sync.Mutex
	files            map[string]*FileCache
	flushFinished    bool
}

func NewItem(id uint64, ss *settings.Settings) *Item {
	var c Item
	c.ss = ss
	c.id = id
	c.historyDepthDays = 30
	c.data = make([]*common_interfaces.ItemValue, 0)
	c.files = make(map[string]*FileCache)
	return &c
}

func (c *Item) Write(value common_interfaces.ItemValue) {
	c.mtx.Lock()
	c.data = append(c.data, &value)
	c.mtx.Unlock()
}

func (c *Item) Read(dtBegin int64, dtEnd int64) *ReadResult {
	var result ReadResult
	result.Items = make([]*common_interfaces.ItemValue, 0)
	c.mtx.Lock()
	result.Items = append(result.Items, c.readFiles(dtBegin, dtEnd)...)
	for _, item := range c.data {
		if item.DT >= dtBegin && item.DT < dtEnd {
			result.Items = append(result.Items, item)
		}
	}
	c.mtx.Unlock()
	return &result
}

func (c *Item) historyPath() string {
	return c.ss.ServerDataPath() + "/history"
}

func (c *Item) readFiles(begin int64, end int64) []*common_interfaces.ItemValue {
	result := make([]*common_interfaces.ItemValue, 0)
	currentDT := begin
	for currentDT < end+86400*1000000 {
		dir := c.historyPath() + "/" + time.Unix(0, currentDT*1000).Format("2006-01-02")
		fullPath := dir + "/" + fmt.Sprintf("%016X", c.id) + ".bin"

		var file *FileCache
		if _, ok := c.files[fullPath]; !ok {
			file = NewFileCache(fullPath)
			c.files[fullPath] = file
		} else {
			file = c.files[fullPath]
		}

		result = append(result, file.Read(begin, end)...)
		currentDT += 86400 * 1000000
	}

	return result
}

type FlushResult struct {
	Error        error
	FullDataSize int
	CountOfItems int
}

func (c *Item) FinishFlush() {
	res := c.Flush()
	//logger.Println("HISTORY flushed ", c.id, "OK items:", res.CountOfItems, "data size:", res.FullDataSize)
	_ = res
	c.flushFinished = true
}

func (c *Item) Flush() FlushResult {

	var result FlushResult
	var err error
	var f *os.File
	var currentDir string

	type ItemToWrite struct {
		DT   time.Time
		data []byte
	}

	// Prepare data for writing
	c.mtx.Lock()

	items := make([]ItemToWrite, 0)
	itemsAsObjects := make([]*common_interfaces.ItemValue, 0)
	for _, item := range c.data {

		itemBinary := make([]byte, 1+1+8+1+len(item.Value)+1+len(item.UOM)+1)
		itemBinary[0] = 0xAA
		itemBinary[1] = byte(1 + 1 + 8 + 1 + len(item.Value) + 1 + len(item.UOM) + 1)
		binary.LittleEndian.PutUint64(itemBinary[1+1:], uint64(item.DT))
		itemBinary[1+1+8] = byte(len(item.Value))
		copy(itemBinary[1+1+8+1:], item.Value)
		itemBinary[1+1+8+1+len(item.Value)] = byte(len(item.UOM))
		copy(itemBinary[1+1+8+1+len(item.Value)+1:], item.UOM)
		itemBinary[1+1+8+1+len(item.Value)+1+len(item.UOM)] = 0x55

		items = append(items, ItemToWrite{DT: time.Unix(0, item.DT*1000), data: itemBinary})
		itemsAsObjects = append(itemsAsObjects, item)
	}
	c.data = make([]*common_interfaces.ItemValue, 0)

	var currentFilePath string

	// Writing
	for index, item := range items {
		dir := c.historyPath() + "/" + item.DT.Format("2006-01-02")
		if dir != currentDir {
			fullPath := dir + "/" + fmt.Sprintf("%016X", c.id) + ".bin"
			if f != nil {
				_ = f.Close()
				currentDir = ""
				f = nil
			}
			_ = os.MkdirAll(dir, 0755)
			f, err = os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				currentDir = dir
			}
			currentFilePath = fullPath
		}

		if file, ok := c.files[currentFilePath]; ok {
			file.Write(itemsAsObjects[index])
		}

		if err == nil {
			var written int
			var n int
			for n < len(item.data) {
				n, err = f.Write(item.data)
				if err != nil {
					result.Error = err
					logger.Println("error write to history file: ", err)
					break
				}
				written += n
			}

			result.FullDataSize += written
		}
	}

	result.CountOfItems = len(items)

	{
		found := true
		for found {
			found = false
			for key, file := range c.files {
				if time.Now().UTC().Sub(file.lastReadDT) > 5*time.Second {
					delete(c.files, key)
					found = true
					break
				}
			}
		}
	}

	c.mtx.Unlock()

	if f != nil {
		_ = f.Close()
		currentDir = ""
		f = nil
	}

	return result
}

func (c *Item) Remove() {
	logger.Println("History Item Remove", strconv.FormatUint(c.id, 16))

	c.mtx.Lock()
	c.data = nil
	c.files = nil

	historyDir := c.historyPath()

	dirs, err := utilities.GetDir(historyDir)
	if err == nil {
		for _, dir := range dirs {
			if dir.Dir {
				files, err := utilities.GetDir(dir.Path)
				if err == nil {
					for _, file := range files {
						if !file.Dir {
							if file.NameWithoutExt == fmt.Sprintf("%016X", c.id) {
								logger.Println("Item removing", file.Path)
								err = os.Remove(file.Path)
								if err != nil {
									logger.Println("Item removing", file.Path, "error", err)
								}
								logger.Println("Item removing", file.Path, "OK")
							}
						}
					}
				}

			}
		}
	}

	c.mtx.Unlock()
}

func (c *Item) CheckDepth() {
	historyDir := c.historyPath()

	dirs, err := utilities.GetDir(historyDir)
	if err == nil {
		for _, dir := range dirs {
			if dir.Dir {
				t, err := time.Parse("2006-01-02", dir.NameWithoutExt)
				if err == nil {
					if time.Now().UTC().Sub(t.UTC()) > time.Duration((c.historyDepthDays+1)*24)*time.Hour {
						files, err := utilities.GetDir(dir.Path)
						if err == nil {
							for _, file := range files {
								if !file.Dir {
									if file.NameWithoutExt == fmt.Sprintf("%016X", c.id) {
										logger.Println("Item CheckDepth removing", file.Path)
										err = os.Remove(file.Path)
										if err != nil {
											logger.Println("Item CheckDepth removing", file.Path, "error", err)
										}
										logger.Println("Item CheckDepth removing", file.Path, "OK")
									}
								}
							}
						}

					}
				}
			}
		}
	}
}
