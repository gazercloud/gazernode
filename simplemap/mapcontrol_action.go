package simplemap

import "encoding/json"

type Action struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (c *MapControl) ExecAction(action *Action) error {
	c.actionElapsed = true
	switch action.Type {
	case "write_item":
		return c.execWriteItemAction(action.Content)
	case "open_map":
		return c.execOpenMapAction(action.Content)
	}

	return nil
}

func (c *MapControl) execWriteItemAction(content string) error {
	type WriteItemAction struct {
		Item  string `json:"item"`
		Value string `json:"value"`
	}
	var a WriteItemAction
	err := json.Unmarshal([]byte(content), &a)
	if err != nil {
		return err
	}

	c.MapWidget.ActionWriteItem(a.Item, a.Value)
	return nil
}

func (c *MapControl) execOpenMapAction(content string) error {
	type OpenMapAction struct {
		ResId string `json:"res_id"`
	}
	var a OpenMapAction
	err := json.Unmarshal([]byte(content), &a)
	if err != nil {
		return err
	}

	c.MapWidget.ActionOpenMap(a.ResId)
	return nil
}
