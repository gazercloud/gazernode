package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"strconv"
)

type WidgetCloudNodes struct {
	uicontrols.Panel
	client         *client.Client
	timer          *uievents.FormTimer
	accountLoading bool
	accountLoaded  bool
	lastSessionKey string

	btnAdd        *uicontrols.Button
	btnRename     *uicontrols.Button
	btnRemove     *uicontrols.Button
	btnRefresh    *uicontrols.Button
	btnSetCurrent *uicontrols.Button

	lvItems *uicontrols.ListView
	//lvAccountInfo *uicontrols.ListView

	lblAccountInfoEmail         *uicontrols.TextBlock
	lblAccountInfoMaxNodesCount *uicontrols.TextBlock

	menuNodes *uicontrols.PopupMenu

	OnNeedToConnect func(nodeId string, sessionKey string)
}

func NewWidgetCloudNodes(parent uiinterfaces.Widget, client *client.Client) *WidgetCloudNodes {
	var c WidgetCloudNodes
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	return &c
}

func (c *WidgetCloudNodes) OnInit() {
	pHeader := c.AddPanelOnGrid(0, 0)
	pHeader.SetPanelPadding(0)

	txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Nodes")
	txtHeader.SetForeColor(c.AccentColor())
	txtHeader.SetFontSize(c.FontSize() * 1.2)

	pButtons := c.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		c.addNode()
	})
	c.btnAdd.SetTooltip("Add node")
	c.btnAdd.SetMinWidth(60)

	c.btnRename = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		c.updateNode()
	})
	c.btnRename.SetTooltip("Rename node")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		c.removeNode()
	})
	c.btnRemove.SetTooltip("Remove node")

	c.btnSetCurrent = pButtons.AddButtonOnGrid(3, 0, "", func(event *uievents.Event) {
		c.setAsCurrentNode()
	})
	c.btnSetCurrent.SetTooltip("Set as a current node")

	pButtons.AddTextBlockOnGrid(4, 0, " | ")

	c.btnRefresh = pButtons.AddButtonOnGrid(5, 0, "", func(event *uievents.Event) {
		c.refresh()
	})
	c.btnRefresh.SetTooltip("Refresh")
	pButtons.AddHSpacerOnGrid(10, 0)

	pContent := c.AddPanelOnGrid(0, 2)
	pContent.SetPanelPadding(0)

	c.lvItems = pContent.AddListViewOnGrid(0, 0)
	c.lvItems.AddColumn("Id", 100)
	c.lvItems.AddColumn("Name", 200)

	c.menuNodes = uicontrols.NewPopupMenu(c.lvItems)
	c.menuNodes.AddItemWithUiResImage("Add node ...", func(event *uievents.Event) {
		c.addNode()
	}, uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, "")
	c.menuNodes.AddItemWithUiResImage("Rename node ...", func(event *uievents.Event) {
		c.updateNode()
	}, uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, "")
	c.menuNodes.AddItemWithUiResImage("Remove node ...", func(event *uievents.Event) {
		c.removeNode()
	}, uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, "")
	c.menuNodes.AddItemWithUiResImage("Set as a current node", func(event *uievents.Event) {
		c.setAsCurrentNode()
	}, uiresources.R_icons_material4_png_maps_pin_drop_materialiconsoutlined_48dp_1x_outline_pin_drop_black_48dp_png, "")
	c.menuNodes.AddItemWithUiResImage("Refresh", func(event *uievents.Event) {
		c.refresh()
	}, uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, "")
	c.menuNodes.AddItemWithUiResImage("Connect ...", func(event *uievents.Event) {
		c.connect()
	}, uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, "")
	c.lvItems.SetContextMenu(c.menuNodes)

	lblAccountInfo := pContent.AddTextBlockOnGrid(0, 1, "Account information")
	lblAccountInfo.SetForeColor(c.AccentColor())
	lblAccountInfo.SetFontSize(c.FontSize() * 1.2)

	/*c.lvAccountInfo = pContent.AddListViewOnGrid(0, 2)
	c.lvAccountInfo.AddColumn("Parameter", 200)
	c.lvAccountInfo.AddColumn("Value", 200)*/

	pAccountInfo := pContent.AddPanelOnGrid(0, 2)
	pAccountInfo.SetBorders(1, c.InactiveColor())
	pAccountInfo.AddTextBlockOnGrid(0, 0, "Account's Email: ")
	c.lblAccountInfoEmail = pAccountInfo.AddTextBlockOnGrid(1, 0, "")
	pAccountInfo.AddTextBlockOnGrid(0, 1, "Max Number Of Nodes: ")
	c.lblAccountInfoMaxNodesCount = pAccountInfo.AddTextBlockOnGrid(1, 1, "")
	pAccountInfo.AddHSpacerOnGrid(2, 0)

	c.UpdateStyle()
}

func (c *WidgetCloudNodes) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil
	c.client = nil

	c.Panel.Dispose()
}

func (c *WidgetCloudNodes) addNode() {
	d := NewFormAddNode(c, c.client)
	d.ShowDialog()
	d.OnAccept = func() {
		c.loadNodes()
	}
}

func (c *WidgetCloudNodes) updateNode() {
	if len(c.lvItems.SelectedItems()) != 1 {
		return
	}
	item := c.lvItems.SelectedItems()[0]

	d := NewFormEditNode(c, c.client, item.TempData, item.Value(1))
	d.ShowDialog()
	d.OnAccept = func() {
		c.loadNodes()
	}
}

func (c *WidgetCloudNodes) removeNode() {
	if len(c.lvItems.SelectedItems()) != 1 {
		return
	}
	uicontrols.ShowQuestionMessageOKCancel(c, "Remove selected node?", "Confirmation", func() {
		item := c.lvItems.SelectedItems()[0]
		c.client.CloudRemoveNode(item.TempData, func(resp nodeinterface.CloudRemoveNodeResponse, err error) {
			c.loadNodes()
		})
	}, nil)

}

func (c *WidgetCloudNodes) setAsCurrentNode() {
	if len(c.lvItems.SelectedItems()) != 1 {
		return
	}
	uicontrols.ShowQuestionMessageOKCancel(c, "Set node as a current?", "Confirmation", func() {
		item := c.lvItems.SelectedItems()[0]
		c.client.CloudSetCurrentNodeId(item.TempData, func(response nodeinterface.CloudSetCurrentNodeIdResponse, err error) {
		})
	}, nil)
}

func (c *WidgetCloudNodes) refresh() {
	c.loadNodes()
}

func (c *WidgetCloudNodes) connect() {
	if len(c.lvItems.SelectedItems()) != 1 {
		return
	}
	item := c.lvItems.SelectedItems()[0]
	if c.OnNeedToConnect != nil {
		c.OnNeedToConnect(item.TempData, c.lastSessionKey)
	}
}

func (c *WidgetCloudNodes) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, activeColor))
	c.btnRename.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, activeColor))
	c.btnSetCurrent.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_pin_drop_materialiconsoutlined_48dp_1x_outline_pin_drop_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, inactiveColor))
	c.btnRename.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, inactiveColor))
	c.btnSetCurrent.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_pin_drop_materialiconsoutlined_48dp_1x_outline_pin_drop_black_48dp_png, inactiveColor))

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *WidgetCloudNodes) SetState(response nodeinterface.CloudStateResponse) {
	c.lastSessionKey = response.SessionKey

	for i := 0; i < c.lvItems.ItemsCount(); i++ {
		if c.lvItems.Item(i).Value(0) == response.NodeId {
			c.lvItems.Item(i).SetForeColorForRow(settings.GoodColor)
		} else {
			c.lvItems.Item(i).SetForeColorForRow(nil)
		}
	}

	if !c.accountLoaded {
		c.loadNodes()
	}
}

func (c *WidgetCloudNodes) loadNodes() {
	if c.accountLoading {
		return
	}

	c.accountLoading = true

	c.client.CloudAccountInfo(func(response nodeinterface.CloudAccountInfoResponse, err error) {
		if err != nil {
			c.accountLoading = false
			return
		}
		if c.lblAccountInfoEmail == nil || c.lblAccountInfoMaxNodesCount == nil {
			c.accountLoading = false
			return
		}

		c.lblAccountInfoEmail.SetText(response.Email)
		c.lblAccountInfoMaxNodesCount.SetText(strconv.FormatInt(response.MaxNodesCount, 10))

		c.client.CloudNodes(func(response nodeinterface.CloudNodesResponse, err error) {
			if err != nil {
				return
			}
			if c.lvItems == nil {
				c.accountLoading = false
				return
			}
			c.lvItems.RemoveItems()
			for _, node := range response.Nodes {
				item := c.lvItems.AddItem2(node.NodeId, node.Name)
				item.TempData = node.NodeId
			}
			c.accountLoaded = true
			c.accountLoading = false
		})

	})
}
