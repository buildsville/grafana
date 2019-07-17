package dtos

type ShareSlack struct {
	RawURL    string `json:rawURL`
	Uid       string `json:uid`
	Slug      string `json:slug`
	PanelName string `json:panelName`
	Channel   string `json:channel`
	Param     string `json:param`
}
