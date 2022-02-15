package common_interfaces

type User struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Properties   map[string]*ItemProperty
}
