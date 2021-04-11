package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/cloud"
)

type PublicChannelListRequest struct {
}

type PublicChannelListResponse struct {
	Channels []cloud.ChannelInfo `json:"channels"`
}

type PublicChannelAddRequest struct {
	ChannelName string `json:"name"`
}

type PublicChannelAddResponse struct {
}

type PublicChannelSetNameRequest struct {
	ChannelId   string `json:"id"`
	ChannelName string `json:"name"`
}

type PublicChannelSetNameResponse struct {
}

type PublicChannelRemoveRequest struct {
	ChannelId string `json:"id"`
}

type PublicChannelRemoveResponse struct {
}

type PublicChannelItemAddRequest struct {
	Channels []string `json:"ids"`
	Items    []string `json:"items"`
}

type PublicChannelItemAddResponse struct {
}

type PublicChannelItemRemoveRequest struct {
	Channels []string `json:"ids"`
	Items    []string `json:"items"`
}

type PublicChannelItemRemoveResponse struct {
}

type PublicChannelItemsStateRequest struct {
	ChannelId string `json:"id"`
}

type PublicChannelItemsStateResponse struct {
	UnitValues []common_interfaces.Item `json:"unit_values"`
}
