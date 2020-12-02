package ascm

type Organization struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		Alias             string        `json:"alias"`
		CuserID           string        `json:"cuserId"`
		ID                int           `json:"id"`
		Internal          bool          `json:"internal"`
		Level             string        `json:"level"`
		MultiCloudStatus  string        `json:"multiCloudStatus"`
		MuserID           string        `json:"muserId"`
		Name              string        `json:"name"`
		ParentID          int           `json:"parentId"`
		SupportRegionList []interface{} `json:"supportRegionList"`
		UUID              string        `json:"uuid"`
		SupportRegions    string        `json:"supportRegions,omitempty"`
		Mtime             int64         `json:"mtime,omitempty"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	RequestID    string `json:"requestId"`
	Success      bool   `json:"success"`
}
