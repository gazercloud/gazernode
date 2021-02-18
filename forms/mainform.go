package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/local_user_storage"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazernode/settings"
	"github.com/gazercloud/gazerui/go-gl/glfw/v3.3/glfw"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/gazercloud/gazerui/uistyles"
	"golang.org/x/image/colornames"
	"image/color"
)

type MainForm struct {
	uiforms.Form
	client      *client.Client
	panelBottom *uicontrols.Panel

	btnPanelUnits  *uicontrols.Button
	btnPanelCharts *uicontrols.Button
	btnPanelCloud  *uicontrols.Button
	btnPanelMaps   *uicontrols.Button
	btnSettings    *uicontrols.Button

	panelUnits  *PanelUnits
	panelCloud  *PanelCloud
	panelCharts *PanelCharts
	panelMaps   *PanelMaps

	currentButton *uicontrols.Button
	buttons       []*uicontrols.Button

	panelMain                    *uicontrols.Panel
	panelFullScreenValue         *PanelFullScreenValue
	controlBeforeFullScreenValue uiinterfaces.Widget

	lblStatistics *uicontrols.TextBlock
	lblAd         *uicontrols.TextBlock

	timer *uievents.FormTimer
}

var MainFormInstance *MainForm

func (c *MainForm) StylizeButton() {
	c.btnPanelUnits.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_image_blur_on_materialiconsoutlined_48dp_1x_outline_blur_on_black_48dp_png, c.btnPanelUnits.AccentColor()))
	c.btnPanelCloud.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialiconsoutlined_48dp_1x_outline_cloud_upload_black_48dp_png, c.btnPanelCloud.AccentColor()))
	c.btnPanelCharts.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_editor_stacked_line_chart_materialiconsoutlined_48dp_1x_outline_stacked_line_chart_black_48dp_png, c.btnPanelCharts.AccentColor()))
	c.btnPanelMaps.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_layers_materialiconsoutlined_48dp_1x_outline_layers_black_48dp_png, c.btnPanelCharts.AccentColor()))

	for _, btn := range c.buttons {
		btn.SetBorders(0, color.White)
		btn.SetBorderBottom(1, c.Panel().ForeColor())
		if c.currentButton != nil {
			if btn.Text() == c.currentButton.Text() {
				btn.SetBorderLeft(3, c.Panel().AccentColor())
			} else {
				btn.SetBorderLeft(3, c.Panel().BackColor())
			}
		}
	}
	c.UpdateLayout()
}

func (c *MainForm) OnInit() {
	MainFormInstance = c
	c.SetTitle("Gazer " + productinfo.Version())
	c.SetIcon(productinfo.Icon())

	winWidth := 1300
	winHeight := 700

	mon := glfw.GetPrimaryMonitor()
	if mon != nil {
		_, _, w, h := mon.GetWorkarea()
		winWidth = w - w/10
		winHeight = h - h/10
	}

	if winWidth < 1000 {
		winWidth = 1000
	}
	if winWidth > 1920 {
		winWidth = 1920
	}
	if winHeight < 500 {
		winHeight = 500
	}
	if winHeight > 1080 {
		winHeight = 1080
	}

	c.Resize(winWidth, winHeight)
	c.client = client.New(c)

	c.Panel().SetPanelPadding(0)

	c.panelFullScreenValue = NewPanelFullScreenValue(c.Panel(), c.client, "")
	c.Panel().AddWidgetOnGrid(c.panelFullScreenValue, 0, 0)
	c.panelFullScreenValue.SetVisible(false)

	c.panelMain = c.Panel().AddPanelOnGrid(0, 0)
	panelLeftMenu := c.panelMain.AddPanelOnGrid(0, 0)
	panelLeftMenu.SetPanelPadding(0)
	panelLeftMenu.SetMinWidth(100)
	panelLeftMenu.SetMaxWidth(100)

	c.btnPanelUnits = panelLeftMenu.AddButtonOnGrid(0, 0, "Units", func(event *uievents.Event) {
		c.panelUnits.SetVisible(true)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(false)
		c.currentButton = c.btnPanelUnits
		c.panelUnits.Activate()
		c.StylizeButton()
	})
	c.btnPanelUnits.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelUnits)

	c.btnPanelCloud = panelLeftMenu.AddButtonOnGrid(0, 1, "Cloud", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(true)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(false)
		c.currentButton = c.btnPanelCloud
		c.StylizeButton()
	})
	c.btnPanelCloud.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelCloud)

	c.btnPanelCharts = panelLeftMenu.AddButtonOnGrid(0, 2, "Charts", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(true)
		c.panelMaps.SetVisible(false)
		c.currentButton = c.btnPanelCharts
		c.StylizeButton()
	})
	c.btnPanelCharts.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelCharts)

	c.btnPanelMaps = panelLeftMenu.AddButtonOnGrid(0, 3, "Maps", func(event *uievents.Event) {
		c.panelUnits.SetVisible(false)
		c.panelCloud.SetVisible(false)
		c.panelCharts.SetVisible(false)
		c.panelMaps.SetVisible(true)
		c.currentButton = c.btnPanelMaps
		c.StylizeButton()
	})
	c.btnPanelMaps.SetMouseCursor(ui.MouseCursorPointer)
	c.buttons = append(c.buttons, c.btnPanelMaps)

	panelLeftMenu.AddVSpacerOnGrid(0, 5)

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

	c.panelUnits.SetVisible(false)
	c.panelCloud.SetVisible(false)
	c.panelCharts.SetVisible(false)
	c.panelMaps.SetVisible(false)

	c.panelUnits.SetPanelPadding(0)
	c.panelCloud.SetPanelPadding(0)
	c.panelCharts.SetPanelPadding(0)
	c.panelMaps.SetPanelPadding(0)

	// Bottom
	c.panelBottom = c.Panel().AddPanelOnGrid(0, 1)
	c.panelBottom.SetMaxHeight(50)
	//c.panelBottom.SetBackColor(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	c.panelBottom.SetPanelPadding(0)
	c.panelBottom.SetCellPadding(0)
	c.btnSettings = c.panelBottom.AddButtonOnGrid(0, 0, "Settings", func(event *uievents.Event) {
		menu := uicontrols.NewPopupMenu(c.Panel())
		menu.AddItem("Statistics", func(event *uievents.Event) {
			formStatistics := NewFormStatistics(c.Panel(), c.client)
			formStatistics.ShowDialog()
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_wysiwyg_materialiconsoutlined_48dp_1x_outline_wysiwyg_black_48dp_png, c.Panel().ForeColor()), "")
		menu.AddItem("Open gazer.cloud", func(event *uievents.Event) {
			client.OpenBrowser("https://gazer.cloud/?ref=menu_settings")
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialiconsoutlined_48dp_1x_outline_open_in_browser_black_48dp_png, c.Panel().ForeColor()), "")
		menuItemTheme := menu.AddItem("Theme", func(event *uievents.Event) {
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.Panel().ForeColor()), "")
		{
			innerMenu := uicontrols.NewPopupMenu(c.Panel())
			innerMenu.AddItem("Dark", func(event *uievents.Event) {
				c.SetTheme(uistyles.StyleDarkBlue)
			}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.Panel().ForeColor()), "")
			innerMenu.AddItem("Light", func(event *uievents.Event) {
				c.SetTheme(uistyles.StyleLight)
			}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_invert_colors_materialiconsoutlined_48dp_1x_outline_invert_colors_black_48dp_png, c.Panel().ForeColor()), "")
			menuItemTheme.SetInnerMenu(innerMenu)
		}
		menu.AddItem("About", func(event *uievents.Event) {
			formAbout := NewFormAbout(c.Panel())
			formAbout.ShowDialog()
		}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_info_materialiconsoutlined_48dp_1x_outline_info_black_48dp_png, c.Panel().ForeColor()), "")
		_, menuPosY := c.btnSettings.RectClientAreaOnWindow()
		menu.ShowMenu(c.btnSettings.X(), menuPosY)
		menuPosY -= menu.Height()
		menu.SetX(c.btnSettings.Width())
		menu.SetY(menuPosY)
	})
	c.btnSettings.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_settings_applications_materialicons_48dp_1x_baseline_settings_applications_black_48dp_png, c.btnSettings.ForeColor()))
	c.btnSettings.SetBorders(5, color.RGBA{0, 0, 0, 0})
	c.btnSettings.SetShowText(false)

	c.lblStatistics = c.panelBottom.AddTextBlockOnGrid(1, 0, "---")

	c.panelBottom.AddHSpacerOnGrid(5, 0)
	c.lblAd = c.panelBottom.AddTextBlockOnGrid(6, 0, "This is a beta version of the software!")
	c.lblAd.SetForeColor(colornames.Red)
	c.lblAd.SetUnderline(true)
	c.lblAd.OnClick = func(ev *uievents.Event) {
		client.OpenBrowser("https://gazer.cloud/ad/beta_version")
	}
	c.lblAd.SetMouseCursor(ui.MouseCursorPointer)
	c.panelBottom.AddTextBlockOnGrid(7, 0, "  ")
	txtBrand := c.panelBottom.AddTextBlockOnGrid(10, 0, "Gazer.Cloud")
	txtBrand.SetBorderRight(10, txtBrand.BackColor())
	txtBrand.OnClick = func(ev *uievents.Event) {
		client.OpenBrowser("https://gazer.cloud/?ref=app_bottom")
	}
	txtBrand.SetMouseCursor(ui.MouseCursorPointer)
	//txtBrand.SetBorderBottom(1, colornames.Blue)
	txtBrand.SetForeColor(color.RGBA{
		R: 0,
		G: 50,
		B: 255,
		A: 255,
	})

	c.currentButton = c.btnPanelUnits
	c.StylizeButton()

	c.Panel().SetOnKeyDown(func(event *uievents.KeyDownEvent) bool {
		if event.Key == glfw.KeyF1 {
			c.btnPanelUnits.Press()
			return true
		}
		if event.Key == glfw.KeyF2 {
			c.btnPanelCloud.Press()
			return true
		}
		if event.Key == glfw.KeyF9 {
			c.SetTheme(uistyles.StyleLight)
			return true
		}
		if event.Key == glfw.KeyF10 {
			c.SetTheme(uistyles.StyleDarkBlue)
			return true
		}
		return false
	})
	c.timer = c.NewTimer(1000, c.timerUpdate)
	c.timer.StartTimer()

	c.SetTheme(c.GetTheme())

	c.btnPanelUnits.Press()
}

func (c *MainForm) Dispose() {
	if c.timer != nil {
		c.timer.StopTimer()
	}
	c.timer = nil

	c.client = nil
	c.panelBottom = nil

	c.btnPanelUnits = nil
	c.btnPanelCharts = nil
	c.btnPanelCloud = nil
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
	c.Form.Dispose()
}

func (c *MainForm) GetTheme() int {
	theme := local_user_storage.Theme()
	if theme == "light" {
		return uistyles.StyleLight
	}
	if theme == "dark_blue" {
		return uistyles.StyleDarkBlue
	}
	return uistyles.StyleDarkBlue
}

func (c *MainForm) SetTheme(theme int) {
	uistyles.CurrentStyle = theme
	c.UpdateStyle()
	c.StylizeButton()

	themeStr := "dark_blue"
	if theme == uistyles.StyleLight {
		themeStr = "light"
	}
	if theme == uistyles.StyleDarkBlue {
		themeStr = "dark_blue"
	}

	local_user_storage.SetTheme(themeStr)
}

func (c *MainForm) OnClose() bool {
	return true
}

func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (c *MainForm) timerUpdate() {
	c.client.GetStatistics(func(statistics common_interfaces.Statistics, err error) {
		if err == nil {
			c.lblStatistics.SetForeColor(c.Panel().InactiveColor())
			c.lblStatistics.SetText("connected to local service")
		} else {
			c.lblStatistics.SetForeColor(settings.BadColor)
			c.lblStatistics.SetText(Substr(err.Error(), 0, 32))
		}
	})
}

func (c *MainForm) ShowFullScreenValue(show bool, itemId string) {
	if show {
		c.controlBeforeFullScreenValue = c.FocusedWidget()
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
