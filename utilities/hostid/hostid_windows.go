package hostid

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/settings"
	"github.com/google/uuid"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"time"
)

type HostId struct {
	UniqueId string    `json:"unique_id"`
	DT       time.Time `json:"dt"`
}

var hostId HostId

func InitHostId() {
	if LoadHostId() != nil {
		err := MakeHostId()
		if err != nil {
			logger.Println("Make Host Id error", err)
		}
	}
}

func LoadHostId() error {
	hostIdString, err := ioutil.ReadFile(settings.ServerDataPath() + "/host_id.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(hostIdString, &hostId)
	if err != nil {
		return err
	}
	return nil
}

func MakeHostId() error {
	var err error
	hostId.UniqueId, err = machineID()
	if err != nil {
		hostId.UniqueId = time.Now().UTC().Format("2006-01-02 15:04:05") + "_" + uuid.New().String()
	}
	hostId.DT = time.Now().UTC()
	var bs []byte
	bs, err = json.Marshal(hostId)
	if err == nil {
		err = ioutil.WriteFile(settings.ServerDataPath()+"/host_id.json", bs, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

// https://github.com/denisbrodbeck/machineid
func machineID() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}
	return s, nil
}
