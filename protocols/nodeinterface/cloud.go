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
	Allow bool   `json:"allow"`
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
	SessionKey       string `json:"session_key"`
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

type CloudGetSettingsResponseItem struct {
	Function string `json:"function"`
	Allow    bool   `json:"allow"`
}

type CloudGetSettingsResponse struct {
	Items []*CloudGetSettingsResponseItem `json:"items"`
}

type CloudSetSettingsRequest struct {
	Items []*CloudGetSettingsResponseItem `json:"items"`
}

type CloudSetSettingsResponse struct {
}

type CloudAccountInfoRequest struct {
}

type CloudAccountInfoResponse struct {
	Email         string `json:"email"`
	MaxNodesCount int64  `json:"max_nodes_count"`
}

type CloudSetCurrentNodeIdRequest struct {
	NodeId string `json:"node_id"`
}

type CloudSetCurrentNodeIdResponse struct {
}

type CloudGetSettingsProfilesRequest struct {
}

type CloudGetSettingsProfilesResponseItem struct {
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
}

type CloudGetSettingsProfilesResponse struct {
	Items []*CloudGetSettingsProfilesResponseItem `json:"items"`
}
