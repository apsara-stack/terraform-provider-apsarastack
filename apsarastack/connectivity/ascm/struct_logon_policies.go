package ascm

type LogonPolicy struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		CuserID  string `json:"cuserId"`
		Default  bool   `json:"default"`
		Enable   bool   `json:"enable"`
		ID       int    `json:"id"`
		IPRanges []struct {
			IPRange       string `json:"ipRange"`
			LoginPolicyID int    `json:"loginPolicyId"`
			Protocol      string `json:"protocol"`
		} `json:"ipRanges"`
		LpID                   string `json:"lpId"`
		MuserID                string `json:"muserId"`
		Name                   string `json:"name"`
		OrganizationVisibility string `json:"organizationVisibility"`
		OwnerOrganizationID    int    `json:"ownerOrganizationId"`
		Rule                   string `json:"rule"`
		TimeRanges             []struct {
			EndTime       string `json:"endTime"`
			LoginPolicyID int    `json:"loginPolicyId"`
			StartTime     string `json:"startTime"`
		} `json:"timeRanges"`
		UserCount   int    `json:"userCount"`
		Description string `json:"description,omitempty"`
		Mtime       int64  `json:"mtime,omitempty"`
	} `json:"data"`
	Message  string `json:"message"`
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
	} `json:"pageInfo"`
	PureListData bool `json:"pureListData"`
	Redirect     bool `json:"redirect"`
	Success      bool `json:"success"`
}
