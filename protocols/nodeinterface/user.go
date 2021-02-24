package nodeinterface

type SessionOpenRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SessionOpenResponse struct {
	SessionToken string `json:"session_token"`
}

type SessionActivateRequest struct {
	SessionToken string `json:"session_token"`
}

type SessionActivateResponse struct {
	SessionToken string `json:"session_token"`
}
