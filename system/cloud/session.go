package cloud

type SessionConfig struct {
	UserName string `json:"user_name"`
	Key      string `json:"key"`
	NodeId   string `json:"node_id"`

	AllowIncomingFunctions []string `json:"allow_incoming_functions"`
}
