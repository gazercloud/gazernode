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

type CloudStateResponse struct {
	Connected bool   `json:"connected"`
	Status    string `json:"status"`
}
