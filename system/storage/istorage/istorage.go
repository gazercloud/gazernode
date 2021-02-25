package istorage

type IStorage interface {
	UserAdd(userName string, password string)
	UserSet(userName string, password string)
	UserRemove(userName string)
	UserList()
}

/*
Unit - string(time) -> guid
DataItem - int -> guid
PublicChannel - string from server -> +guid
Resource - guid
Session - string(token) -> +guid
User - name -> +guid

HistoryItem ->DataItemId+Timestamp
*/
