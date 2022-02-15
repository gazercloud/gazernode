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

type SessionRemoveRequest struct {
	SessionToken string `json:"session_token"`
}

type SessionRemoveResponse struct {
}

type SessionListRequest struct {
	UserName string `json:"user_name"`
}

type SessionListResponseItem struct {
	SessionToken    string `json:"session_token"`
	UserName        string `json:"user_name"`
	SessionOpenTime int64  `json:"session_open_time"`
}

type SessionListResponse struct {
	Items []SessionListResponseItem `json:"items"`
}

type UserListRequest struct {
}

type UserListResponse struct {
	Items []string `json:"items"`
}

type UserAddRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserAddResponse struct {
}

type UserSetPasswordRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserSetPasswordResponse struct {
}

type UserRemoveRequest struct {
	UserName string `json:"user_name"`
}

type UserRemoveResponse struct {
}

type UserPropSetRequest struct {
	UserName string     `json:"user_name"`
	Props    []PropItem `json:"props"`
}

type UserPropSetResponse struct {
}

type UserPropGetRequest struct {
	UserName string `json:"user_name"`
}

type UserPropGetResponse struct {
	Props []PropItem `json:"props"`
}
