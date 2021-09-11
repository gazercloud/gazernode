package cloud_accounts

import (
	"sort"
	"sync"
)

type CloudAccount struct {
	Email      string
	SessionKey string
}

type CloudAccounts struct {
	mtx         sync.Mutex
	accountsMap map[string]*CloudAccount
}

var c CloudAccounts

func init() {
	c.accountsMap = make(map[string]*CloudAccount)
}

func Set(email string, sessionKey string) {
	c.mtx.Lock()
	if v, ok := c.accountsMap[email]; ok {
		v.SessionKey = sessionKey
	} else {
		var ca CloudAccount
		ca.Email = email
		ca.SessionKey = sessionKey
		c.accountsMap[email] = &ca
	}
	c.mtx.Unlock()
}

func List() []*CloudAccount {
	result := make([]*CloudAccount, 0)
	c.mtx.Lock()
	for _, v := range c.accountsMap {
		var ca CloudAccount
		ca = *v
		result = append(result, &ca)
	}
	c.mtx.Unlock()
	sort.Slice(result, func(i, j int) bool {
		return result[i].Email < result[j].Email
	})
	return result
}
