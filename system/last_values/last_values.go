package last_values

import (
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/settings"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

func Write(ss *settings.Settings, items []*common_interfaces.Item) {
	dir := ss.ServerDataPath() + "/last_values/"
	fullPath := dir + "/" + fmt.Sprintf("%016X", time.Now().UTC().UnixNano())
	_ = os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		bs, err := json.MarshalIndent(items, "", " ")
		if err == nil {
			f.Write(bs)
		}
		_ = f.Close()
	}
}

func Read(ss *settings.Settings) []*common_interfaces.Item {
	result := make([]*common_interfaces.Item, 0)
	dir := ss.ServerDataPath() + "/last_values/"
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
