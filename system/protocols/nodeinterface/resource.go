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
	Id      string `json:"id"`
	Suffix  string `json:"suffix"`
	Offset  int64  `json:"offset"`
	Content []byte `json:"content"`
}

type ResourceSetResponse struct {
}

type ResourceGetRequest struct {
	Id     string `json:"id"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
}

type ResourceGetResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Offset  int64  `json:"offset"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
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

type ResourceListRequest struct {
	Type     string `json:"type"`
	Filter   string `json:"filter"`
	Offset   int    `json:"offset"`
	MaxCount int    `json:"max_count"`
}

type ResourceListResponse struct {
	Items common_interfaces.ResourcesInfo `json:"items"`
}

type ResourcePropSetRequest struct {
	Id    string     `json:"id"`
	Props []PropItem `json:"props"`
}

type ResourcePropSetResponse struct {
}

type ResourcePropGetRequest struct {
	Id string `json:"id"`
}

type ResourcePropGetResponse struct {
	Props []PropItem `json:"props"`
}
