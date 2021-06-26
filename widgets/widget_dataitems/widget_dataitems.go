package widget_dataitems

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/tree_items_parser"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"time"
)

type WidgetDataItems struct {
	uicontrols.Panel
	selectedItem string
	btnRefresh   *uicontrols.Button
	tvItems      *uicontrols.TreeView
	client       *client.Client
	timer        *uievents.FormTimer
	text1        string
	text2        string

	OnAdd func(id string)
}

func NewWidgetDataItems(parent uiinterfaces.Widget, client *client.Client, text1 string, text2 string) *WidgetDataItems {
	var c WidgetDataItems
	c.client = client
	c.text1 = text1
	c.text2 = text2
	c.InitControl(parent, &c)
	c.LoadItems()
	return &c
}

func (c *WidgetDataItems) OnInit() {
	panelToolbox := c.AddPanelOnGrid(0, 0)
	panelToolbox.SetPanelPadding(0)
	c.btnRefresh = panelToolbox.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		c.LoadItems()
	})
	c.btnRefresh.SetTooltip("Refresh")

	btnAdd := panelToolbox.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		if c.selectedItem != "" {
			if c.OnAdd != nil {
				c.OnAdd(c.selectedItem)
			}
		}
	})

	btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, c.AccentColor()))
	btnAdd.SetTooltip("Add selected item to the chart group")
	panelToolbox.AddHSpacerOnGrid(2, 0)

	c.tvItems = c.AddTreeViewOnGrid(0, 1)
	c.tvItems.AddColumn("Value", 70)
	c.tvItems.AddColumn("UOM", 70)
	c.tvItems.AddColumn("Time", 70)
	c.tvItems.OnSelectedNode = func(treeView *uicontrols.TreeView, node *uicontrols.TreeNode) {
		var ok bool
		c.selectedItem, ok = node.UserData.(string)
		if !ok {
			c.selectedItem = ""
		}
	}
	c.tvItems.OnBeginDrag = func(treeView *uicontrols.TreeView, node *uicontrols.TreeNode) interface{} {
		return node.UserData
	}

	c.AddTextBlockOnGrid(0, 2, c.text1)
	c.AddTextBlockOnGrid(0, 3, c.text2)

	c.timer = c.Window().NewTimer(500, func() {
		nodes := c.tvItems.VisibleNodes()
		for _, node := range nodes {
			value := c.client.GetItemValue(node.UserData.(string)).Value
			uom := c.client.GetItemValue(node.UserData.(string)).UOM
			c.tvItems.SetNodeValue(node, 1, value)
			c.tvItems.SetNodeValue(node, 2, uom)
			c.tvItems.SetNodeValue(node, 3, time.Unix(0, c.client.GetItemValue(node.UserData.(string)).DT*1000).Format("15:04:05"))
		}
	})
	c.timer.StartTimer()

}

func (c *WidgetDataItems) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.tvItems = nil
	c.client = nil
	c.Panel.Dispose()
}

func (c *WidgetDataItems) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *WidgetDataItems) addNode(parentNode *uicontrols.TreeNode, item *tree_items_parser.TreeItem, path string) *uicontrols.TreeNode {
	node := c.tvItems.AddNode(parentNode, item.ShortName)
	node.UserData = path
	for _, i := range item.Children {
		childPath := path
		if len(childPath) > 0 {
			childPath = childPath + "/" + i.ShortName
		} else {
			childPath = i.ShortName
		}
		c.addNode(node, i, childPath)
	}
	return node
}

func (c *WidgetDataItems) LoadItems() {
	c.client.GetAllItems(func(items []common_interfaces.ItemGetUnitItems, err error) {
		itemsNames := make([]string, 0)
		for _, i := range items {
			itemsNames = append(itemsNames, i.Name)
		}
		rootItem := tree_items_parser.ParseItems(itemsNames)
		c.tvItems.RemoveAllNodes()
		rootNode := c.addNode(nil, rootItem, "")
		c.tvItems.ExpandNode(rootNode)
	})
}

func (c *WidgetDataItems) SelectedItem() string {
	return c.selectedItem
}
