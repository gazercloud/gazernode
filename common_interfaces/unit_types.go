package common_interfaces

type UnitTypes struct {
	TotalCount    int            `json:"total_count"`
	InFilterCount int            `json:"in_filter_count"`
	Types         []UnitTypeInfo `json:"types"`
}
