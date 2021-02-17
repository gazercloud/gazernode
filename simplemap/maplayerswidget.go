package simplemap

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type MapLayersWidget struct {
	uicontrols.Panel

	view    *MapControlView
	lvItems *uicontrols.ListView

	btnAddLayer    *uicontrols.Button
	btnRemoveLayer *uicontrols.Button
	ctxMenu        *uicontrols.PopupMenu

	OnSelectedLayer func(layer *MapControlViewLayer)
	OnChangedLayer  func(layer *MapControlViewLayer)
}

func NewMapLayersWidget(parent uiinterfaces.Widget) *MapLayersWidget {
	var c MapLayersWidget
	c.Panel.InitControl(parent, &c)
	c.lvItems = uicontrols.NewListView(&c)
	c.lvItems.SetAnchors(uicontrols.ANCHOR_ALL)
	c.lvItems.AddColumn("Name", 100)
	c.lvItems.AddColumn("Current", 30)
	c.lvItems.AddColumn("Visible", 30)
	c.AddWidget(c.lvItems)

	/*c.btnAddLayer = uicontrols.NewButton(&c, 0, 0, 50, 22, "Add", func(event *uievents.Event) {
	})

	c.btnRemoveLayer = uicontrols.NewButton(&c, 60, 0, 50, 22, "Remove", func(event *uievents.Event) {
	})*/

	c.lvItems.OnSelectionChanged = func() {
		//selectedLayer := c.lvItems.SelectedItem().UserData.(*MapControlViewLayer)
	}

	c.ctxMenu = uicontrols.NewPopupMenu(&c)
	c.ctxMenu.AddItem("Add layer", func(event *uievents.Event) {
		c.view.AddLayer("New layer")
		c.updateList()
	}, nil, "")
	c.ctxMenu.AddItem("Remove layer", func(event *uievents.Event) {
		c.view.RemoveLayer(c.lvItems.SelectedItem().UserData("data").(*MapControlViewLayer))
		c.updateList()
	}, nil, "")
	c.ctxMenu.AddItem("Rename layer", func(event *uievents.Event) {
		layer := c.lvItems.SelectedItem().UserData("data").(*MapControlViewLayer)
		layer.name_ = "renamed"
		c.updateList()
	}, nil, "")
	c.ctxMenu.AddItem("Set as current", func(event *uievents.Event) {
		layer := c.lvItems.SelectedItem().UserData("data").(*MapControlViewLayer)
		if c.OnSelectedLayer != nil {
			c.OnSelectedLayer(layer)
		}
		c.updateList()
	}, nil, "")
	c.ctxMenu.AddItem("Set visible", func(event *uievents.Event) {
		layer := c.lvItems.SelectedItem().UserData("data").(*MapControlViewLayer)
		layer.visible_ = true
		if c.OnSelectedLayer != nil {
			c.OnChangedLayer(layer)
		}
		c.updateList()
	}, nil, "")
	c.ctxMenu.AddItem("Set invisible", func(event *uievents.Event) {
		layer := c.lvItems.SelectedItem().UserData("data").(*MapControlViewLayer)
		layer.visible_ = false
		if c.OnSelectedLayer != nil {
			c.OnChangedLayer(layer)
		}
		c.updateList()
	}, nil, "")

	c.lvItems.SetContextMenu(c.ctxMenu)

	return &c
}

func (c *MapLayersWidget) Dispose() {
	c.view = nil
	c.lvItems = nil

	c.btnAddLayer = nil
	c.btnRemoveLayer = nil
	c.ctxMenu = nil

	c.OnSelectedLayer = nil
	c.OnChangedLayer = nil

	c.Panel.Dispose()
}

func (c *MapLayersWidget) SetView(view *MapControlView) {
	c.view = view
	c.updateList()
}

func (c *MapLayersWidget) updateList() {
	c.lvItems.RemoveItems()
	if c.view == nil {
		return
	}

	for index, layer := range c.view.layers_ {
		lvItem := c.lvItems.AddItem(layer.name_)
		lvItem.SetUserData("data", layer)
		if layer == c.view.currentLayer() {
			c.lvItems.SetItemValue(index, 1, "Yes")
		} else {
			c.lvItems.SetItemValue(index, 1, "")
		}

		if layer.visible_ {
			c.lvItems.SetItemValue(index, 2, "Yes")
		} else {
			c.lvItems.SetItemValue(index, 2, "")
		}
	}
}
