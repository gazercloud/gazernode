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

type DialogRemoteAccessSettings struct {
	uicontrols.Dialog

	response *nodeinterface.CloudGetSettingsResponse

	btnAllow   *uicontrols.Button
	btnDeny    *uicontrols.Button
	btnDenyAll *uicontrols.Button
	btnRefresh *uicontrols.Button

	client  *client.Client
	lvItems *uicontrols.ListView

	lvProfiles      *uicontrols.ListView
	btnApplyProfile *uicontrols.Button

	btnOK *uicontrols.Button
}

func NewDialogRemoteAccessSettings(parent uiinterfaces.Widget, client *client.Client) *DialogRemoteAccessSettings {
	var c DialogRemoteAccessSettings
	c.client = client
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_hardware_security_materialiconsoutlined_48dp_1x_outline_security_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pTopButtons := pRight.AddPanelOnGrid(0, 0)
	pTopButtons.SetPanelPadding(0)

	c.btnAllow = pTopButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		selectedItems := c.lvItems.SelectedItems()
		for _, selectedItem := range selectedItems {
			selectedItem.SetValue(1, "allow")
		}
		c.updateListViewColors()
	})
	c.btnAllow.SetMinWidth(60)
	c.btnAllow.SetTooltip("Allow")

	c.btnDeny = pTopButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		selectedItems := c.lvItems.SelectedItems()
		for _, selectedItem := range selectedItems {
			selectedItem.SetValue(1, "deny")
		}
		c.updateListViewColors()
	})
	c.btnDeny.SetMinWidth(60)
	c.btnDeny.SetTooltip("Deny")

	c.btnDenyAll = pTopButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		for i := 0; i < c.lvItems.ItemsCount(); i++ {
			c.lvItems.Item(i).SetValue(1, "deny")
		}
		c.updateListViewColors()
	})
	c.btnDenyAll.SetMinWidth(60)
	c.btnDenyAll.SetTooltip("Deny All")

	pTopButtons.AddTextBlockOnGrid(3, 0, " | ")

	c.btnRefresh = pTopButtons.AddButtonOnGrid(4, 0, "", func(event *uievents.Event) {
		c.loadFunctions()
	})
	c.btnRefresh.SetTooltip("Refresh")

	pTopButtons.AddHSpacerOnGrid(10, 0)

	c.lvItems = pRight.AddListViewOnGrid(0, 1)
	c.lvItems.AddColumn("Function", 200)
	c.lvItems.AddColumn("Rule", 200)
	c.lvItems.OnSelectionChanged = func() {
		c.updateButtons()
	}

	pRightProfiles := pContent.AddPanelOnGrid(2, 0)

	pTopButtonsProfiles := pRightProfiles.AddPanelOnGrid(0, 0)
	pTopButtonsProfiles.SetPanelPadding(0)

	c.btnApplyProfile = pTopButtonsProfiles.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		selectedItems := c.lvProfiles.SelectedItems()
		if len(selectedItems) == 1 {
			profile := selectedItems[0].UserData("profile").(*nodeinterface.CloudGetSettingsProfilesResponseItem)
			for _, f := range profile.Functions {
				for i := 0; i < c.lvItems.ItemsCount(); i++ {
					if c.lvItems.Item(i).Value(0) == f {
						c.lvItems.Item(i).SetValue(1, "allow")
						break
					}
				}
			}
		}
		c.updateListViewColors()
	})
	c.btnApplyProfile.SetMinWidth(60)
	c.btnApplyProfile.SetTooltip("Apply profile")

	pTopButtonsProfiles.AddHSpacerOnGrid(10, 0)

	c.lvProfiles = pRightProfiles.AddListViewOnGrid(0, 1)
	c.lvProfiles.AddColumn("Profile Name", 300)
	c.lvProfiles.OnSelectionChanged = func() {
		c.updateButtons()
	}

	pButtons.AddHSpacerOnGrid(0, 0)

	btnApply := pButtons.AddButtonOnGrid(1, 0, "Apply", func(event *uievents.Event) {
		c.Apply(false)
	})
	btnApply.SetMinWidth(70)

	c.btnOK = pButtons.AddButtonOnGrid(2, 0, "OK", nil)
	c.TryAccept = func() bool {
		c.Apply(true)
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(3, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.OnShow = func() {
		c.lvItems.Focus()
	}

	c.UpdateStyle()

	return &c
}

func (c *DialogRemoteAccessSettings) updateButtons() {
	if len(c.lvItems.SelectedItems()) > 0 {
		c.btnAllow.SetEnabled(true)
		c.btnDeny.SetEnabled(true)
		c.btnDenyAll.SetEnabled(true)
	} else {
		c.btnAllow.SetEnabled(false)
		c.btnDeny.SetEnabled(false)
		c.btnDenyAll.SetEnabled(false)
	}

	if len(c.lvProfiles.SelectedItems()) > 0 {
		c.btnApplyProfile.SetEnabled(true)
	} else {
		c.btnApplyProfile.SetEnabled(false)
	}
}

func (c *DialogRemoteAccessSettings) Apply(closeAfter bool) {
	var req nodeinterface.CloudSetSettingsRequest
	for _, item := range c.response.Items {
		allow := false
		for i := 0; i < c.lvItems.ItemsCount(); i++ {
			if c.lvItems.Item(i).Value(0) == item.Function {
				if c.lvItems.Item(i).Value(1) == "allow" {
					allow = true
				}
				break
			}
		}
		req.Items = append(req.Items, &nodeinterface.CloudGetSettingsResponseItem{
			Function: item.Function,
			Allow:    allow,
		})
	}

	c.btnOK.SetEnabled(false)

	c.client.CloudSetSettings(req, func(response nodeinterface.CloudSetSettingsResponse, err error) {
		if err == nil {
			c.btnOK.SetEnabled(true)
			if closeAfter {
				c.TryAccept = nil
				c.Accept()
			}
		} else {
			c.btnOK.SetEnabled(true)
			uicontrols.ShowErrorMessage(c, err.Error(), "error")
		}
	})

}

func (c *DialogRemoteAccessSettings) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Remote Access Settings")
	c.Resize(800, 600)

	c.loadFunctions()
	c.loadProfiles()
}

func (c *DialogRemoteAccessSettings) loadFunctions() {
	c.client.CloudGetSettings(func(response nodeinterface.CloudGetSettingsResponse, err error) {
		if err == nil {
			c.response = &response
			c.lvItems.RemoveItems()
			for _, item := range response.Items {
				lvItem := c.lvItems.AddItem(item.Function)
				if item.Allow {
					lvItem.SetValue(1, "allow")
				} else {
					lvItem.SetValue(1, "deny")
				}
			}

			c.updateListViewColors()
		}
		c.updateButtons()
	})
}

func (c *DialogRemoteAccessSettings) loadProfiles() {
	c.client.CloudGetSettingsProfiles(func(response nodeinterface.CloudGetSettingsProfilesResponse, err error) {
		if err == nil {
			c.lvProfiles.RemoveItems()
			for i, item := range response.Items {
				lvItem := c.lvProfiles.AddItem(item.Name)
				lvItem.SetUserData("profile", response.Items[i])
			}
		}
		c.updateButtons()
	})
}

func (c *DialogRemoteAccessSettings) updateListViewColors() {
	for i := 0; i < c.lvItems.ItemsCount(); i++ {
		if c.lvItems.Item(i).Value(1) == "allow" {
			c.lvItems.Item(i).SetForeColorForCell(1, settings.GoodColor)
		} else {
			c.lvItems.Item(i).SetForeColorForCell(1, c.InactiveColor())
		}
	}
}

func (c *DialogRemoteAccessSettings) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAllow.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_check_circle_outline_materialicons_48dp_1x_baseline_check_circle_outline_black_48dp_png, activeColor))
	c.btnDeny.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_highlight_off_materialiconsoutlined_48dp_1x_outline_highlight_off_black_48dp_png, activeColor))
	c.btnDenyAll.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, activeColor))
	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))

	c.btnAllow.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_check_circle_outline_materialicons_48dp_1x_baseline_check_circle_outline_black_48dp_png, inactiveColor))
	c.btnDeny.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_action_highlight_off_materialiconsoutlined_48dp_1x_outline_highlight_off_black_48dp_png, inactiveColor))
	c.btnDenyAll.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, inactiveColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))

	c.btnApplyProfile.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_arrow_back_materialiconsoutlined_48dp_1x_outline_arrow_back_black_48dp_png, activeColor))
	c.btnApplyProfile.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_arrow_back_materialiconsoutlined_48dp_1x_outline_arrow_back_black_48dp_png, inactiveColor))

}
