package nodeinterface

type CloudLoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type CloudLoginResponse struct {
}

type CloudLogoutRequest struct {
}

type CloudLogoutResponse struct {
}

type CloudStateRequest struct {
}

type CloudStateResponseItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type CloudStateResponse struct {
	UserName         string `json:"user_name"`
	NodeId           string `json:"node_id"`
	Connected        bool   `json:"connected"`
	LoggedIn         bool   `json:"logged_in"`
	LoginStatus      string `json:"login_status"`
	ConnectionStatus string `json:"connection_status"`
	IAmStatus        string `json:"i_am_status"`
	CurrentRepeater  string `json:"current_repeater"`
	Counters         []CloudStateResponseItem
}

type CloudNodesRequest struct {
}

type CloudNodesResponseItem struct {
	NodeId string
	Name   string
}

type CloudNodesResponse struct {
	Nodes []CloudNodesResponseItem `json:"nodes"`
}

type CloudAddNodeRequest struct {
	Name string `json:"name"`
}

type CloudAddNodeResponse struct {
	NodeId string `json:"node_id"`
}

type CloudUpdateNodeRequest struct {
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

type CloudUpdateNodeResponse struct {
}

type CloudRemoveNodeRequest struct {
	NodeId string `json:"node_id"`
}

type CloudRemoveNodeResponse struct {
}

type CloudGetSettingsRequest struct {
}

type CloudGetSettingsResponse struct {
	AllowWriteItem bool `json:"allow_write_item"`
}

type CloudSetSettingsRequest struct {
	NodeId         string `json:"node_id"`
	AllowWriteItem bool   `json:"allow_write_item"`
}

type CloudSetSettingsResponse struct {
}
