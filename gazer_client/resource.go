package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) ResAdd(name string, tp string, content []byte) (string, error) {
	var call Call
	var req nodeinterface2.ResourceAddRequest
	req.Name = name
	req.Type = tp
	req.Content = content
	call.function = nodeinterface2.FuncResourceAdd
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceAddResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Id, err
}

func (c *GazerNodeClient) ResSet(id string, thumbnail []byte, content []byte) error {
	var call Call
	var req nodeinterface2.ResourceSetRequest
	req.Id = id
	req.Content = content
	call.function = nodeinterface2.FuncResourceSet
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceSetResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

/*func (c *GazerNodeClient) ResGet(id string) (*common_interfaces.ResourcesItem, error) {
	var call Call
	var req nodeinterface2.ResourceGetRequest
	req.Id = id
	call.function = nodeinterface2.FuncResourceGet
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceGetResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Item, err
}*/

func (c *GazerNodeClient) ResList(tp string, filter string, offset int, maxCount int) (common_interfaces.ResourcesInfo, error) {
	var call Call
	var req nodeinterface2.ResourceListRequest
	req.Type = tp
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount
	call.function = nodeinterface2.FuncResourceList
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceListResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) ResRemove(id string) error {
	var call Call
	var req nodeinterface2.ResourceRemoveRequest
	req.Id = id

	call.function = nodeinterface2.FuncResourceRemove
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceRemoveResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

/*func (c *GazerNodeClient) ResRename(id string, name string) error {
	var call Call
	var req nodeinterface2.ResourceRenameRequest
	req.Id = id
	//req.Name = name
	call.function = nodeinterface2.FuncResourceRename
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ResourceRenameResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}*/
