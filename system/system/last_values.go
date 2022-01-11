package system

import (
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

func (c *System) WriteLastValues(items []*common_interfaces.Item) {
	c.mtx.Lock()
	bs, err := json.MarshalIndent(items, "", " ")
	c.mtx.Unlock()

	dir := c.ss.ServerDataPath() + "/last_values/"
	fullPath := dir + "/" + fmt.Sprintf("%016X", time.Now().UTC().UnixNano())
	_ = os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		_, _ = f.Write(bs)
		_ = f.Close()
	}
}

func (c *System) ReadLastValues() []*common_interfaces.Item {
	result := make([]*common_interfaces.Item, 0)
	dir := c.ss.ServerDataPath() + "/last_values/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result
	}

	ids := make([]int64, 0)
	for _, file := range files {
		iVal, err := strconv.ParseInt(file.Name(), 16, 64)
		if err == nil {
			ids = append(ids, iVal)
		}
	}

	if len(ids) < 1 {
		return result
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	lastId := ids[len(ids)-1]
	bs, err := ioutil.ReadFile(dir + "/" + fmt.Sprintf("%016X", lastId))
	if err == nil {
		_ = json.Unmarshal(bs, &result)
	}

	return result
}

func (c *System) RemoveOldLastValuesFiles() {
	dir := c.ss.ServerDataPath() + "/last_values/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	ids := make([]int64, 0)
	for _, file := range files {
		iVal, err := strconv.ParseInt(file.Name(), 16, 64)
		if err == nil {
			ids = append(ids, iVal)
		}
	}

	if len(ids) < 1 {
		return
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	idsToRemove := make([]int64, 0)
	for len(ids) > 3 {
		idsToRemove = append(idsToRemove, ids[0])
		ids = ids[1:]
	}

	for _, id := range idsToRemove {
		_ = os.Remove(dir + "/" + fmt.Sprintf("%016X", id))
	}
}
