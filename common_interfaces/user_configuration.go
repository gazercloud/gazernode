package common_interfaces

type UserConfiguration struct {
	Name         string          `json:"name"`
	PasswordHash string          `json:"password_hash"`
	Properties   []*ItemProperty `json:"p"`
}
