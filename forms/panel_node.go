package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/local_user_storage"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
)

type PanelNode struct {
	uicontrols.Panel

	client *client.Client
	//panelConnection *uicontrols.Panel
	panelBottom *uicontrols.Panel

	btnPanelUnits  *uicontrols.Button
	btnPanelCharts *uicontrols.Button
	btnPanelCloud  *uicontrols.Button
	btnPanelMaps   *uicontrols.Button
	btnPanelUsers  *uicontrols.Button

	btnSettings *uicontrols.Button

	panelLeftMenu *uicontrols.Panel

	panelUnits  *PanelUnits
	panelCloud  *PanelCloud
	panelCharts *PanelCharts
	panelMaps   *PanelMaps
	panelUsers  *PanelUsers

	currentButton *uicontrols.Button
	buttons       []*uicontrols.Button

	panelMain                    *uicontrols.Panel
	panelFullScreenValue         *PanelFullScreenValue
	controlBeforeFullScreenValue uiinterfaces.Widget

	lblStatistics *uicontrols.TextBlock
	lblAd         *uicontrols.TextBlock

	connectionIndex int

	timer *uievents.FormTimer

	imgConnectionOK    image.Image
	imgConnectionError image.Image

	imgBottomStatus *uicontrols.ImageBox
}

func NewPanelNode(parent uiinterfaces.Widget, client *client.Client, connectionIndex int) *PanelNode {
	var c PanelNode
	c.client = client
	c.connectionIndex = connectionIndex
	c.InitControl(parent, &c)

	c.client.OnSessionOpen = func() {
		var conn local_user_storage.NodeConnection
		conn.UserName = c.client.UserName()
		conn.Address = c.client.Address()
		conn.SessionToken = c.client.SessionToken()
		local_user_storage.Instance().SetConnection(c.connectionIndex, conn)
		c.FullRefresh()

		if MainFormInstance.currentNodeWidget == &c {
			c.client.ListOfUnits(func(items []nodeinterface.UnitListResponseItem, err error) {
				if len(items) == 0 {
					dialog := NewFormWizard(&c, c.client)
					dialog.ShowDialog()
					dialog.OnAccept = func() {
						c.FullRefresh()
					}
				}
			})
		}
	}

	c.client.OnSessionClose = func() {
		var conn local_user_storage.NodeConnection
		conn.UserName = c.client.UserName()
		conn.Address = c.client.Address()
		conn.SessionToken = c.client.SessionToken()
		local_user_storage.Instance().SetConnection(c.connectionIndex, conn)
		c.FullRefresh()
	}

	if c.client.SessionToken() == "" {
		c.client.SessionOpen(client.UserName(), client.Password(), nil)
	}

	c.UpdateStyle()

	return &c
}

func (c *PanelNode) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil
	c.panelBottom = nil

	c.btnPanelUnits = nil
	c.btnPanelCharts = nil
	c.btnPanelCloud = nil
	c.btnPanelUsers = nil
	c.btnSettings = nil

	c.panelUnits = nil
	c.panelCloud = nil
	c.panelCharts = nil

	c.currentButton = nil
	c.buttons = nil

	c.panelMain = nil
	c.panelFullScreenValue = nil
	c.controlBeforeFullScreenValue = nil

	c.lblStatistics = nil
	c.Panel.Dispose()
}

func (c *PanelNode) FullRefresh() {
	c.panelUnits.FullRefresh()
	c.panelCharts.FullRefresh()
	c.panelMaps.FullRefresh()
	c.panelCloud.FullRefresh()
	c.panelUsers.FullRefresh()
}

func (c *PanelNode) StylizeButton() {
	c.btnPanelUnits.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_image_blur_on_materialiconsoutlined_48dp_1x_outline_blur_on_black_48dp_png, c.btnPanelUnits.AccentColor()))
	c.btnPanelCloud.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialiconsoutlined_48dp_1x_outline_cloud_upload_black_48dp_png, c.btnPanelCloud.AccentColor()))
	c.btnPanelCharts.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_editor_stacked_line_chart_materialiconsoutlined_48dp_1x_outline_stacked_line_chart_black_48dp_png, c.btnPanelCharts.AccentColor()))
	c.btnPanelMaps.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_layers_materialiconsoutlined_48dp_1x_outline_layers_black_48dp_png, c.btnPanelCharts.AccentColor()))
	c.btnPanelUsers.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_perm_identity_materialiconsoutlined_48dp_1x_outline_perm_identity_black_48dp_png, c.btnPanelCharts.AccentColor()))

	for _, btn := range c.buttons {
		btn.SetBorders(0, color.White)
		if btn == c.btnPanelCloud {
			btn.SetBorderTop(1, c.ForeColor())
			btn.SetBorderBottom(1, c.ForeColor())
		} else {
			btn.SetBorderBottom(1, c.ForeColor())
		}
		if c.currentButton != nil {
			if btn.Text() == c.currentButton.Text() {
				btn.SetBorderLeft(3, c.AccentColor())
			} else {
				btn.SetBorderLeft(3, c.BackColor())
			}
		}
	}
	c.UpdateLayout()
}

func (c *PanelNode) OnInit() {
	c.SetPanelPadding(0)

	c.imgConnectionOK = uiresources.ResImgCol(uiresources.R_icons_material4_png_action_verified_user_materialicons_48dp_1x_baseline_verified_user_black_48dp_png, settings.GoodColor)
	c.imgConnectionOK = resize.Resize(24, 24, c.imgConnectionOK, resize.Bilinear)
	c.imgConnectionError = uiresources.ResImgCol(uiresources.R_icons_material4_png_action_verified_user_materialicons_48dp_1x_baseline_verified_user_black_48dp_png, colornames.Red)
	c.imgConnectionError = resize.Resize(24, 24, c.imgConnectionError, resize.Bilinear)

	c.panelFullScreenValue = NewPanelFullScreenValue(c, c.client, "")
	c.AddWidgetOnGrid(c.panelFullScreenValue, 0, 0)
	c.panelFullScreenValue.SetVisible(false)

	c.panelMain = c.AddPanelOnGrid(0, 0)
	c.panelMain.SetPanelPadding(0)
	c.panelLeftMenu = c.panelMain.AddPanelOnGrid(0, 0)
	c.panelLeftMenu.SetPanelPadding(0)
	c.panelLeftMenu.SetCellPadding(0)
	c.panelLeftMenu.SetMinWidth(100)
	c.panelLeftMenu.SetMaxWidth(100)

	c.btnPanelUnits = c.panelLeftMenu.AddButtonOnGrid(0, 0, "Units", func(event *uievents.Event) {
		c.panelUnits.SetVisible(true)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(false)
		c.panelUsers.SetVisible(false)
		c.currentButton = c.btnPanelUnits
		c.panelUnits.Activate()
		c.StylizeButton()
	})
	c.btnPanelUnits.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelUnits)

	c.btnPanelCharts = c.panelLeftMenu.AddButtonOnGrid(0, 2, "Charts", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(true)
		c.panelMaps.SetVisible(false)
		c.panelUsers.SetVisible(false)
		c.currentButton = c.btnPanelCharts
		c.StylizeButton()
	})
	c.btnPanelCharts.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelCharts)

	c.btnPanelMaps = c.panelLeftMenu.AddButtonOnGrid(0, 3, "Maps", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(true)
		c.panelUsers.SetVisible(false)
		c.currentButton = c.btnPanelMaps
		c.StylizeButton()
	})
	c.btnPanelMaps.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelMaps)

	c.panelLeftMenu.AddVSpacerOnGrid(0, 5)

	c.btnPanelCloud = c.panelLeftMenu.AddButtonOnGrid(0, 6, "Public\r\nChannels", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(true)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(false)
		c.panelUsers.SetVisible(false)
		c.currentButton = c.btnPanelCloud
		c.StylizeButton()
	})
	c.btnPanelCloud.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelCloud)

	c.btnPanelUsers = c.panelLeftMenu.AddButtonOnGrid(0, 7, "Users", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(false)
		c.panelUsers.SetVisible(true)
		c.currentButton = c.btnPanelUsers
		c.StylizeButton()
	})
	c.btnPanelUsers.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelUsers)

	c.StylizeButton()

	panelContent := c.panelMain.AddPanelOnGrid(1, 0)
	panelContent.SetName("PanelContent")
	panelContent.SetPanelPadding(0)
	c.panelUnits = NewPanelUnits(panelContent, c.client)
	panelContent.AddWidgetOnGrid(c.panelUnits, 0, 0)
	c.panelCloud = NewPanelCloud(panelContent, c.client)
	panelContent.AddWidgetOnGrid(c.panelCloud, 0, 1)
	c.panelCharts = NewPanelCharts(panelContent, c.client)
	panelContent.AddWidgetOnGrid(c.panelCharts, 0, 2)
	c.panelMaps = NewPanelMaps(panelContent, c.client)
	panelContent.AddWidgetOnGrid(c.panelMaps, 0, 3)
	c.panelUsers = NewPanelUsers(panelContent, c.client)
	panelContent.AddWidgetOnGrid(c.panelUsers, 0, 4)

	c.panelMaps.OnActionOpenMap = func(resId string) {
		c.client.ResGet(resId, func(item *common_interfaces.ResourcesItem, err error) {
			if err != nil {
				return
			}
			if item == nil {
				return
			}
			if item.Info.Type == "simple_map" {
				c.btnPanelMaps.Press()
				c.panelMaps.SelectMap(resId)
			}
			if item.Info.Type == "chart_group" {
				c.btnPanelCharts.Press()
				c.panelCharts.SelectChartGroup(resId)
			}
		})
	}

	c.panelUnits.SetVisible(false)
	c.panelCloud.SetVisible(false)
	c.panelCharts.SetVisible(false)
	c.panelMaps.SetVisible(false)
	c.panelUsers.SetVisible(false)

	c.panelUnits.SetPanelPadding(0)
	c.panelCloud.SetPanelPadding(0)
	c.panelCharts.SetPanelPadding(0)
	c.panelMaps.SetPanelPadding(0)
	c.panelUsers.SetPanelPadding(0)

	// Bottom
	c.panelBottom = c.AddPanelOnGrid(0, 1)
	c.panelBottom.SetBorderTop(1, c.ForeColor())
	c.panelBottom.SetMaxHeight(50)
	//c.panelBottom.SetBackColor(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	c.panelBottom.SetPanelPadding(0)
	c.panelBottom.SetCellPadding(0)
	c.btnSettings = c.panelBottom.AddButtonOnGrid(20, 0, "Settings", func(event *uievents.Event) {
		dialog := NewServiceDialog(c.panelBottom, c.client)
		dialog.ShowDialog()
		/*menu := uicontrols.NewPopupMenu(c)
		menu.AddItem("Statistics", func(event *uievents.Event) {
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_wysiwyg_materialiconsoutlined_48dp_1x_outline_wysiwyg_black_48dp_png, c.ForeColor()), "")
		menu.AddItem("Open gazer.cloud", func(event *uievents.Event) {

		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialiconsoutlined_48dp_1x_outline_open_in_browser_black_48dp_png, c.ForeColor()), "")
		menuItemTheme := menu.AddItem("Theme", func(event *uievents.Event) {
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.ForeColor()), "")
		{
			innerMenu := uicontrols.NewPopupMenu(c)
			innerMenu.AddItem("Dark", func(event *uievents.Event) {
				MainFormInstance.SetTheme(uistyles.StyleDarkBlue)
			}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.ForeColor()), "")
			innerMenu.AddItem("Light", func(event *uievents.Event) {
				MainFormInstance.SetTheme(uistyles.StyleLight)
			}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.ForeColor()), "")
			menuItemTheme.SetInnerMenu(innerMenu)
		}
		menu.AddItem("About", func(event *uievents.Event) {
			formAbout := NewFormAbout(c)
			formAbout.ShowDialog()
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_info_materialiconsoutlined_48dp_1x_outline_info_black_48dp_png, c.ForeColor()), "")
		_, menuPosY := c.btnSettings.RectClientAreaOnWindow()
		menu.ShowMenu(c.btnSettings.X(), menuPosY)
		menuPosY -= menu.Height()
		menu.SetX(c.Window().Width() - menu.Width() - 10)
		menu.SetY(menuPosY)*/
	})
	c.btnSettings.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_settings_applications_materialicons_48dp_1x_baseline_settings_applications_black_48dp_png, c.btnSettings.ForeColor()))
	c.btnSettings.SetBorders(5, color.RGBA{0, 0, 0, 0})
	c.btnSettings.SetShowText(false)

	c.imgBottomStatus = c.panelBottom.AddImageBoxOnGrid(1, 0, nil)
	c.imgBottomStatus.SetFixedSize(24, 24)

	c.lblStatistics = c.panelBottom.AddTextBlockOnGrid(2, 0, "---")
	c.lblStatistics.OnClick = func(ev *uievents.Event) {
		dialog := NewNodeConnectionDialog(c, c.client)
		dialog.OnAccept = func() {
			c.client.SessionOpen(dialog.Connection.UserName, dialog.Connection.Password, nil)
		}
		dialog.ShowDialog()
	}
	c.lblStatistics.SetMouseCursor(ui.MouseCursorPointer)
	c.lblStatistics.SetUnderline(true)
	c.lblStatistics.SetMinHeight(24)

	c.panelBottom.AddHSpacerOnGrid(5, 0)
	c.lblAd = c.panelBottom.AddTextBlockOnGrid(6, 0, "")
	c.lblAd.SetForeColor(settings.GoodColor)
	c.lblAd.SetUnderline(true)
	c.lblAd.OnClick = func(ev *uievents.Event) {
		client.OpenBrowser(adFromSite.Url)
	}
	c.lblAd.SetMouseCursor(ui.MouseCursorPointer)
	c.panelBottom.AddTextBlockOnGrid(7, 0, "  ")

	c.currentButton = c.btnPanelUnits
	c.StylizeButton()

	c.SetOnKeyDown(func(event *uievents.KeyDownEvent) bool {
		if event.Key == glfw.KeyF1 {
			c.btnPanelUnits.Press()
			return true
		}
		if event.Key == glfw.KeyF2 {
			c.btnPanelCloud.Press()
			return true
		}
		return false
	})
	c.timer = c.Window().NewTimer(1000, c.timerUpdate)
	c.timer.StartTimer()

	c.btnPanelUnits.Press()
}

func (c *PanelNode) ShowNavigation() {
	c.panelLeftMenu.SetVisible(true)
}

func (c *PanelNode) HideNavigation() {
	c.panelLeftMenu.SetVisible(false)
}

func (c *PanelNode) timerUpdate() {
	if c.client == nil {
		return
	}
	c.client.GetStatistics(func(statistics common_interfaces.Statistics, err error) {
		if c.lblStatistics != nil {
			if err == nil {
				c.lblStatistics.SetForeColor(settings.GoodColor)
				txt := "Connected (" + c.client.UserName() + "@" + c.client.Address() + ")"
				c.lblStatistics.SetText(Substr(txt, 0, 64))
				c.imgBottomStatus.SetImage(c.imgConnectionOK)
			} else {
				if c.client.SessionToken() == "" {
					c.lblStatistics.SetForeColor(colornames.Red)
					txt := "Authentication required (" + c.client.Address() + ")"
					c.lblStatistics.SetText(Substr(txt, 0, 64))
					c.imgBottomStatus.SetImage(c.imgConnectionError)
				} else {
					c.lblStatistics.SetForeColor(colornames.Red)
					txt := "Error: " + err.Error() + " (" + c.client.UserName() + "@" + c.client.Address() + ")"
					c.lblStatistics.SetText(Substr(txt, 0, 64))
					c.imgBottomStatus.SetImage(c.imgConnectionError)
				}
			}
			c.imgBottomStatus.SetMinWidth(48)
		}
	})
	c.lblAd.SetText(Substr(adFromSite.Text, 0, 64))
}

func (c *PanelNode) ShowFullScreenValue(show bool, itemId string) {
	if show {
		c.controlBeforeFullScreenValue = c.Window().FocusedWidget()
		c.panelFullScreenValue.SetItemId(itemId)
		c.panelFullScreenValue.SetVisible(true)
		c.panelFullScreenValue.Focus()
		c.panelMain.SetVisible(false)
	} else {

		c.panelFullScreenValue.SetVisible(false)
		c.panelMain.SetVisible(true)

		if c.controlBeforeFullScreenValue != nil {
			c.controlBeforeFullScreenValue.Focus()
		}
	}
}

type IShowFullScreen interface {
	ShowFullScreenValue(show bool)
}
