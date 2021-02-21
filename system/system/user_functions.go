package system

type User struct {
	Name         string
	PasswordHash string
}

var DefaultUserName string
var DefaultUserPassword string

func init() {
	DefaultUserName = "admin"
	DefaultUserPassword = "admin"
}

type UserSession struct {
	SessionId string
	UserName  string
}

func (c *System) UserAdd(name string, password string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) UserRename(oldName string, newName string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) UserSetPassword(name string, password string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) UserRemove(name string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) CheckSession(sessionId string) (string, error) {
	c.mtx.Lock()
	c.mtx.Unlock()
	return DefaultUserName, nil
}

func (c *System) OpenSession(name string, password string) (string, error) {
	c.mtx.Lock()
	c.mtx.Unlock()
	return "", nil
}
