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
	FuncCloudLogin               = "cloud_login"
	FuncCloudLogout              = "cloud_logout"
	FuncCloudState               = "cloud_state"
	FuncCloudNodes               = "cloud_nodes"
	FuncCloudAddNode             = "cloud_add_node"
	FuncCloudUpdateNode          = "cloud_update_node"
	FuncCloudRemoveNode          = "cloud_remove_node"
	FuncCloudGetSettings         = "cloud_get_settings"
	FuncCloudSetSettings         = "cloud_set_settings"
	FuncCloudAccountInfo         = "cloud_account_info"
	FuncCloudSetCurrentNodeId    = "cloud_set_current_node_id"
	FuncCloudGetSettingsProfiles = "cloud_get_settings_profiles"

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
	FuncServiceInfo        = "service_info"

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

func ApiFunctions() []string {
	res := make([]string, 0)
	res = append(res, FuncUnitTypeList)
	res = append(res, FuncUnitTypeCategories)
	res = append(res, FuncUnitTypeConfigMeta)

	res = append(res, FuncUnitAdd)
	res = append(res, FuncUnitRemove)
	res = append(res, FuncUnitState)
	res = append(res, FuncUnitStateAll)
	res = append(res, FuncUnitItemsValues)
	res = append(res, FuncUnitList)
	res = append(res, FuncUnitStart)
	res = append(res, FuncUnitStop)
	res = append(res, FuncUnitSetConfig)
	res = append(res, FuncUnitGetConfig)

	res = append(res, FuncDataItemList)
	res = append(res, FuncDataItemListAll)
	res = append(res, FuncDataItemWrite)
	res = append(res, FuncDataItemHistory)
	res = append(res, FuncDataItemHistoryChart)
	res = append(res, FuncDataItemRemove)

	res = append(res, FuncCloudLogin)
	res = append(res, FuncCloudLogout)
	res = append(res, FuncCloudState)
	res = append(res, FuncCloudNodes)
	res = append(res, FuncCloudAddNode)
	res = append(res, FuncCloudUpdateNode)
	res = append(res, FuncCloudRemoveNode)
	res = append(res, FuncCloudGetSettings)
	res = append(res, FuncCloudSetSettings)
	res = append(res, FuncCloudAccountInfo)
	res = append(res, FuncCloudSetCurrentNodeId)
	res = append(res, FuncCloudGetSettingsProfiles)

	res = append(res, FuncPublicChannelList)
	res = append(res, FuncPublicChannelAdd)
	res = append(res, FuncPublicChannelSetName)
	res = append(res, FuncPublicChannelRemove)
	res = append(res, FuncPublicChannelItemAdd)
	res = append(res, FuncPublicChannelItemRemove)
	res = append(res, FuncPublicChannelItemsState)
	res = append(res, FuncPublicChannelStart)
	res = append(res, FuncPublicChannelStop)

	res = append(res, FuncServiceLookup)
	res = append(res, FuncServiceStatistics)
	res = append(res, FuncServiceApi)
	res = append(res, FuncServiceSetNodeName)
	res = append(res, FuncServiceNodeName)
	res = append(res, FuncServiceInfo)

	res = append(res, FuncResourceAdd)
	res = append(res, FuncResourceSet)
	res = append(res, FuncResourceGet)
	res = append(res, FuncResourceRemove)
	res = append(res, FuncResourceRename)
	res = append(res, FuncResourceList)

	res = append(res, FuncSessionOpen)
	res = append(res, FuncSessionActivate)
	res = append(res, FuncSessionRemove)
	res = append(res, FuncSessionList)

	res = append(res, FuncUserList)
	res = append(res, FuncUserAdd)
	res = append(res, FuncUserSetPassword)
	res = append(res, FuncUserRemove)

	return res
}

type ApiRole struct {
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
}

func ApiRoles() []ApiRole {
	res := make([]ApiRole, 0)

	res = append(res, ApiRole{
		Code: "access_administrator",
		Name: "Access Administrator",
		Functions: []string{
			FuncUserList,
			FuncUserAdd,
			FuncUserSetPassword,
			FuncUserRemove,
			FuncSessionRemove,
			FuncSessionList,
		},
	})

	res = append(res, ApiRole{
		Code: "resource_administrator",
		Name: "Resource Administrator",
		Functions: []string{
			FuncResourceAdd,
			FuncResourceSet,
			FuncResourceRemove,
			FuncResourceRename,
		},
	})

	res = append(res, ApiRole{
		Code: "public_channels_administrator",
		Name: "Public Channels Administrator",
		Functions: []string{
			FuncPublicChannelAdd,
			FuncPublicChannelSetName,
			FuncPublicChannelRemove,
			FuncPublicChannelItemAdd,
			FuncPublicChannelItemRemove,
		},
	})

	res = append(res, ApiRole{
		Code: "read_only",
		Name: "ReadOnly",
		Functions: []string{
			FuncUnitTypeList,
			FuncUnitTypeCategories,
			FuncUnitTypeConfigMeta,
			FuncUnitState,
			FuncUnitStateAll,
			FuncUnitItemsValues,
			FuncUnitList,
			FuncDataItemList,
			FuncDataItemListAll,
			FuncDataItemHistory,
			FuncDataItemHistoryChart,
			FuncPublicChannelList,
			FuncPublicChannelItemsState,
			FuncServiceLookup,
			FuncServiceStatistics,
			FuncServiceNodeName,
			FuncResourceGet,
			FuncResourceList,
			FuncSessionOpen,
			FuncSessionActivate,
		},
	})

	return res
}
