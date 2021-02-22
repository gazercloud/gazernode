package nodeinterface

type SessionOpenRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SessionOpenResponse struct {
	SessionToken string `json:"session_token"`
}
