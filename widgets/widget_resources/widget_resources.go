package widget_resources

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/tree_items_parser"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type WidgetResources struct {
	uicontrols.Panel
	selectedItem string
	btnRefresh   *uicontrols.Button
	tvItems      *uicontrols.TreeView
	client       *client.Client

	nodeMaps   *uicontrols.TreeNode
	nodeCharts *uicontrols.TreeNode
}

func NewWidgetResources(parent uiinterfaces.Widget, client *client.Client) *WidgetResources {
	var c WidgetResources
	c.client = client
	c.InitControl(parent, &c)
	c.LoadItems()
	return &c
}

func (c *WidgetResources) OnInit() {
	panelToolbox := c.AddPanelOnGrid(0, 0)
	panelToolbox.SetPanelPadding(0)
	c.btnRefresh = panelToolbox.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		c.LoadItems()
	})
	c.btnRefresh.SetTooltip("Refresh")
	panelToolbox.AddHSpacerOnGrid(1, 0)
	c.tvItems = c.AddTreeViewOnGrid(0, 1)
	c.tvItems.SetColumnWidth(0, 300)
	c.tvItems.OnSelectedNode = func(treeView *uicontrols.TreeView, node *uicontrols.TreeNode) {
		var ok bool
		var resInfo *common_interfaces.ResourcesItemInfo
		resInfo, ok = node.UserData.(*common_interfaces.ResourcesItemInfo)
		if ok {
			c.selectedItem = resInfo.Id
		} else {
			c.selectedItem = ""
		}
	}
	rootNode := c.tvItems.AddNode(nil, "Root")
	c.nodeMaps = c.tvItems.AddNode(rootNode, "Maps")
	c.nodeCharts = c.tvItems.AddNode(rootNode, "Charts")

	c.UpdateStyle()
}

func (c *WidgetResources) Dispose() {
	c.tvItems = nil
	c.client = nil
	c.Panel.Dispose()
}

func (c *WidgetResources) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *WidgetResources) addNode(parentNode *uicontrols.TreeNode, item *tree_items_parser.TreeItem, path string) *uicontrols.TreeNode {
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

func (c *WidgetResources) LoadItems() {
	c.tvItems.RemoveNodes(c.nodeMaps)
	c.tvItems.RemoveNodes(c.nodeCharts)

	c.client.ResList("simple_map", "", 0, 100000, func(info common_interfaces.ResourcesInfo, err error) {
		c.tvItems.RemoveNodes(c.nodeMaps)
		for _, item := range info.Items {
			node := c.tvItems.AddNode(c.nodeMaps, item.Name)
			itemForPointer := item
			node.UserData = &itemForPointer
		}
	})
	c.client.ResList("chart_group", "", 0, 100000, func(info common_interfaces.ResourcesInfo, err error) {
		c.tvItems.RemoveNodes(c.nodeCharts)
		for _, item := range info.Items {
			node := c.tvItems.AddNode(c.nodeCharts, item.Name)
			itemForPointer := item
			node.UserData = &itemForPointer
		}
	})
}

func (c *WidgetResources) SelectedItem() string {
	return c.selectedItem
}
