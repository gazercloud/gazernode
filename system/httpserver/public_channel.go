package httpserver

/*
func (c *HttpServer) PublicChannelList(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelListRequest
	var resp nodeinterface.PublicChannelListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Channels, err = c.system.GetCloudChannels()
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelAdd(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelAddRequest
	var resp nodeinterface.PublicChannelAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.AddCloudChannel(req.ChannelName)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelSetName(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelSetNameRequest
	var resp nodeinterface.PublicChannelSetNameResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.EditCloudChannel(req.ChannelId, req.ChannelName)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelRemoveRequest
	var resp nodeinterface.PublicChannelRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.RemoveCloudChannel(req.ChannelId)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelItemAdd(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelItemAddRequest
	var resp nodeinterface.PublicChannelItemAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.CloudAddItems(req.Channels, req.Items)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelItemRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelItemRemoveRequest
	var resp nodeinterface.PublicChannelItemRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.CloudRemoveItems(req.Channels, req.Items)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelItemsState(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelItemsStateRequest
	var resp nodeinterface.PublicChannelItemsStateResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.UnitValues, err = c.system.GetCloudChannelValues(req.ChannelId)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelStart(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelStartRequest
	var resp nodeinterface.PublicChannelStartResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.StartPublicChannels(req.Ids)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) PublicChannelStop(request []byte) (response []byte, err error) {
	var req nodeinterface.PublicChannelStopRequest
	var resp nodeinterface.PublicChannelStopResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.StopPublicChannels(req.Ids)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
*/
