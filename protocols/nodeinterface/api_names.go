package nodeinterface

const (
	// *** UnitType ***
	FuncUnitTypeList       = "unit_types_list"
	FuncUnitTypeCategories = "unit_categories"
	FuncUnitTypeConfigMeta = "unit_config_by_type"

	// *** Unit ***
	FuncUnitAdd         = "add_unit"
	FuncUnitRemove      = "remove_unit"
	FuncUnitState       = "unit_state"
	FuncUnitItemsValues = "unit_values"
	FuncUnitList        = "units"
	FuncUnitStart       = "start_units"
	FuncUnitStop        = "stop_units"
	FuncUnitSetConfig   = "set_unit_config"
	FuncUnitGetConfig   = "unit_config"

	// *** Data Item ***
	FuncDataItemList    = "items"
	FuncDataItemListAll = "all_items"
	FuncDataItemWrite   = "write"
	FuncDataItemHistory = "history"

	// *** Public Channel ***
	FuncPublicChannelList       = "cloud_channels"
	FuncPublicChannelAdd        = "add_cloud_channel"
	FuncPublicChannelSetName    = "edit_cloud_channel"
	FuncPublicChannelRemove     = "remove_cloud_channel"
	FuncPublicChannelItemAdd    = "cloud_add_items"
	FuncPublicChannelItemRemove = "cloud_remove_items"
	FuncPublicChannelItemsState = "cloud_channel_values"

	// *** Service ***
	FuncServiceLookup     = "lookup"
	FuncServiceStatistics = "statistics"

	// *** Resources ***
	FuncResourceAdd    = "res_add"
	FuncResourceSet    = "res_set"
	FuncResourceGet    = "res_get"
	FuncResourceRemove = "res_remove"
	FuncResourceRename = "res_rename"
	FuncResourceList   = "res_list"
)
