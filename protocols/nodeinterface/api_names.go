package nodeinterface

const (
	// *** UnitType ***
	FuncUnitTypeList       = "unit_type_list"
	FuncUnitTypeCategories = "unit_type_categories"
	FuncUnitTypeConfigMeta = "unit_type_config_meta"

	// *** Unit ***
	FuncUnitAdd         = "unit_add"
	FuncUnitRemove      = "unit_remove"
	FuncUnitState       = "unit_state"
	FuncUnitStateAll    = "unit_state_all"
	FuncUnitItemsValues = "unit_items_values"
	FuncUnitList        = "unit_list"
	FuncUnitStart       = "unit_start"
	FuncUnitStop        = "unit_stop"
	FuncUnitSetConfig   = "unit_set_config"
	FuncUnitGetConfig   = "unit_get_config"

	// *** Data Item ***
	FuncDataItemList         = "data_item_list"
	FuncDataItemListAll      = "data_item_list_all"
	FuncDataItemWrite        = "data_item_write"
	FuncDataItemHistory      = "data_item_history"
	FuncDataItemHistoryChart = "data_item_history_chart"
	FuncDataItemRemove       = "data_item_remove"

	// *** Cloud ***
	FuncCloudLogin            = "cloud_login"
	FuncCloudLogout           = "cloud_logout"
	FuncCloudState            = "cloud_state"
	FuncCloudNodes            = "cloud_nodes"
	FuncCloudAddNode          = "cloud_add_node"
	FuncCloudUpdateNode       = "cloud_update_node"
	FuncCloudRemoveNode       = "cloud_remove_node"
	FuncCloudGetSettings      = "cloud_get_settings"
	FuncCloudSetSettings      = "cloud_set_settings"
	FuncCloudAccountInfo      = "cloud_account_info"
	FuncCloudSetCurrentNodeId = "cloud_set_current_node_id"

	// *** Public Channel ***
	FuncPublicChannelList       = "public_channel_list"
	FuncPublicChannelAdd        = "public_channel_add"
	FuncPublicChannelSetName    = "public_channel_set_name"
	FuncPublicChannelRemove     = "public_channel_remove"
	FuncPublicChannelItemAdd    = "public_channel_item_add"
	FuncPublicChannelItemRemove = "public_channel_item_remove"
	FuncPublicChannelItemsState = "public_channel_item_state"
	FuncPublicChannelStart      = "public_channel_start"
	FuncPublicChannelStop       = "public_channel_stop"

	// *** Service ***
	FuncServiceLookup      = "service_lookup"
	FuncServiceStatistics  = "service_statistics"
	FuncServiceApi         = "service_api"
	FuncServiceSetNodeName = "service_set_node_name"
	FuncServiceNodeName    = "service_node_name"

	// *** Resource ***
	FuncResourceAdd    = "resource_add"
	FuncResourceSet    = "resource_set"
	FuncResourceGet    = "resource_get"
	FuncResourceRemove = "resource_remove"
	FuncResourceRename = "resource_rename"
	FuncResourceList   = "resource_list"

	// *** User ***
	FuncSessionOpen     = "session_open"
	FuncSessionActivate = "session_activate"
	FuncSessionRemove   = "session_remove"
	FuncSessionList     = "session_list"

	FuncUserList        = "user_list"
	FuncUserAdd         = "user_add"
	FuncUserSetPassword = "user_set_password"
	FuncUserRemove      = "user_remove"
)
