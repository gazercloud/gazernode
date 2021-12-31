package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
)

type ResourceAddRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

type ResourceAddResponse struct {
	Id string `json:"id"`
}

type ResourceSetRequest struct {
	Id        string `json:"id"`
	Thumbnail []byte `json:"thumbnail"`
	Content   []byte `json:"content"`
}

type ResourceSetResponse struct {
}

type ResourceGetRequest struct {
	Id string `json:"id"`
}

type ResourceGetResponse struct {
	Item *common_interfaces.ResourcesItem `json:"item"`
}

type ResourceGetThumbnailRequest struct {
	Id string `json:"id"`
}

type ResourceGetThumbnailResponse struct {
	Item *common_interfaces.ResourcesItem `json:"item"`
}

type ResourceRemoveRequest struct {
	Id string `json:"id"`
}

type ResourceRemoveResponse struct {
}

type ResourceRenameRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ResourceRenameResponse struct {
}

type ResourceListRequest struct {
	Type     string `json:"type"`
	Filter   string `json:"filter"`
	Offset   int    `json:"offset"`
	MaxCount int    `json:"max_count"`
}

type ResourceListResponse struct {
	Items common_interfaces.ResourcesInfo `json:"items"`
}
