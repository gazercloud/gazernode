package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
)

type UnitAddRequest struct {
	UnitType string `json:"type"`
	UnitName string `json:"name"`
	Config   string `json:"config"`
}

type UnitAddResponse struct {
	UnitId string `json:"id"`
}

type UnitRemoveRequest struct {
	Units []string `json:"ids"`
}

type UnitRemoveResponse struct {
}

type UnitStateRequest struct {
	UnitId string `json:"id"`
}

type UnitStateResponse struct {
	UnitId          string                               `json:"id"`
	UnitDisplayName string                               `json:"name"`
	Type            string                               `json:"type"`
	TypeName        string                               `json:"type_name"`
	Status          string                               `json:"status"`
	MainItem        string                               `json:"main_item"`
	Value           string                               `json:"value"`
	UOM             string                               `json:"uom"`
	Items           []common_interfaces.ItemGetUnitItems `json:"items"`
}

type UnitStateAllRequest struct {
}

type UnitStateAllResponse struct {
	Items []UnitStateAllResponseItem `json:"items"`
}

type UnitStateAllResponseItem struct {
	UnitId          string `json:"id"`
	UnitDisplayName string `json:"name"`
	Type            string `json:"type"`
	TypeName        string `json:"type_name"`
	Status          string `json:"status"`
	MainItem        string `json:"main_item"`
	Value           string `json:"value"`
	UOM             string `json:"uom"`
}

type UnitItemsValuesRequest struct {
	UnitName string `json:"name"`
}

type UnitItemsValuesResponse struct {
	Items []common_interfaces.ItemGetUnitItems `json:"items"`
}

type UnitListRequest struct {
}

type UnitListResponseItem struct {
	Id             string `json:"id"`
	DisplayName    string `json:"name"`
	Type           string `json:"type"`
	TypeForDisplay string `json:"type_for_display"`
	Config         string `json:"config"`
	Enable         bool   `json:"enable"`
}

type UnitListResponse struct {
	Items []UnitListResponseItem `json:"items"`
}

type UnitStartRequest struct {
	Ids []string `json:"ids"`
}

type UnitStartResponse struct {
}

type UnitStopRequest struct {
	Ids []string `json:"ids"`
}

type UnitStopResponse struct {
}

type UnitSetConfigRequest struct {
	UnitId     string `json:"id"`
	UnitName   string `json:"name"`
	UnitConfig string `json:"config"`
}

type UnitSetConfigResponse struct {
}

type UnitGetConfigRequest struct {
	UnitId string `json:"id"`
}

type UnitGetConfigResponse struct {
	UnitId         string `json:"id"`
	UnitName       string `json:"name"`
	UnitType       string `json:"type"`
	UnitConfig     string `json:"config"`
	UnitConfigMeta string `json:"config_meta"`
}

type UnitPropSetRequest struct {
	UnitId string     `json:"unit_id"`
	Props  []PropItem `json:"props"`
}

type UnitPropSetResponse struct {
}

type UnitPropGetRequest struct {
	UnitId string `json:"unit_id"`
}

type UnitPropGetResponse struct {
	Props []PropItem `json:"props"`
}
