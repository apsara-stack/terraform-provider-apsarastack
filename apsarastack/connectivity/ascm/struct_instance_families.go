package ascm

type InstanceFamily struct {
	Success bool `json:"success"`
	Data    []struct {
		ID          interface{} `json:"id"`
		GmtCreate   string      `json:"gmtCreate"`
		GmtModified string      `json:"gmtModified"`
		Creator     string      `json:"creator"`
		Modifier    string      `json:"modifier"`
		IsDeleted   string      `json:"isDeleted"`
		PageStart   int         `json:"pageStart"`
		PageSize    int         `json:"pageSize"`
		PageSort    string      `json:"pageSort"`
		PageOrder   string      `json:"pageOrder"`
		OrderBy     struct {
			ID string `json:"id"`
		} `json:"orderBy"`
		RegionID        interface{} `json:"regionId"`
		SpecFrom        interface{} `json:"specFrom"`
		BaseVersion     interface{} `json:"baseVersion"`
		Status          interface{} `json:"status"`
		SeriesID        string      `json:"seriesId"`
		SeriesName      string      `json:"seriesName"`
		SeriesNameLabel string      `json:"seriesNameLabel"`
		ResourceType    string      `json:"resourceType"`
		Deleted         bool        `json:"deleted"`
	} `json:"data"`
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	HTTPCode  interface{} `json:"httpCode"`
	IP        interface{} `json:"ip"`
	RequestID interface{} `json:"requestId"`
	HTTPOk    bool        `json:"httpOk"`
}
