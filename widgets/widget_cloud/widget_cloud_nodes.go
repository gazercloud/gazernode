package widget_cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type WidgetCloudNodes struct {
	uicontrols.Panel
	client *client.Client
	timer  *uievents.FormTimer

	btnAdd        *uicontrols.Button
	btnRename     *uicontrols.Button
	btnRemove     *uicontrols.Button
	btnRefresh    *uicontrols.Button
	btnSetCurrent *uicontrols.Button

	lvItems *uicontrols.ListView

	OnNeedSetCurrent func(nodeId string)
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
	txtHeader.SetFontSize(16)

	pButtons := c.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		d := NewFormAddNode(c, c.client, "simple_map")
		d.ShowDialog()
		d.OnAccept = func() {
			c.loadNodes()
		}
	})
	c.btnAdd.SetTooltip("Add map ...")
	c.btnAdd.SetMinWidth(60)

	c.btnRename = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		item := c.lvItems.SelectedItems()[0]

		d := NewFormEditNode(c, c.client, item.TempData, item.Value(1))
		d.ShowDialog()
		d.OnAccept = func() {
			c.loadNodes()
		}
	})
	c.btnRename.SetTooltip("Rename chart group")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove selected node?", "Confirmation", func() {
			item := c.lvItems.SelectedItems()[0]
			c.client.CloudRemoveNode(item.TempData, func(resp nodeinterface.CloudRemoveNodeResponse, err error) {
				c.loadNodes()
			})
		}, nil)
	})
	c.btnRemove.SetTooltip("Remove selected map")

	c.btnSetCurrent = pButtons.AddButtonOnGrid(3, 0, "", func(event *uievents.Event) {
		if len(c.lvItems.SelectedItems()) != 1 {
			return
		}
		uicontrols.ShowQuestionMessageOKCancel(c, "Set node as a current?", "Confirmation", func() {
			item := c.lvItems.SelectedItems()[0]
			c.OnNeedSetCurrent(item.TempData)
		}, nil)
	})
	c.btnSetCurrent.SetTooltip("Set as a current node")

	pButtons.AddTextBlockOnGrid(4, 0, " | ")

	c.btnRefresh = pButtons.AddButtonOnGrid(5, 0, "", func(event *uievents.Event) {
		c.loadNodes()
	})
	c.btnRefresh.SetTooltip("Refresh")
	pButtons.AddHSpacerOnGrid(10, 0)

	pContent := c.AddPanelOnGrid(0, 2)
	pContent.SetPanelPadding(0)

	c.lvItems = pContent.AddListViewOnGrid(0, 0)
	c.lvItems.AddColumn("Id", 100)
	c.lvItems.AddColumn("Name", 100)

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

func (c *WidgetCloudNodes) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, activeColor))
	c.btnRename.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, activeColor))
	c.btnSetCurrent.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialiconsoutlined_48dp_1x_outline_add_black_48dp_png, inactiveColor))
	c.btnRename.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, inactiveColor))
	c.btnSetCurrent.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialiconsoutlined_48dp_1x_outline_create_black_48dp_png, inactiveColor))

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *WidgetCloudNodes) timerUpdate() {
	if !c.IsVisible() {
		return
	}

	c.client.CloudState(func(response nodeinterface.CloudStateResponse, err error) {
		if err != nil {
			return
		}
	})
}

func (c *WidgetCloudNodes) SetState(response nodeinterface.CloudStateResponse) {
	//logger.Println("WidgetCloudNodes SetState", response.NodeId)
	for i := 0; i < c.lvItems.ItemsCount(); i++ {
		if c.lvItems.Item(i).Value(0) == response.NodeId {
			c.lvItems.Item(i).SetForeColorForRow(settings.GoodColor)
		} else {
			c.lvItems.Item(i).SetForeColorForRow(nil)
		}
	}
}

func (c *WidgetCloudNodes) loadNodes() {
	c.client.CloudNodes(func(response nodeinterface.CloudNodesResponse, err error) {
		if err != nil {
			return
		}
		c.lvItems.RemoveItems()
		for _, node := range response.Nodes {
			item := c.lvItems.AddItem2(node.NodeId, node.Name)
			item.TempData = node.NodeId
		}
	})
}
