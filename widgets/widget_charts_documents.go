package widgets

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/tree_items_parser"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"time"
)

type WidgetChartsDocuments struct {
	uicontrols.Panel
	selectedItem string
	tvItems      *uicontrols.TreeView
	client       *client.Client
	timer        *uievents.FormTimer
}

func NewWidgetChartsDocuments(parent uiinterfaces.Widget, client *client.Client) *WidgetChartsDocuments {
	var c WidgetChartsDocuments
	c.client = client
	c.InitControl(parent, &c)
	c.LoadItems()
	return &c
}

func (c *WidgetChartsDocuments) OnInit() {
	c.tvItems = c.AddTreeViewOnGrid(0, 0)
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

	c.timer = c.Window().NewTimer(500, func() {
		nodes := c.tvItems.VisibleNodes()
		for _, node := range nodes {
			c.tvItems.SetNodeValue(node, 1, c.client.GetItemValue(node.UserData.(string)).Value)
			c.tvItems.SetNodeValue(node, 2, c.client.GetItemValue(node.UserData.(string)).UOM)
			c.tvItems.SetNodeValue(node, 3, time.Unix(0, c.client.GetItemValue(node.UserData.(string)).DT*1000).Format("15:04:05"))
		}
	})
	c.timer.StartTimer()
}

func (c *WidgetChartsDocuments) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.tvItems = nil
	c.client = nil
	c.Panel.Dispose()
}

func (c *WidgetChartsDocuments) addNode(parentNode *uicontrols.TreeNode, item *tree_items_parser.TreeItem) *uicontrols.TreeNode {
	node := c.tvItems.AddNode(parentNode, item.ShortName)
	node.UserData = item.FullName
	for _, i := range item.Children {
		c.addNode(node, i)
	}
	return node
}

func (c *WidgetChartsDocuments) LoadItems() {
	c.client.GetAllItems(func(items []common_interfaces.ItemGetUnitItems, err error) {
		itemsNames := make([]string, 0)
		for _, i := range items {
			itemsNames = append(itemsNames, i.Name)
		}
		rootItem := tree_items_parser.ParseItems(itemsNames)
		rootNode := c.addNode(nil, rootItem)
		c.tvItems.ExpandNode(rootNode)
	})
}

func (c *WidgetChartsDocuments) SelectedItem() string {
	return c.selectedItem
}
