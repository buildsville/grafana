package dtos

type ShareSlack struct {
	RawURL    string `json:rawURL`
	Uid       string `json:uid`
	Slug      string `json:slug`
	PanelId   int    `json:panelId`
	PanelName string `json:panelName`
	Channel   string `json:channel`
	From      int    `json:from`
	To        int    `json:to`
	Theme     string `json:theme`
}
