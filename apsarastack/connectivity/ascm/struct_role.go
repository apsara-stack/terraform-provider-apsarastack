package ascm

type Role struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		Active                 bool   `json:"active"`
		ArID                   string `json:"arId"`
		Code                   string `json:"code"`
		Default                bool   `json:"default"`
		Description            string `json:"description,omitempty"`
		Enable                 bool   `json:"enable"`
		ID                     int    `json:"id"`
		OrganizationVisibility string `json:"organizationVisibility"`
		OwnerOrganizationID    int    `json:"ownerOrganizationId"`
		RAMRole                bool   `json:"rAMRole"`
		RoleLevel              int64  `json:"roleLevel"`
		RoleName               string `json:"roleName"`
		RoleRange              string `json:"roleRange"`
		RoleType               string `json:"roleType"`
		UserCount              int    `json:"userCount"`
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
