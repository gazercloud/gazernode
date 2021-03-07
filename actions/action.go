package actions

type Action struct {
	Type    string `json:"type"`
	Comment string `json:"comment"`
	Content string `json:"content"`
}

type WriteItemAction struct {
	Item  string `json:"item"`
	Value string `json:"value"`
}

type OpenMapAction struct {
	ResId string `json:"res_id"`
}
