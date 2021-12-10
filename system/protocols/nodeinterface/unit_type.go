package nodeinterface

type UnitTypeListRequest struct {
	Category string `json:"category"`
	Filter   string `json:"filter"`
	Offset   int    `json:"offset"`
	MaxCount int    `json:"max_count"`
}

type UnitTypeListResponseItem struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	DisplayName string `json:"display_name"`
	Help        string `json:"help"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
}

type UnitTypeListResponse struct {
	TotalCount    int                        `json:"total_count"`
	InFilterCount int                        `json:"in_filter_count"`
	Types         []UnitTypeListResponseItem `json:"types"`
}

type UnitTypeCategoriesRequest struct {
}

type UnitTypeCategoriesResponseItem struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Image       []byte `json:"image"`
}

type UnitTypeCategoriesResponse struct {
	Items []UnitTypeCategoriesResponseItem `json:"items"`
}

type UnitTypeConfigMetaRequest struct {
	UnitType string `json:"type"`
}

type UnitTypeConfigMetaResponse struct {
	UnitType           string `json:"type"`
	UnitTypeConfigMeta string `json:"config_meta"`
}
