package main

import "fmt"

func getStdUserName(userID uint16, userName string) string {
	return fmt.Sprintf("User: %s; ID: %d", userName, userID)
}

func getTabNameBasedOnBaseTextAndNoty(baseText string, noty uint) string {
	return fmt.Sprintf("%s (%d)", baseText, noty)
}

func updateTabsBasedOnNotyVals(tabs *MutableStructAboutTabs) {
	tabs.LinkTab.Text = getTabNameBasedOnBaseTextAndNoty(tabs.LinksBT, tabs.LinksNoty)
	tabs.ConnectionTab.Text = getTabNameBasedOnBaseTextAndNoty(tabs.ConnsBT, tabs.ConnsNoty)
	tabs.UsersTab.Text = getTabNameBasedOnBaseTextAndNoty(tabs.UsersBT, tabs.UsersNoty)
	tabs.AppTabs.Refresh()
}
