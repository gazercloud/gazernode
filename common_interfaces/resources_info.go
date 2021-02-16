package common_interfaces

type ResourcesInfo struct {
	TotalCount    int                 `json:"total_count"`
	InFilterCount int                 `json:"in_filter_count"`
	Items         []ResourcesItemInfo `json:"items"`
}
