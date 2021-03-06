package resources

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/settings"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Resources struct {
	mtx sync.Mutex
}

func NewResources() *Resources {
	var c Resources
	return &c
}

func (c *Resources) dir() string {
	return settings.ServerDataPath() + "/res"
}

func (c *Resources) fileName(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func (c *Resources) Rename(id string, newName string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var err error

	dir := c.dir()

	var bs []byte

	bs, err = ioutil.ReadFile(dir + "/" + id + ".info")
	if err != nil {
		return errors.New("no resource found")
	}

	var info common_interfaces.ResourcesItemInfo
	err = json.Unmarshal(bs, &info)
	if err != nil {
		return errors.New("no resource found")
	}

	info.Name = newName
	bs, _ = json.MarshalIndent(info, "", " ")
	err = ioutil.WriteFile(dir+"/"+id+".info", bs, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (c *Resources) Remove(id string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var err error

	dir := c.dir()

	_, err = ioutil.ReadFile(dir + "/" + id + ".info")
	if err != nil {
		return errors.New("no resource found")
	}

	err = os.Remove(dir + "/" + id + ".info")
	if err != nil {
		return err
	}
	err = os.Remove(dir + "/" + id + ".content")
	if err != nil {
		return err
	}

	return nil
}

func (c *Resources) Add(name string, tp string, content []byte) (string, error) {
	var err error
	c.mtx.Lock()
	defer c.mtx.Unlock()

	dir := c.dir()
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	filePathInfo := dir + "/" + "_"
	filePathContent := dir + "/" + "_"

	foundId := false
	var id string
	for i := 0; i < 100; i++ {
		id = uuid.New().String()
		filePathInfo = dir + "/" + id + ".info"
		filePathContent = dir + "/" + id + ".content"
		if _, err = os.Stat(filePathInfo); os.IsNotExist(err) {
			foundId = true
			break
		}
	}

	if !foundId {
		return "", errors.New("no id found")
	}

	var info common_interfaces.ResourcesItemInfo
	info.Id = id
	info.Name = name
	info.Type = tp
	bs, _ := json.MarshalIndent(info, "", " ")
	err = ioutil.WriteFile(filePathInfo, bs, 0666)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filePathContent, content, 0666)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *Resources) Set(id string, thumbnail []byte, content []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var err error
	var bs []byte

	dir := c.dir()

	bs, err = ioutil.ReadFile(dir + "/" + id + ".info")
	if err != nil {
		return errors.New("no resource found")
	}

	var info common_interfaces.ResourcesItemInfo
	err = json.Unmarshal(bs, &info)
	if err != nil {
		return errors.New("no resource found")
	}

	info.Thumbnail = thumbnail

	bs, _ = json.MarshalIndent(info, "", " ")
	err = ioutil.WriteFile(dir+"/"+id+".info", bs, 0666)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir+"/"+id+".content", content, 0666)
	if err != nil {
		return errors.New("can not save resource")
	}

	return nil
}

func (c *Resources) Get(id string) (*common_interfaces.ResourcesItem, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var err error
	var bs []byte

	dir := c.dir()

	bs, err = ioutil.ReadFile(dir + "/" + id + ".info")
	if err != nil {
		return nil, errors.New("no resource found")
	}

	var info common_interfaces.ResourcesItemInfo
	err = json.Unmarshal(bs, &info)
	if err != nil {
		return nil, errors.New("no resource found")
	}

	bs, err = ioutil.ReadFile(dir + "/" + id + ".content")
	if err != nil {
		return nil, errors.New("no resource found")
	}

	var result common_interfaces.ResourcesItem
	result.Info = info
	result.Content = bs
	return &result, nil
}

func SplitWithoutEmpty(req string, sep rune) []string {
	return strings.FieldsFunc(req, func(r rune) bool {
		return r == sep
	})
}

func (c *Resources) List(tp string, filter string, offset int, maxCount int) common_interfaces.ResourcesInfo {
	var result common_interfaces.ResourcesInfo
	result.Items = make([]common_interfaces.ResourcesItemInfo, 0)
	c.mtx.Lock()
	defer c.mtx.Unlock()

	filterParts := SplitWithoutEmpty(strings.ToLower(filter), ' ')

	//allItems := make([]common_interfaces.ResourcesItemInfo, 0)
	var err error
	var files []os.FileInfo
	files, err = ioutil.ReadDir(c.dir())
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".info") {
			var bs []byte
			bs, err = ioutil.ReadFile(c.dir() + "/" + file.Name())
			if err == nil {
				var info common_interfaces.ResourcesItemInfo
				err = json.Unmarshal(bs, &info)
				if err == nil {
					result.TotalCount++

					inFilter := 0
					for _, filterPart := range filterParts {
						if strings.Contains(strings.ToLower(info.Name), filterPart) {
							inFilter++
						}
					}
					if inFilter == len(filterParts) && (tp == "" || tp == info.Type) {
						if result.InFilterCount >= offset && len(result.Items) < maxCount {
							result.Items = append(result.Items, info)
						}
						result.InFilterCount++
					}
				}
			}
		}
	}

	return result
}
